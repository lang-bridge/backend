package httperr

import (
	"errors"
)

type badRequest struct {
	err error
}

func (b badRequest) Error() string {
	return b.err.Error()
}

func (b badRequest) Unwrap() error {
	return b.err
}

func BadRequest(err error) error {
	return &badRequest{err: err}
}

func IsBadRequest(err error) bool {
	target := new(badRequest)
	return errors.As(err, &target)
}
