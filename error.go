package glcm

import "errors"

var (
	ErrServiceNotRunning = errors.New("service not running")
)

var (
	ErrRegisterServiceAlreadyExists = errors.New("service already exists")
	ErrRunnerAlreadyRunning         = errors.New("runner already running")
	ErrRegisterNilService           = errors.New("can not register nil service")
	ErrUnsupportedOS                = errors.New("unsupported OS")
	ErrSocketNoService              = errors.New("no service provided")
)
