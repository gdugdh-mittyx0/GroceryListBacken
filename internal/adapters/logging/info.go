package logging

type Info struct {
	log        Logger
	key        string
	httpStatus int
}

func NewInfo(log Logger, key string, httpStatus int) Info {
	return Info{
		log:        log,
		key:        key,
		httpStatus: httpStatus,
	}
}

func (i Info) Log(msg string) {
	i.log.WithFields(Fields{
		"key":         i.key,
		"http_status": i.httpStatus,
	}).Infof(msg)
}
