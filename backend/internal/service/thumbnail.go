// Package service contains stuff related to the background processing
// of the application
package service

import (
	"bitwise74/video-api/pkg/util"
	"context"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
)

// MakeThumbnail creates a thumbnail from a multipart.File
func MakeThumbnail(input, userID string, j *JobQueue) (p string, err error) {
	zap.L().Debug("Creating thumbnail for video")

	done := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	thumbPath := path.Join(os.TempDir(), util.RandStr(10)+".webp")
	zap.L().Debug("Writing thumbnail file", zap.String("path", thumbPath))

	err = j.Enqueue(&FFmpegJob{
		ID:     util.RandStr(5),
		UserID: userID,
		Args: &[]string{
			"-loglevel", "error",
			"-ss", "0",
			"-i", input,
			"-frames:v", "1",
			"-vf", "scale=1280:-1",
			"-q:v", "1",
			"-compression_level", "4",
			thumbPath,
		},
		Done: done,
		Ctx:  ctx,
	})
	if err != nil {
		return "", err
	}

	select {
	case err := <-done:
		if err != nil {
			return "", err
		}
	case <-ctx.Done():
		return "", ctx.Err()
	}

	return thumbPath, nil
}
