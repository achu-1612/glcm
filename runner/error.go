package runner

import "errors"

var (
	ErrRegisterServiceAlreadyExists = errors.New("service already exists")
	ErrRunnerAlreadyRunning         = errors.New("runner already running")
	ErrRegisterNilService           = errors.New("can not register nil service")
)
