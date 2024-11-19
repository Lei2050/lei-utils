package io

import (
	"io"
	"net"
	"runtime"

	"github.com/pkg/errors"
)

type timeoutError interface {
	Timeout() bool // Is it a timeout error
}

type temperaryError interface {
	Temporary() bool
}

// IsTimeout checks if the error is a timeout error
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}

	err = errors.Cause(err)
	ne, ok := err.(timeoutError)
	return ok && ne.Timeout()
}

// IsTimeout checks if the error is a timeout error
func IsTemporary(err error) bool {
	if err == nil {
		return false
	}

	err = errors.Cause(err)
	ne, ok := err.(temperaryError)
	return ok && ne.Temporary()
}

func WriteFull(conn io.Writer, data []byte) error {
	left := len(data)
	for left > 0 {
		n, err := conn.Write(data)
		if n == left && err == nil {
			return nil
		}

		if n > 0 {
			data = data[n:]
			left -= n
		}

		if err != nil {
			if !IsTemporary(err) {
				return err
			} else {
				runtime.Gosched()
			}
		}
	}
	return nil
}

func ReadFull(reader io.Reader, data []byte) error {
	left := len(data)
	for left > 0 {
		n, err := reader.Read(data)
		if n == left && err == nil { // handle most common case first
			return nil
		}

		if n > 0 {
			data = data[n:]
			left -= n
		}

		if err != nil {
			if !IsTemporary(err) {
				return err
			} else {
				runtime.Gosched()
			}
		}
	}
	return nil
}

type flushable interface {
	Flush() error
}

func TryFlush(conn net.Conn) error {
	if f, ok := conn.(flushable); ok {
		for {
			err := f.Flush()
			if err == nil || !IsTemporary(err) {
				return err
			} else {
				runtime.Gosched()
			}
		}
	} else {
		return nil
	}
}
