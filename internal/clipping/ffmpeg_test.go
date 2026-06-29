// tests for ffmpeg

package clipping

import (
	"strings"
	"testing"

	"github.com/Suthar345Piyush/videoclippingpipeline/internal/config"
)

// tests for build args command

func TestBuildArgs_BasicFlags(t *testing.T) {

	c := NewCutter(&config.FFmpegConfig{BinaryPath: "ffmpeg"})

	ci := ClipInput{
		SourcePath: "/src/video.mp4",
		OutputPath: "/out/clip.mp4",
		StartTime:  30.0,
		EndTime:    60.0,
		Preset:     DefaultPreset,
	}

	args := c.buildArgs(ci)

	joined := strings.Join(args, " ")

	// starting time -ss must before -i (source file path)

	ssIdx := indexOfArgs(args, "-ss")
	iIdx := indexOfArgs(args, "-i")

	if ssIdx == -1 {
		t.Fatal("-ss flag not found in args")
	}

	if iIdx == -1 {
		t.Fatal("-i flag not found in args")
	}

	// -ss must appear before -i in main command

	if ssIdx > iIdx {
		t.Errorf("-ss (Index %d) must appear before -i (Index %d) for input seek", ssIdx, iIdx)
	}

	// verifying every value of ffmpeg command is present, testing with some values

	for _, want := range []string{
		"-y",
		"-ss", "30.000",
		"-to", "60.000",
		"-i", "/src/video.mp4",
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		"-movflags", "+faststart",
		"-avoid_negative_ts", "make_zero",
		"/out/clip.mp4",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("expected arg %q in args: %s", want, joined)
		}
	}
}

// seperate test on build args (ffmpeg command), for -ss and -to both should be present before -i, so that ffmpeg will not decode the whole which is slow operation and we don't want our service to be that much slower

func TestBuildArgs_StartEndTimeBeforeInput(t *testing.T) {

	c := NewCutter(&config.FFmpegConfig{})
	args := c.buildArgs(ClipInput{
		SourcePath: "/src/v.mp4",
		OutputPath: "/out/c.mp4",
		StartTime:  120.5,
		EndTime:    180.0,
		Preset:     DefaultPreset,
	})

	ssIdx := indexOfArgs(args, "-ss")
	toIdx := indexOfArgs(args, "-to")
	iIdx := indexOfArgs(args, "-i")

	if ssIdx == -1 || toIdx == -1 || iIdx == -1 {
		t.Fatalf("missing flag: -ss=%d -to=%d -i=%d", ssIdx, toIdx, iIdx)
	}

	if ssIdx > iIdx {
		t.Fatalf("-ss must come before -i, -ss=%d -i=%d", ssIdx, iIdx)
	}

	if toIdx > iIdx {
		t.Fatalf("-to must come before -i, -to=%d -i=%d", toIdx, iIdx)
	}

}

// index of args  function

func indexOfArgs(slice []string, target string) int {
	for i, s := range slice {
		if s == target {
			return i
		}
	}
	return -1
}
