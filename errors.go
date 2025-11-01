package salamoonder

import "errors"

var (
	ErrTaskNotReady = errors.New("task not ready")
	ErrNoApiKey     = errors.New("no api key")
	ErrUnsupportedTaskOptionsType = errors.New("unsupported task options type")
)