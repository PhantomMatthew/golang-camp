package log

import "go.uber.org/zap"

// Options config
type Options struct {
	zap.Config
	EnableKafka bool `json:"enable_kafka"`
	LogFileDir    string `json:"log_file_dir"`
	AppName       string `json:"app_name"`
	FatalFileName string `json:"fatal_file_name"`
	ErrorFileName string `json:"error_file_name"`
	WarnFileName  string `json:"warn_file_name"`
	InfoFileName  string `json:"info_file_name"`
	DebugFileName string `json:"debug_file_name"`
	MaxSize       int    `json:"max_size"` // MB
	MaxBackups    int    `json:"max_backups"`
	MaxAge        int    `json:"max_age"` // days
}