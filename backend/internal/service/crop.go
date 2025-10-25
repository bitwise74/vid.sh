package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Shorhand
func i2S(i int) string {
	return strconv.Itoa(i)
}

// Crop runs vips to crop a provided image. Returns the path to the processed image.
// Callers must remember to remove the file after using it. A path should always be
// returned and if it's not that means the file creation has failed.
// The cropping options should be checked by the caller
func Crop(path string, opts []int) (string, error) {
	if len(opts) != 4 {
		return "", fmt.Errorf("invalid number of cropping options provided")
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*15))
	defer cancel()

	tmpFile, err := os.CreateTemp("", "profile-pic-*.webp")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file for vips, %w", err)
	}
	defer tmpFile.Close()

	cmd := exec.CommandContext(ctx, "vips", "crop", path, tmpFile.Name(), i2S(opts[0]), i2S(opts[1]), i2S(opts[2]), i2S(opts[3]))

	if err := cmd.Start(); err != nil {
		return tmpFile.Name(), fmt.Errorf("failed to start vips, %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return tmpFile.Name(), fmt.Errorf("vips failed, %w", err)
	}

	return tmpFile.Name(), nil
}
