package logging

type Error struct {
	log        Logger
	err        error
	key        string
	httpStatus int
}

func NewError(log Logger, err error, key string, httpStatus int) Error {
	return Error{
		log:        log,
		err:        err,
		key:        key,
		httpStatus: httpStatus,
	}
}

func (e Error) Log(msg string) {
	e.log.WithFields(Fields{
		"key":         e.key,
		"error":       e.err.Error(),
		"http_status": e.httpStatus,
	}).Errorf(msg)
}
