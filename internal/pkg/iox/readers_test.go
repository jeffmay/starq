package iox_test

import (
	"io"
	"strings"
	"testing"

	"github.com/jeffmay/starq/internal/pkg/iox"

	"github.com/stretchr/testify/require"
)

func TestProxyReadCloser(t *testing.T) {
	expected := "hello"
	reader := strings.NewReader(expected)
	var closed bool
	pw := iox.ProxyReadCloser(reader, func() error {
		closed = true
		return nil
	})
	bytes, err := io.ReadAll(pw)
	require.NoError(t, err)
	require.Equal(t, expected, string(bytes))
	require.False(t, closed)
	require.NoError(t, pw.Close())
	require.True(t, closed)
}
