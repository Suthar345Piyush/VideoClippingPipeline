package config

// main struct use in load function

type Config struct {
	FFmpegPath   string `mapstructure:"ffmpeg_path"`
	FFprobePath  string `mapstructure:"ffprobe_path"`
	StoragePath  string `mapstructure:"storage_path"`
	ClipDuration int    `mapstructure:"clip_duration"`
	DatabasePath string `mapstructure:"database_path"`
}

// Load function

func Load() (*Config, error) {

}
