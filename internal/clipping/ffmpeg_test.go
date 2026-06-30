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

// tests for start and end time formatting

func TestBuildArgs_TimeFormatting(t *testing.T) {

	c := NewCutter(&config.FFmpegConfig{})

	args := c.buildArgs(ClipInput{
		SourcePath: "/src/input.mp4",
		OutputPath: "/out/output.mp4",
		StartTime:  1.5,
		EndTime:    90.123,
		Preset:     DefaultPreset,
	})

	joined := strings.Join(args, " ")

	if !strings.Contains(joined, "1.500") {
		t.Errorf("start time should be '1.500', args: %s", joined)
	}

	if !strings.Contains(joined, "90.123") {
		t.Errorf("end time should be '90.123', args: %s", joined)
	}

}

// tests for the fixed scale width

func TestBuildArgs_ScaleWidth(t *testing.T) {

	c := NewCutter(&config.FFmpegConfig{})

	preset := DefaultPreset
	preset.ScaleWidth = 1280

	args := c.buildArgs(ClipInput{
		SourcePath: "/src/input.mp4",
		OutputPath: "/out/output.mp4",
		StartTime:  0,
		EndTime:    10,
		Preset:     preset,
	})

	joined := strings.Join(args, " ")

	if !strings.Contains(joined, "-vf") {
		t.Errorf("expected -vf flag when scale width > 0")
	}

	if !strings.Contains(joined, "scale=1280:-2") {
		t.Errorf("expected 'scale=1280:-2' in args, got: %s", joined)
	}

}

// when scale width = 0, then no -vf flag should appear

func TestBuildArgs_NoScaleWidth(t *testing.T) {

	c := NewCutter(&config.FFmpegConfig{})
	args := c.buildArgs(ClipInput{
		SourcePath: "/src/input.mp4",
		OutputPath: "/out/output.mp4",
		StartTime:  0,
		EndTime:    10,
		Preset:     DefaultPreset,
	})

	for _, a := range args {
		if a == "-vf" {
			t.Errorf("unexpected -vf flag when scale width = 0")
		}
	}
}

// their should be a output path to write to it final output ffmpeg will require it

func TestBuildArgs_OutputIsLastArgs(t *testing.T) {

	c := NewCutter(&config.FFmpegConfig{})
	args := c.buildArgs(ClipInput{
		SourcePath: "/src/input.mp4",
		OutputPath: "/out/clip.mp4",
		StartTime:  0,
		EndTime:    10,
		Preset:     DefaultPreset,
	})

	if args[len(args)-1] != "/out/clip.mp4" {
		t.Errorf("last arg must be output path, but got %q", args[len(args)-1])
	}
}

// last build args tests, will be hardcoded codec and other EncodePreset struct

func TestBuildArgs_CustomCodecAndCRF(t *testing.T) {
	c := NewCutter(&config.FFmpegConfig{})

	preset := EncodePreset{
		VideoCodec:   "libx264",
		Preset:       "slow",
		CRF:          18,
		AudioCodec:   "copy",
		AudioBitrate: "192k",
	}

	args := c.buildArgs(ClipInput{
		SourcePath: "/src/input.mp4",
		OutputPath: "/out/output.mp4",
		StartTime:  0,
		EndTime:    10,
		Preset:     preset,
	})

	joined := strings.Join(args, " ")

	for _, want := range []string{"libx246", "slow", "18", "copy", "192k"} {
		if !strings.Contains(joined, want) {
			t.Errorf("expected %q in args: %s", want, joined)
		}
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
