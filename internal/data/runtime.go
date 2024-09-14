package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%d mins", r)
	s = strconv.Quote(s)
	return []byte(s), nil
}

func (r *Runtime) UnmarshalJSON(b []byte) error {
	value, err := strconv.Unquote(string(b))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(value, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	n, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(n)
	return nil
}
