package runner

import "errors"

var (
	ErrServiceAlreadyExists = errors.New("service already exists")
	ErrRunnerAlreadyRunning = errors.New("runner already running")
)
