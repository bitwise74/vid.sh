// Package types contains any shared types that don't fit anywhere else as they could cause import cycles
package types

import (
	"bitwise74/video-api/aws"
	"bitwise74/video-api/db"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/pkg/security"

	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	DB       *db.DB
	Rdb      *redis.Client
	Argon    *security.ArgonHash
	S3       *aws.S3Client
	JobQueue *service.JobQueue
	Uploader *service.Uploader
}
