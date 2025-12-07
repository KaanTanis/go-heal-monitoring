package types

import "time"

type TargetConfig struct {
	URL      string        `yaml:"url"`
	Name     string        `yaml:"name"`
	Interval time.Duration `yaml:"interval"`
}

type Config struct {
	Targets []TargetConfig `yaml:"targets"`
}

type LogEntry struct {
	TimeStamp  string `json:"timestamp"`
	Level      string `json:"level"`
	TargetName string `json:"target_name"`
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	DurationMs int64  `json:"duration_ms"`
	Error      string `json:"error,omitempty"`
}