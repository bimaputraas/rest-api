package pkgerrors

import (
	"errors"
)

const (
	ErrInvalidArgument Type = 1
	ErrNotFound        Type = 2
	ErrInternal        Type = 3
	ErrIllegal         Type = 4
	ErrUnimplemented   Type = 5
)

type (
	Type  int
	Error interface {
		Code() Type
		error
	}

	errImpl struct {
		code Type
		msg  string
	}
)

func ParseError(err error) (Error, bool) {
	var r *errImpl
	ok := errors.As(err, &r)
	if !ok {
		return &errImpl{}, false
	}

	return r, true
}

func Code(err error) (Type, bool) {
	parsed, ok := ParseError(err)
	if !ok {
		return 0, false
	}
	return parsed.Code(), true
}

func InvalidArgument(err error) error {
	return &errImpl{code: ErrInvalidArgument, msg: err.Error()}
}

func NotFound(err error) error {
	return &errImpl{code: ErrNotFound, msg: err.Error()}
}

func Internal(err error) error {
	return &errImpl{code: ErrInternal, msg: err.Error()}
}

func Illegal(err error) error {
	return &errImpl{code: ErrIllegal, msg: err.Error()}
}

func Unimplemented() error {
	return &errImpl{code: ErrUnimplemented, msg: "unimplemented"}
}

func (e *errImpl) Code() Type {
	return e.code
}
func (e *errImpl) Error() string {
	return e.msg
}
