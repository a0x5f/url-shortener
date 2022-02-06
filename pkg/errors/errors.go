package errors

import (
	"errors"
)

func Wrap(err error, msg string) error {
	return errors.New(msg + " : " + err.Error())
}
