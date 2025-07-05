package rdparser

import (
	"errors"
	"fmt"
)

type Error struct {
	base error
	msg  string
}

func NewError(base error, msg string) error {
	return &Error{base: base, msg: msg}
}

func (err Error) Error() string {
	return fmt.Sprintf("%s - %s", err.base, err.msg)
}

func (err Error) Unwrap() error {
	return err.base
}

func Catch(kind error, addr *error) {
	if v := recover(); v != nil {
		if err, ok := v.(error); ok && errors.Is(err, kind) {
			*addr = err
			return
		}
		panic(v)
	}
}
