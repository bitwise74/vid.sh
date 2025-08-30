package app

import (
	"bitwise74/video-api/app/ffmpeg"
	"bitwise74/video-api/app/file"
	"bitwise74/video-api/app/root"
	"bitwise74/video-api/app/user"
	"bitwise74/video-api/aws"
	"bitwise74/video-api/db"
	"bitwise74/video-api/internal"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/pkg/middleware"
	"bitwise74/video-api/pkg/security"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// TODO: use redis
var store = persist.NewMemoryStore(time.Minute)

func NewRouter() (*gin.Engine, error) {
	d := &internal.Deps{
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
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "TurnstileToken", "Range"},
			ExposeHeaders:    []string{"Content-Length", "Content-Range"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		gin.Recovery(),
		middleware.NewRequestIDMiddleware(),
		gin.Logger(),
	)

	router.HandleMethodNotAllowed = true
	router.RedirectFixedPath = true

	// maxUploadSize, _ := strconv.ParseInt(os.Getenv("UPLOAD_MAX_SIZE"), 10, 64)
	rateLimit, _ := strconv.Atoi(os.Getenv("SECURITY_RATE_LIMIT"))

	jwt := middleware.NewJWTMiddleware(db)
	turnstile := middleware.NewTurnstileMiddleware()
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
	}

	u := m.Group("/users")
	{
		// GET /api/users		-> Returns the basic info of a user
		u.GET("", jwt, func(c *gin.Context) { user.UserFetch(c, d) })

		// POST /api/users 		-> Registers a new user
		u.POST("", func(c *gin.Context) { user.UserRegister(c, d) })

		// POST /api/users/login 	-> Logs in a user and returns a JWT token
		u.POST("/login", func(c *gin.Context) { user.UserLogin(c, d) })

		// POST /api/users/verify	-> Verifies a new user
		u.POST("/verify", func(c *gin.Context) { user.UserVerify(c, d) })

		// DELETE /api/users/:id 	-> Deletes a user account
		// users.DELETE("/:id", jwt)
	}

	ff := m.Group("/files", jwt)
	{
		// GET /api/files/:id/owns	-> Checks if a user owns a file
		ff.GET("/:id/owns", cacheFor(5*60), func(c *gin.Context) { file.FileOwns(c, d) })

		// GET /api/files/:id		-> Returns a file by it's ID if the user owns it
		ff.GET("/:id", func(c *gin.Context) { file.FileFetch(c, d) })

		// GET /api/files/bulk 		-> Returns a user's files in bulk
		ff.GET("/bulk", func(c *gin.Context) { file.FileFetchBulk(c, d) })

		// POST /api/files         	-> Uploads a new file and stores it in the database
		ff.POST("", func(c *gin.Context) { file.FileUpload(c, d) })

		// PATCH /api/files/:id		-> Updates a file
		ff.PATCH("/:id", func(c *gin.Context) { file.FileEdit(c, d) })

		// DELETE /api/files/:id	-> Deletes a file owned by a user
		ff.DELETE("/:id", func(c *gin.Context) { file.FileDelete(c, d) })

		// GET /api/files/search	-> Searches for files saved in the database
		ff.GET("/search", cacheFor(15), func(c *gin.Context) { file.FileSearch(c, d) })
	}

	f := m.Group("/ffmpeg", jwt)
	{
		// GET /api/ffmpeg/start	-> Starts an FFmpeg job
		f.GET("/start", ffmpeg.FFMpegStart)

		// GET /api/ffmpeg/progress	-> Returns the progress of a job
		f.GET("/progress", func(c *gin.Context) { ffmpeg.FFmpegProcess(c, d) })

		// POST /api/ffmpeg/process	-> Processes a file provided in a multipart form
		f.POST("/process", turnstile, ffmpeg.FFMpegStart)
	}

	d.Argon = security.New()
	s3, err := aws.NewS3()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client, %w", err)
	}

	d.S3 = s3
	d.Uploader = service.NewUploader(d.JobQueue, s3)

	// Start FFmpeg job queue
	d.JobQueue.StartWorkerPool()

	// Check for useless tokens every day because they expire rarely
	go service.TokenCleanup(time.Hour*24, db)

	// Check for expired accounts rarely because they have a week to verify
	go service.AccountCleanup(time.Hour*24*7, db, s3.C)

	return router, nil
}

func cacheFor(sec int) gin.HandlerFunc {
	return cache.CacheByRequestURI(store, time.Second*time.Duration(sec))
}
