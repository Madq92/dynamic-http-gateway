package common

type BaseError struct {
	Code     int
	Message  string
	Internal error
}

func (baseError *BaseError) Error() string {
	return baseError.Message
}
