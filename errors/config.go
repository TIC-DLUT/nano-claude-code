package errors

import "errors"

var (
	ReadInConfigError    = errors.New("failed to get config from the config file")
	ConfigFileWriteError = errors.New("failed to write config file")
	ConfigDirCreateError = errors.New("failed to create config directory")
)
