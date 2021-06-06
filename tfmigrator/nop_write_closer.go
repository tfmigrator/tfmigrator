package tfmigrator

import "io"

type nopWriteCloser struct{}

func (wc *nopWriteCloser) Close() error {
	return nil
}

func (wc *nopWriteCloser) Write(p []byte) (n int, err error) {
	return io.Discard.Write(p) //nolint:wrapcheck
}
