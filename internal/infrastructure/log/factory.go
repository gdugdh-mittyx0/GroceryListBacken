package log

import (
	"errors"
	"glbackend/internal/adapters/logging"
)

// logging.Logger
var (
	errInvalidLoggerInstance = errors.New("invalid log instance")
)

const (
	InstanceLogrusLogger int = iota
)

func NewLoggerFactory(instance int) (logging.Logger, error) {
	switch instance {
	case InstanceLogrusLogger:
		return NewLogrusLogger(), nil
	default:
		return nil, errInvalidLoggerInstance
	}
}
