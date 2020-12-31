package api

import (
	"net"
	"net/url"
	"syscall"
)

type errRetryFilter func(error) bool

func isErrorNetRetryable(err error) bool {
	nerr, ok := err.(net.Error)
	if ok {
		if nerr.Timeout() || nerr.Temporary() {
			return true
		}

		switch t := nerr.(type) {
		case *net.AddrError, net.InvalidAddrError, *net.OpError, net.UnknownNetworkError:
			return t.Temporary() || t.Timeout()
		}
	}
	return false
}

func IsErrorRetryable(err error) bool {
	filters := []errRetryFilter{isErrorUrlRetryable, isErrorNetRetryable, isErrorSysRetryable}
	for _, filter := range filters {
		if filter(err) {
			return true
		}
	}
	return false
}

func isErrorSysRetryable(err error) bool {
	if serr, ok := err.(syscall.Errno); ok {
		switch serr {
		case syscall.ECONNABORTED, syscall.ECONNREFUSED, syscall.ECONNRESET, syscall.ENETDOWN, syscall.ETIMEDOUT:
			return true
		}
	}
	return false
}

func isErrorUrlRetryable(err error) bool {
	if uerr, ok := err.(*url.Error); ok {
		return uerr.Timeout() || uerr.Temporary()
	}
	return false
}
