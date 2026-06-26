// ffmpeg configs and full clip cutting code with the ffmpeg command

package clipping

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Suthar345Piyush/videoclippingpipeline/internal/config"
)

// cutter struct - it will take ClipInput and returns the output

type Cutter struct {
	binaryPath string
	timeout    time.Duration
}

// cutter function with main ffmpeg config, it will return cutter struct

func NewCutter(cfg *config.FFmpegConfig) *Cutter {

	path := cfg.BinaryPath
	if path == "" {
		path = "ffmpeg"
	}

	return &Cutter{
		binaryPath: path,
		timeout:    10 * time.Minute,
	}

}

// a custom 10 minute cutter - timer

func (c *Cutter) WithTimeout(d time.Duration) *Cutter {

	if d <= 0 {
		d = 10 * time.Minute
	}

	return &Cutter{
		binaryPath: c.binaryPath,
		timeout:    d,
	}

}

/* main cutter function this function used to extract the clip out of video

ffmpeg command is going to used will be -

Complete Command :-

  ffmpeg -y -ss <start> -to <end> -i <inputFilePath> -c:v libx264 -preset fast -crf 23 -c:a aac -b:a 128k -movflags +faststart -avoid_negative_ts make_zero <output>

	details about command ---

	-ss -> start time

	-to -> end time

	-i -> info of source file

	-c:v -> video encoder in our case is libx264

	-preset -> fast

	-crf -> constant rate factor, using ffmpeg default 23

	-c:a -> audio encoder in our case is aac

	- b:a -> audio bitrate, using ffmpeg default 128k

	-movflags +faststart -> for front web streaming and online video viewing

	-avoid_negative_ts make_zero -> this will be the option, to avoid negative timestamp, make zero will shift all the timestamps, so video's first timestamp will become zero

*/

func (c *Cutter) Cut(ctx context.Context, ci ClipInput) (Result, error) {

	if err := ci.validate(); err != nil {
		return Result{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)

	defer cancel()

	args := c.buildArgs(ci)

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			return Result{}, fmt.Errorf("clipping: ffmpeg command timed out after %s cutting %q [%.3f, %.3f]", c.timeout, ci.SourcePath, ci.StartTime, ci.EndTime)
		}

		return Result{}, fmt.Errorf("clipping: ffmpeg command failed for %q [%.3f, %.3f]: %w\nstderr: %s", ci.SourcePath, ci.StartTime, ci.EndTime, err, strings.TrimSpace(stderr.String()))

	}

	return Result{
		OutputPath: ci.OutputPath,
		Duration:   ci.EndTime - ci.StartTime,
	}, nil

}

// seperate function for ffmpeg command with clipInput, and it returns slice of string

func (c *Cutter) buildArgs(ci ClipInput) []string {

	start := strconv.FormatFloat(ci.StartTime, 'f', 3, 64)
	end := strconv.FormatFloat(ci.EndTime, 'f', 3, 64)

	args := []string{
		"-y",
		"-ss", start,
		"-to", end,
		"-i", ci.SourcePath,
		"-c:v", ci.Preset.VideoCodec,
		"-preset", ci.Preset.Preset,
		"-crf", strconv.Itoa(ci.Preset.CRF),
		"-c:a", ci.Preset.AudioCodec,
		"-b:a", ci.Preset.AudioBitrate,
		"-movflags", "+faststart",
		"-avoid_negative_ts", "make_zero",
	}

	// scale width used for scaling video, rotating , crop and all
	// appending the -vf - video filter flag with -2 value

	if ci.Preset.ScaleWidth > 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=%d:-2", ci.Preset.ScaleWidth))
	}

	args = append(args, ci.OutputPath)

	return args

}
