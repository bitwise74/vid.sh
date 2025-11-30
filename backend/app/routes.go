package app

import (
	"bitwise74/video-api/app/ffmpeg"
	"bitwise74/video-api/app/file"
	"bitwise74/video-api/app/mail"
	"bitwise74/video-api/app/profile"
	"bitwise74/video-api/app/root"
	"bitwise74/video-api/app/user"
	"bitwise74/video-api/aws"
	"bitwise74/video-api/db"
	"bitwise74/video-api/internal/redis"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/internal/types"
	"bitwise74/video-api/pkg/middleware"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() (*gin.Engine, error) {
	d := &types.Dependencies{
		JobQueue: service.NewJobQueue(),
	}

	router := gin.New()

	db, err := db.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite database, %w", err)
	}
	d.DB = db

	origins := strings.Split(os.Getenv("HOST_CORS"), ",")

	router.Use(
		cors.New(cors.Config{
			AllowOrigins:     origins,
			AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "TurnstileToken", "Range", "Access-Control-Allow-Headers", "auth_token"},
			ExposeHeaders:    []string{"Content-Length", "Content-Range"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		gin.Recovery(),
		middleware.NewRequestIDMiddleware(),
	)

	if os.Getenv("APP_GIN_LOGS") == "true" {
		router.Use(gin.Logger())
	}

	router.HandleMethodNotAllowed = true
	router.RedirectFixedPath = true

	rateLimit, _ := strconv.Atoi(os.Getenv("SECURITY_RATE_LIMIT"))
	bodySizeLimit, _ := strconv.Atoi(os.Getenv("UPLOAD_MAX_SIZE"))

	jwt := middleware.NewJWTMiddleware(db.Gorm)
	turnstile := middleware.NewTurnstileMiddleware()
	bodySizeLimiter := middleware.NewBodySizeLimiter(int64(bodySizeLimit))
	rateLimiter := middleware.RateLimiterMiddleware(middleware.RateLimiterConfig{
		RequestsPerSecond: rateLimit,
		Burst:             rateLimit * 2,
		CleanupInterval:   time.Second,
	})

	m := router.Group("/api", rateLimiter)
	{
		// HEAD /api/heartbeat 		-> Used to check if the server is alive
		m.HEAD("/heartbeat", root.Heartbeat)

		// HEAD /api/validate		-> Validates a JWT token
		m.GET("/validate", jwt, root.Validate)

		// GET /api/oembed		-> Returns oembed data for discord embeds
		m.GET("/oembed", root.OEmbed)
	}

	u := m.Group("/users")
	{
		// GET /api/users		-> Returns the basic info of a user
		u.GET("", jwt, func(c *gin.Context) { user.Fetch(c, d) })

		// POST /api/users 		-> Registers a new user
		u.POST("", func(c *gin.Context) { user.Register(c, d) })

		// POST /api/users/login 	-> Logs in a user and returns a JWT token
		u.POST("/login", func(c *gin.Context) { user.Login(c, d) })

		// POST /api/users/verify	-> Verifies a new user
		u.POST("/verify", func(c *gin.Context) { user.Verify(c, d) })

		// POST /api/users/logout	-> Deletes all cookies to logout
		u.POST("/logout", jwt, func(c *gin.Context) { user.Logout(c, d) })

		// DELETE /api/users/:id 	-> Deletes a user account
		u.DELETE("/:id", jwt, func(c *gin.Context) { user.Delete(c, d) })

		// PATCH /api/users/update		-> Updates a user's profile
		u.PATCH("/update", jwt, func(c *gin.Context) { user.Update(c, d) })

		// POST /api/users/reset-password	-> Resets a user's password
		u.POST("/reset-password", func(c *gin.Context) { user.ResetPassword(c, d) })
	}

	ma := m.Group("/mail")
	{
		// POST /api/mail/verify	-> Sends a verification email
		ma.POST("/verify", jwt, func(c *gin.Context) { mail.VerificationMail(c, d) })

		// POST /api/mail/password	-> Sends a password reset email
		ma.POST("/reset-passwd", func(c *gin.Context) { mail.PasswdReset(c, d) })
	}

	ff := m.Group("/files")
	{
		// GET /api/files/:id/owns	-> Checks if a user owns a file
		ff.GET("/:id/owns", jwt, func(c *gin.Context) { file.FileOwns(c, d) })

		// GET /api/files/:id		-> Returns a file by it's file_key if the user owns it
		ff.GET("/:id", func(c *gin.Context) { file.Fetch(c, d) })

		// POST /api/files/bulk 	-> Returns a user's files in bulk
		ff.POST("/bulk", jwt, func(c *gin.Context) { file.FetchBulk(c, d) })

		// POST /api/files         	-> Uploads a new file and stores it in the database
		ff.POST("", jwt, bodySizeLimiter, func(c *gin.Context) { file.Upload(c, d) })

		// PATCH /api/files/:id		-> Updates a file
		ff.PATCH("/:id", jwt, func(c *gin.Context) { file.Edit(c, d) })

		// DELETE /api/files/		-> Deletes multiple files
		ff.DELETE("", jwt, func(c *gin.Context) { file.Delete(c, d) })

		// POST /api/files/search	-> Searches for files saved in the database
		ff.POST("/search", jwt, func(c *gin.Context) { file.Search(c, d) })
	}

	f := m.Group("/ffmpeg", jwt)
	{
		// GET /api/ffmpeg/start	-> Starts an FFmpeg job
		f.GET("/start", func(c *gin.Context) { ffmpeg.Start(c, d) })

		// GET /api/ffmpeg/progress	-> Returns the progress of a job
		f.POST("/process", bodySizeLimiter, func(c *gin.Context) { ffmpeg.Process(c, d) })

		// POST /api/ffmpeg/process	-> Processes a file provided in a multipart form
		f.GET("/progress", turnstile, func(c *gin.Context) { ffmpeg.Progress(c, d) })
	}

	// Profiles are split out separately because they contain some specific things
	p := m.Group("/profile")
	{
		p.GET("/:username", func(c *gin.Context) { profile.Fetch(c, d) })
	}

	s3, err := aws.NewS3()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client, %w", err)
	}

	err = redis.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis client, %w", err)
	}

	d.S3 = s3
	d.Uploader = service.NewUploader(d.JobQueue, s3)

	// Start FFmpeg job queue
	d.JobQueue.StartWorkerPool()

	// Check for useless tokens every week because they expire rarely
	go service.StaleTokenCleanup(time.Hour*24*7, db.Gorm)

	// Check for expired accounts rarely because they have a week to verify
	go service.UnverifiedUserCleanup(time.Hour*24*7, db.Gorm)

	return router, nil
}
