package iox

import (
	"errors"
	"io"
)

type multiWriteCloser struct {
	multiWriter io.Writer
	writers     []io.WriteCloser
}

func (mwc *multiWriteCloser) Write(p []byte) (n int, err error) {
	return mwc.multiWriter.Write(p)
}

func (mwc *multiWriteCloser) Close() error {
	var allErrors []error
	for _, writer := range mwc.writers {
		if err := writer.Close(); err != nil {
			allErrors = append(allErrors, err)
		}
	}
	return errors.Join(allErrors...)
}

func MultiWriteCloser(writers ...io.WriteCloser) io.WriteCloser {
	onlyWriters := make([]io.Writer, len(writers))
	for i, w := range writers {
		onlyWriters[i] = w
	}
	multiWriter := io.MultiWriter(onlyWriters...)
	return &multiWriteCloser{
		multiWriter,
		writers,
	}
}

type nopWriteCloser struct {
	io.Writer
}

func (nwc *nopWriteCloser) Close() error {
	return nil
}

func NopWriterCloser(writer io.Writer) io.WriteCloser {
	if w, ok := writer.(io.WriteCloser); ok {
		return w
	}
	return &nopWriteCloser{
		Writer: writer,
	}
}

type proxyWriteCloser struct {
	io.Writer
	closeFn func() error
}

var _ io.WriteCloser = new(proxyWriteCloser)

func ProxyWriteCloser(writer io.Writer, closeFn func() error) io.WriteCloser {
	return &proxyWriteCloser{
		Writer:  writer,
		closeFn: closeFn,
	}
}

func (prc *proxyWriteCloser) Close() error {
	return prc.closeFn()
}
