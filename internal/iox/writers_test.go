package iox_test

import (
	"errors"
	"io"
	"starq/internal/iox"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testWriteCloser struct {
	writer       strings.Builder
	closed       bool
	errorOnClose error
}

var _ io.WriteCloser = new(testWriteCloser)

func (t *testWriteCloser) Write(p []byte) (n int, err error) {
	return t.writer.Write(p)
}

func (t *testWriteCloser) Close() error {
	if t.errorOnClose != nil {
		return t.errorOnClose
	}
	t.closed = true
	return nil
}

func TestMultiWriteCloser(t *testing.T) {
	one := new(testWriteCloser)
	two := new(testWriteCloser)
	mwc := iox.MultiWriteCloser(one, two)
	expected := "hello"
	i, err := io.WriteString(mwc, expected)
	require.NoError(t, err)
	require.Equal(t, i, len(expected))
	require.Equal(t, one.writer.String(), expected)
	require.Equal(t, two.writer.String(), expected)
	err = mwc.Close()
	require.NoError(t, err)
	require.True(t, one.closed)
	require.True(t, two.closed)
}

func TestMultiWriteCloserCapturesCloseErrors(t *testing.T) {
	one := new(testWriteCloser)
	one.errorOnClose = errors.New("one failed")
	two := new(testWriteCloser)
	two.errorOnClose = errors.New("two failed")
	three := new(testWriteCloser) // succeeds
	mwc := iox.MultiWriteCloser(one, two, three)
	expected := "hello"
	i, err := io.WriteString(mwc, expected)
	require.NoError(t, err)
	require.Equal(t, i, len(expected))
	err = mwc.Close()
	require.Error(t, err)
	closeErrors := err.(interface{ Unwrap() []error }).Unwrap()
	require.Len(t, closeErrors, 2)
	require.Equal(t, closeErrors[0].Error(), "one failed")
	require.Equal(t, closeErrors[1].Error(), "two failed")
	require.True(t, three.closed)
}

func TestNopWriteCloser(t *testing.T) {
	var writer strings.Builder
	nwc := iox.NopWriterCloser(&writer)
	expected := "hello"
	i, err := io.WriteString(nwc, expected)
	require.NoError(t, err)
	require.Equal(t, i, len(expected))
	require.Equal(t, writer.String(), expected)
	require.NoError(t, nwc.Close())
}

func TestProxyWriteCloser(t *testing.T) {
	var writer strings.Builder
	var closed bool
	pw := iox.ProxyWriteCloser(&writer, func() error {
		closed = true
		return nil
	})
	expected := "hello"
	i, err := io.WriteString(pw, expected)
	require.NoError(t, err)
	require.Equal(t, len(expected), i)
	require.Equal(t, expected, writer.String())
	require.False(t, closed)
	require.NoError(t, pw.Close())
	require.True(t, closed)
}
