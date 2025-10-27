// Package types contains any shared types that don't fit anywhere else as they could cause import cycles
package types

import (
	"bitwise74/video-api/aws"
	"bitwise74/video-api/db"
	"bitwise74/video-api/internal/service"
)

type Dependencies struct {
	DB       *db.DB
	S3       *aws.S3Client
	JobQueue *service.JobQueue
	Uploader *service.Uploader
}
