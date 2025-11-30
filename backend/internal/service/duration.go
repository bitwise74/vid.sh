package service

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func GetDuration(p string) (d float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	zap.L().Debug("Running FFprobe to determine video duration")

	cmd := exec.CommandContext(ctx, "ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-i", p)

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("ffprobe failed, %w (%s)", err, stdErr.String())
	}

	durStr := strings.TrimSpace(stdOut.String())
	d, err = strconv.ParseFloat(strings.TrimSpace(durStr), 64)
	if err != nil {
		return 0, fmt.Errorf("malformed duration: %w (%s)", err, stdErr.String())
	}

	zap.L().Debug("FFprobe finished")
	return d, nil
}
