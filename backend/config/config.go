// Package config contains code to set the default values and read
// config files to be used throughout the whole application
package config

import (
	"bitwise74/video-api/pkg/util"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var validLogLevels = []string{"debug", "info", "warn", "error"}

const (
	gray  = "\x1b[90m"
	reset = "\x1b[0m"
)

func genSecret() string {
	b := make([]byte, 64)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func checkRunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		panic(err)
	}

	return true
}

func makeLogger() {
	// Configure logger as soon as we know the log level is valid
	atom := zap.NewAtomicLevel()

	switch os.Getenv("APP_LOG_LEVEL") {
	case "debug":
		atom.SetLevel(zap.DebugLevel)
	case "warn":
		atom.SetLevel(zap.WarnLevel)
	case "error":
		atom.SetLevel(zap.ErrorLevel)
	case "fatal":
		atom.SetLevel(zap.FatalLevel)
	default:
		atom.SetLevel(zap.InfoLevel)
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = atom
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(gray + t.Format("15:04:05.000") + reset)
	}
	cfg.EncoderConfig.EncodeCaller = func(ec zapcore.EntryCaller, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(gray + ec.TrimmedPath() + reset)
	}

	cfg.DisableStacktrace = true

	log, _ := cfg.Build()
	zap.ReplaceGlobals(log)
}

// Setup prepares everything config-related so that the app can
// start working. Function will return an error if something
// is critically wrong and the application can't run because of
// that.
func Setup() error {
	if !checkRunningInDocker() {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("failed to read .env file, %w", err)
		}
	}

	// Check if vips is installed
	err := exec.Command("vips").Run()
	if err != nil {
		return errors.New("vips not found")
	}

	if l := os.Getenv("APP_LOG_LEVEL"); l == "" || !slices.Contains(validLogLevels, l) {
		os.Setenv("APP_LOG_LEVEL", "info")
	}

	makeLogger()

	if os.Getenv("HOST_PORT") == "" {
		return errors.New("no port provided")
	}

	if os.Getenv("HOST_DOMAIN") == "" {
		return errors.New("no domain provided")
	}

	if os.Getenv("HOST_CORS") == "" {
		return errors.New("no cors origins provided")
	}

	if os.Getenv("HOST_SSL_ENABLED") == "true" {
		if os.Getenv("HOST_SSL_CERTIFICATE_PATH") == "" {
			return errors.New("no SSL certificate provided")
		}

		if os.Getenv("HOST_SSL_CERTIFICATE_KEY_PATH") == "" {
			return errors.New("no SSL certificate key provided")
		}
	}

	// Test if ffmpeg is present
	err = exec.Command("ffmpeg", "-version").Run()
	if err != nil {
		if path := os.Getenv("FFMPEG_PATH"); path != "" {
			zap.L().Debug("FFmpeg not found in path. Trying user provided path")
			err = exec.Command(path, "-version").Run()
			if err != nil {
				return errors.New("ffmpeg not found")
			}
		} else {
			return errors.New("ffmpeg not found")
		}
	}

	if os.Getenv("FFMPEG_USE_GPU") == "true" {
		gpu, err := util.DetectGPU()
		if err != nil {
			zap.L().Warn("Failed to detect GPU, ffmpeg won't use it to encode/decode", zap.Error(err))
			os.Setenv("FFMPEG_USE_GPU", "false")
		}

		if gpu == "" {
			os.Setenv("FFMPEG_USE_GPU", "false")
			zap.L().Warn("No GPU detected. If it exists ffmpeg won't be able to use it to encode/decode")
		}

		switch gpu {
		case "nvidia":
			os.Setenv("FFMPEG_HWACCEL", "cuda")
			os.Setenv("FFMPEG_ENCODER", "h264_nvenc")
		case "amd":
			os.Setenv("FFMPEG_HWACCEL", "qsv")
			os.Setenv("FFMPEG_ENCODER", "h264_qsv")
		case "intel":
			os.Setenv("FFMPEG_HWACCEL", "vaapi")
			os.Setenv("FFMPEG_ENCODER", "h264_vaapi")
		default:
			zap.L().Warn("Unknown GPU detected, ffmpeg won't use it to encode/decode", zap.String("gpu", gpu))
			os.Setenv("FFMPEG_USE_GPU", "false")
		}

		zap.L().Debug("Detected GPU", zap.String("vendor", gpu))
	}

	if val, err := strconv.Atoi(os.Getenv("FFMPEG_MAX_JOBS")); err != nil {
		return errors.New("FFMPEG_MAX_JOBS is not a valid integer")
	} else if val <= 0 {
		return errors.New("FFMPEG_MAX_JOBS must be set least 1")
	}

	if val, err := strconv.Atoi(os.Getenv("FFMPEG_WORKERS")); err != nil {
		return errors.New("FFMPEG_WORKERS is not a valid integer")
	} else if val <= 0 {
		return errors.New("FFMPEG_WORKERS must be set least 1")
	}

	if os.Getenv("SECURITY_JWT_SECRET") == "" {
		zap.L().Warn("You haven't set a JWT secret, so it has been generated for you. Please set it as an environment variable or in the config.toml file.", zap.String("secret", genSecret()))
		os.Exit(0)
	}

	if val, err := strconv.Atoi(os.Getenv("SECURITY_RATE_LIMIT")); err != nil || val <= 0 {
		os.Setenv("SECURITY_RATE_LIMIT", "15")
	}

	if v := os.Getenv("MAIL_CONFIRMATIONS_ENABLE"); v != "true" {
		zap.L().Warn("Email verifications are disabled. Users won't be able to reset password and attackers will be able to create infinite accounts")
	} else {
		if os.Getenv("MAIL_SENDER_ADDRESS") == "" {
			return errors.New("no sender address provided")
		}

		if os.Getenv("MAIL_HOST") == "" {
			return errors.New("no mail host provided")
		}

		if os.Getenv("MAIL_PORT") == "" {
			return errors.New("no mail host port provided")
		}

		if os.Getenv("MAIL_PASSWORD") == "" {
			return errors.New("no mail host password provided")
		}
	}

	if s := os.Getenv("STORAGE_TYPE"); s == "s3" {
		if os.Getenv("ACCESS_KEY_ID") == "" {
			return errors.New("no access key id provided")
		}

		if os.Getenv("SECRET_ACCESS_KEY") == "" {
			return errors.New("no secret access key provided")
		}

		if os.Getenv("REGION") == "" {
			return errors.New("no region provided")
		}

		if os.Getenv("BUCKET") == "" {
			return errors.New("no bucket provided")
		}
	} else {
		return errors.New("invalid STORAGE_TYPE provided")
	}

	if os.Getenv("TURNSTILE_ENABLE") == "false" {
		zap.L().Warn("Turnstile is disabled. FFmpeg endpoints won't be guarded against bots")
	} else {
		if os.Getenv("TURNSTILE_SECRET_TOKEN") == "" {
			return errors.New("no turnstile secret token provided")
		}
	}

	return nil
}
