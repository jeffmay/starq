package iox

import "io"

type proxyReadCloser struct {
	io.Reader
	closeFn func() error
}

var _ io.ReadCloser = new(proxyReadCloser)

func ProxyReadCloser(reader io.Reader, closeFn func() error) io.ReadCloser {
	if r, ok := reader.(io.ReadCloser); ok {
		return r
	}
	return &proxyReadCloser{
		Reader:  reader,
		closeFn: closeFn,
	}
}

func (prc *proxyReadCloser) Close() error {
	return prc.closeFn()
}
