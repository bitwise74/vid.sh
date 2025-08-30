package internal

import (
	"bitwise74/video-api/aws"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/pkg/security"

	"gorm.io/gorm"
)

type Deps struct {
	DB       *gorm.DB
	Argon    *security.ArgonHash
	S3       *aws.S3Client
	JobQueue *service.JobQueue
	Uploader *service.Uploader
}
