// models (structs) for clipping

package clipping

import "fmt"

// struct for encode preset - which have some ffmpeg parameters for output clip

type EncodePreset struct {

	// video codec - using libx264 encoder

	VideoCodec string

	// preset are basically encoding settings - there are in total three preset (fast, balanced, slow)

	Preset string

	// constant rate factor (crf) - it is rate control mode for encoders like libx264, it's default value is 23

	CRF int

	// audio codec - just like we have video encoder (libx264), same we have for audio as well (aac)

	AudioCodec string

	// audio bitrate - it means data processed or transmitted per unit time - kbps, the default value will be 128

	AudioBitrate string

	// scale width - means scaling video horizontally , 0 means no scaling , > 0 means scaling is done

	ScaleWidth int
}

// some default values to set for the encode preset

var DefaultPreset = EncodePreset{
	VideoCodec:   "libx264",
	Preset:       "fast",
	CRF:          23,
	AudioCodec:   "aac",
	AudioBitrate: "128k",
}

// we have provide input to cutter function which cuts clips out the full length video

// input items are like

type ClipInput struct {
	SourcePath string // path of the actual source video

	OutputPath string // output path where clip will be set

	// start time and end time of the clip both in float64 according to db values

	StartTime float64
	EndTime   float64

	// preset that we set

	Preset EncodePreset
}

// validate function to check for default's and incoming values

func (ci *ClipInput) validate() error {

	if ci.SourcePath == "" {
		return fmt.Errorf("clipping: SourcePath must not be empty")
	}

	if ci.OutputPath == "" {
		return fmt.Errorf("clipping: OutputPath must not be empty")
	}

	if ci.StartTime < 0 {
		return fmt.Errorf("clipping: StartTime must be >= 0", ci.StartTime)
	}

	if ci.EndTime <= ci.StartTime {
		return fmt.Errorf("clipping: EndTime must be greater than EndTime", ci.EndTime, ci.StartTime)
	}

	if ci.Preset.VideoCodec == "" {
		ci.Preset = DefaultPreset
	}

	return nil

}

// final output returns as clip with the path and duration of clip itself

type Result struct {
	OutputPath string
	Duration   float64 // duration will be  = endtime - startime
}
