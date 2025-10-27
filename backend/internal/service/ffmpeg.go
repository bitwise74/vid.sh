package service

import (
	"bitwise74/video-api/pkg/util"
	"bitwise74/video-api/pkg/validators"
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"

	"go.uber.org/zap"
)

type FFmpegJob struct {
	ID       string
	UserID   string
	FilePath string
	Output   io.Writer
	UseGPU   bool
	Opts     *validators.ProcessingOptions
	Args     *[]string
	Ctx      context.Context
	Done     chan error
}

type FFMpegJobStats struct {
	Progress float64
	JobID    string
	State    string
	Stopped  bool
}

type JobQueue struct {
	jobs    chan *FFmpegJob
	running atomic.Int32
	workers int64
}

var (
	ErrJobQueueFull = errors.New("job queue is full")
	ErrJobExists    = errors.New("a job is already running for this user")
)

// NewJobQueue initializes a new job queue that limits the
// max amount of jobs that can be queued at once
func NewJobQueue() *JobQueue {
	maxJobs, _ := strconv.ParseInt(os.Getenv("FFMPEG_MAX_JOBS"), 10, 32)
	workers, _ := strconv.ParseInt(os.Getenv("FFMPEG_WORKERS"), 10, 32)

	zap.L().Debug("Initializing job queue", zap.Int64("max_jobs", maxJobs))

	return &JobQueue{
		jobs:    make(chan *FFmpegJob, maxJobs),
		workers: workers,
	}
}

func (q *JobQueue) StartWorkerPool() {
	for range q.workers {
		go q.worker()
	}
}

func (q *JobQueue) worker() {
	for job := range q.jobs {
		err := q.runFFmpegJob(job)

		job.Done <- err
		close(job.Done)

		q.running.Add(-1)
		ProgressMap.Delete(job.UserID)

		if err != nil {
			zap.L().Error("FFmpeg job finished with an error",
				zap.String("user_id", job.UserID),
				zap.String("job_id", job.ID),
				zap.Error(err))
		} else {
			zap.L().Debug("FFmpeg job finished successfully")
		}
	}
}

func (q *JobQueue) Enqueue(job *FFmpegJob) error {
	if q.running.Load() >= int32(cap(q.jobs)) {
		return errors.New("job queue full")
	}

	ProgressMap.Store(job.UserID, FFMpegJobStats{
		Progress: 0.0,
		JobID:    job.ID,
		State:    "Processing video...",
		Stopped:  false,
	})

	q.jobs <- job
	q.running.Add(1)
	zap.L().Debug("New ffmpeg job enqueued", zap.Int32("enqueued", q.running.Load()), zap.String("user_id", job.UserID))
	return nil
}

func (q *JobQueue) MakeFFmpegFlags(opts *validators.ProcessingOptions, p string) ([]string, float64, error) {
	args := []string{}

	encoder := os.Getenv("FFMPEG_ENCODER")
	if encoder == "" {
		encoder = "libx264"
	}

	var duration float64
	var err error

	args = append(args, "-i", p)

	if opts.TrimStart > 0 {
		args = append(args, "-ss", util.FloatToTimestamp(opts.TrimStart))
	}

	if opts.TrimEnd > 0 && opts.TrimStart >= 0 {
		args = append(args, "-to", util.FloatToTimestamp(opts.TrimEnd))
		duration = opts.TrimEnd - opts.TrimStart
	} else {
		duration, err = GetDuration(p)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to run ffprobe to determine video duration: %w", err)
		}
	}

	args = append(args, "-c:v", encoder)

	if opts.ShouldCrop {
		cropStr := fmt.Sprintf("crop=%d:%d:%d:%d", opts.CropW, opts.CropH, opts.CropX, opts.CropY)
		args = append(args, "-vf", cropStr)
	}

	if opts.LosslessExport {
		switch encoder {
		case "libx264":
			args = append(args, "preset", "slow", "-crf", "18", "-pix_fmt", "yuv420p")
		case "h264_nvenc", "hevc_nvenc":
			args = append(args, "-preset", "p7", "-rc", "vbr", "-cq", "19", "-b:v", "0")
		default:
			args = append(args, "-crf", "10")
		}
	} else if opts.TargetSize > 0 {
		totalKilobits := opts.TargetSize * 8388.608
		totalBitrateKbps := totalKilobits / duration
		videoBitrateKbps := totalBitrateKbps - 128
		if videoBitrateKbps <= 0 {
			videoBitrateKbps = 5
		}
		videoBitrateStr := fmt.Sprintf("%.0fK", videoBitrateKbps)
		bufSizeStr := fmt.Sprintf("%dk", int(videoBitrateKbps*2))
		args = append(args,
			"-b:v", videoBitrateStr,
			"-maxrate", videoBitrateStr,
			"-bufsize", bufSizeStr,
		)
	}

	args = append(args,
		"-c:a", "copy",
		"-movflags", "+frag_keyframe+empty_moov+faststart",
		"-loglevel", "error",
	)

	args = append(args,
		"-f", "mp4",
		"pipe:1",
		"-progress", "pipe:2",
		"-nostats",
	)

	return args, duration, nil
}

// The encoder is always appended
func addHWAccelFlags(args []string) []string {
	useGPU, _ := strconv.ParseBool(os.Getenv("FFMPEG_USE_GPU"))
	hwaccel := os.Getenv("FFMPEG_HWACCEL")

	if useGPU && hwaccel != "" {
		for i, arg := range args {
			if arg == "-i" && i > 0 {
				args = append(args[:i], append([]string{"-hwaccel", hwaccel}, args[i:]...)...)
				break
			}
		}
	}

	return args
}

func (q *JobQueue) runFFmpegJob(job *FFmpegJob) error {
	var duration float64
	var err error

	if job.Args == nil {
		if job.Opts == nil {
			return errors.New("no arguments provided")
		}

		var args []string

		args, duration, err = q.MakeFFmpegFlags(job.Opts, job.FilePath)
		if err != nil {
			return err
		}

		job.Args = &args
	}

	if job.UseGPU {
		*job.Args = addHWAccelFlags(*job.Args)
	}

	cmd := exec.CommandContext(job.Ctx, "ffmpeg", *job.Args...)

	zap.L().Debug("Running FFmpeg command", zap.String("cmd", cmd.String()))

	stderrPipe, _ := cmd.StderrPipe()
	defer stderrPipe.Close()

	stderrBuf := &bytes.Buffer{}

	go func() {
		scanner := bufio.NewScanner(io.TeeReader(stderrPipe, stderrBuf))
		for scanner.Scan() {
			line := scanner.Text()

			if line == "progress=end" {
				break
			}

			if after, ok := strings.CutPrefix(line, "out_time_ms="); ok {
				msStr := after
				outTimeMs, err := strconv.ParseFloat(msStr, 64)
				if err == nil {
					percent := (outTimeMs / (duration * 1000)) / 10
					ProgressMap.Store(job.UserID, FFMpegJobStats{
						JobID:    job.ID,
						Progress: percent,
						State:    "Processing video...",
						Stopped:  false,
					})
				}
			}
		}

		ProgressMap.Store(job.UserID, FFMpegJobStats{
			JobID:    job.ID,
			Progress: 100.0,
			State:    "Finishing...",
			Stopped:  false,
		})
	}()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg, %w", err)
	}

	_, err = io.Copy(job.Output, stdout)
	if err != nil {
		return fmt.Errorf("streaming error, %w", err)
	}

	if err := cmd.Wait(); err != nil {
		zap.L().Error("FFmpeg failed", zap.Error(err), zap.String("stderr", stderrBuf.String()))
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return nil
}
