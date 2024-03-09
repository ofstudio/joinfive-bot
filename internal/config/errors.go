package config

import "errors"

var (
	ErrEnvParse = errors.New("config: failed to parse environment variables")
)
