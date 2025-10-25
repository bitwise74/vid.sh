package service

import (
	a "bitwise74/video-api/aws"
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/pkg/util"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

const minMultipartSize = 12 << 20

type Uploader struct {
	S3       *a.S3Client
	JobQueue *JobQueue
}

func NewUploader(j *JobQueue, s *a.S3Client) *Uploader {
	return &Uploader{
		JobQueue: j,
		S3:       s,
	}
}

// Do should be used with a file that's ready for upload and was checked. It creates a thumbnail for the video file and uploads both files. Providing an override value will instead update an existing file. Files are deleted after upload
func (u *Uploader) Do(p, name, userID string, override ...string) (*model.File, error) {
	videoFile, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("failed to open video file, %w", err)
	}
	defer os.Remove(p)
	defer videoFile.Close()

	videoStat, _ := videoFile.Stat()

	thumbPath, err := MakeThumbnail(p, userID, u.JobQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3, %w", err)
	}

	thumbFile, err := os.Open(thumbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open thumbnail file, %w", err)
	}
	defer os.Remove(thumbPath)
	defer thumbFile.Close()

	// Prepare things for background operations
	var wg sync.WaitGroup
	wg.Add(3)

	key := util.RandStr(10)

	errors := make(chan error, 3)
	uploadedKeys := []string{}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	defer cancel()

	if len(override) > 0 {
		key = override[0]
	}

	// Thumbnail upload
	go func() {
		defer wg.Done()
		zap.L().Debug("Starting upload_thumbnail subprocess")

		_, err := u.S3.C.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       u.S3.Bucket,
			Key:          aws.String(key + ".webp"),
			Body:         thumbFile,
			CacheControl: aws.String("public, max-age=86400, stale-while-revalidate=3600"),
			ContentType:  aws.String("image/webp"),
		})
		if err != nil {
			errors <- fmt.Errorf("failed to upload thumbnail, %w", err)
			return
		}

		uploadedKeys = append(uploadedKeys, key+".webp")
		errors <- nil
	}()

	// Video upload
	go func() {
		defer wg.Done()
		zap.L().Debug("Starting upload_video subprocess")

		var uploader *manager.Uploader
		if videoStat.Size() > minMultipartSize {
			uploader = manager.NewUploader(u.S3.C, func(u *manager.Uploader) {
				u.Concurrency = 5
				u.PartSize = 6 << 20
			})
		}

		objectInput := &s3.PutObjectInput{
			Bucket:        u.S3.Bucket,
			Key:           aws.String(key + ".mp4"),
			Body:          videoFile,
			ContentLength: aws.Int64(videoStat.Size()),
			ContentType:   aws.String("video/mp4"),
			CacheControl:  aws.String("public, max-age=3600, stale-while-revalidate=60"),
		}

		var err error
		if uploader != nil {
			_, err = uploader.Upload(ctx, objectInput)
		} else {
			_, err = u.S3.C.PutObject(ctx, objectInput)
		}
		if err != nil {
			errors <- fmt.Errorf("failed to upload video to s3, %w", err)
			return
		}

		uploadedKeys = append(uploadedKeys, key+".mp4")
		errors <- nil
	}()

	var duration float64

	go func() {
		defer wg.Done()
		zap.L().Debug("Starting ffprobe_duration subprocess")
		var err error

		duration, err = GetDuration(p)
		if err != nil {
			errors <- fmt.Errorf("failed to get video duration: %w", err)
			return
		}

		errors <- nil
	}()

	for range 3 {
		if err := <-errors; err != nil {
			cancel()

			for _, id := range uploadedKeys {
				_, err := u.S3.C.DeleteObject(context.Background(), &s3.DeleteObjectInput{
					Bucket: u.S3.Bucket,
					Key:    aws.String(id),
				})
				if err != nil {
					zap.L().Error("Failed to cleanup after faile uploads", zap.String("id", id), zap.Error(err))
				} else {
					zap.L().Debug("Cleaned up after failed upload", zap.String("id", id))
				}
			}

			return nil, err
		}
	}

	wg.Wait()

	fileEnt := &model.File{
		UserID:       userID,
		FileKey:      key + ".mp4",
		OriginalName: name,
		Format:       "video/mp4",
		Size:         videoStat.Size(),
		Tags:         []string{},
		State:        "ready",
		Version:      1,
		Duration:     duration,
		CreatedAt:    time.Now().Unix(),
	}

	return fileEnt, nil
}
