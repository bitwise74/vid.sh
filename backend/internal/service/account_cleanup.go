package service

import (
	"bitwise74/video-api/internal/model"
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type partialUser struct {
	ID int
}

// AccountCleanup automatically deletes accounts that were registered before mail
// verification and didn't verify after 30 days from the update.
// Also deletes accounts that should have verified their account
// after registration but didn't
func AccountCleanup(t time.Duration, db *gorm.DB, c *s3.Client) {
	ticker := time.NewTicker(t)

	zap.L().Debug("Account cleanup attached", zap.Duration("tick_every", t))

	go func() {
		for range ticker.C {
			var toCleanUserIds []string

			err := db.
				Model(model.User{}).
				Where("expires_at < ?", time.Now()).
				Select("id").
				Find(&toCleanUserIds).
				Error
			if err != nil {
				zap.L().Error("Failed to query db for users to clean", zap.Error(err))
				continue
			}

			if len(toCleanUserIds) == 0 {
				continue
			}

			// If we have any users to delete also get their files to delete from S3
			var toCleanFileKeys []string

			err = db.
				Model(model.File{}).
				Where("user_id LIKE ?", toCleanUserIds).
				Select("file_key", "thumb_key").
				Find(&toCleanFileKeys).
				Error
			if err != nil {
				zap.L().Error("Failed to query db for files to delete", zap.Error(err))
				continue
			}

			// Try deleting files from S3
			if len(toCleanFileKeys) > 0 {
				// S3 can delete at most 1000 files in one batch request.
				// We have to split the slice in len = 1000 parts
				for start := 0; start < len(toCleanFileKeys); start += 1000 {
					end := min(start+1000, len(toCleanFileKeys))

					objects := make([]types.ObjectIdentifier, end-start)
					for i, key := range toCleanFileKeys[start:end] {
						objects[i] = types.ObjectIdentifier{Key: &key}
					}

					if _, err := c.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
						Bucket: aws.String(os.Getenv("BUCKET")),
						Delete: &types.Delete{
							Objects: objects,
						},
					}); err != nil {
						zap.L().Error("Failed to delete files from S3", zap.Error(err))
					}
				}
			}

			// Delete users now
			err = db.
				Where("user_id LIKE ?", toCleanUserIds).
				Delete(model.User{}).
				Error
			if err != nil {
				zap.L().Error("Failed to delete users from database", zap.Error(err))
			}

			zap.L().Debug("Account cleanup finished")
		}
	}()
}
