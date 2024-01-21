package starq_test

import (
	"starq/internal/starq"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSample(t *testing.T) {
	out := new(strings.Builder)
	opts := MakeOpts(
		WithConfigFile("./fake/simple.yaml"),
		WithOutputWriter(out),
	)
	err := starq.Run(opts)
	require.NoError(t, err)
}
