package executor

import (
	"errors"
	"reflect"
	"time"
)

var (
	// ErrNotChanType means object is not channel, but used as channel
	ErrNotChanType = errors.New("not chan type")
	// ErrChanClosed means the channel waited from is closed
	ErrChanClosed = errors.New("chan closed")
	// ErrTimeout means waiting is timeout
	ErrTimeout = errors.New("wait chan timeout")
	// ErrOutOfRange should never be used
	ErrOutOfRange = errors.New("index out of range") // should never happen
)

// ReadChanWithTimeout read something from channel with timeout
func ReadChanWithTimeout(ch interface{}, timeout time.Duration) (interface{}, error) {
	v := reflect.ValueOf(ch)
	if v.Kind() != reflect.Chan {
		return nil, ErrNotChanType
	}

	t := time.NewTimer(timeout)
	defer t.Stop()

	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(t.C),
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: v,
		},
	}

	chosen, recv, ok := reflect.Select(cases)
	if !ok {
		return nil, ErrChanClosed
	}

	switch chosen {
	case 0:
		return nil, ErrTimeout
	case 1:
		return recv.Interface(), nil
	default: // should never happen
		return nil, ErrOutOfRange
	}
}

// WriteChanWithTimeout write something to channel with timeout
func WriteChanWithTimeout(ch interface{}, data interface{}, timeout time.Duration) error {
	v := reflect.ValueOf(ch)
	if v.Kind() != reflect.Chan {
		return ErrNotChanType
	}

	t := time.NewTimer(timeout)
	defer t.Stop()

	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(t.C),
		},
		{
			Dir:  reflect.SelectSend,
			Chan: v,
			Send: reflect.ValueOf(data),
		},
	}

	chosen, _, _ := reflect.Select(cases)

	switch chosen {
	case 0:
		return ErrTimeout
	case 1:
		return nil
	default: // should never happen
		return ErrOutOfRange
	}
}

// ClearChannel clear all content in channel
func ClearChannel(ch interface{}) error {
	v := reflect.ValueOf(ch)
	if v.Kind() != reflect.Chan {
		return ErrNotChanType
	}

	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: v,
		},
		{
			Dir: reflect.SelectDefault,
		},
	}

	for {
		chosen, _, ok := reflect.Select(cases)

		switch chosen {
		case 0:
			if !ok {
				return ErrChanClosed
			}
			continue
		case 1:
			return nil
		}
	}
}
