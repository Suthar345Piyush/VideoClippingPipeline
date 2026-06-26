// ffmpeg configs and full clip cutting code with the ffmpeg command

package clipping

import (
	"context"
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

func (c *Cutter) Cut(ctx context.Context, in ClipInput) (Result, error) {

}
