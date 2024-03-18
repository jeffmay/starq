package cmd_test

import (
	"io"
	"starq/cmd"
	"starq/sample"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fastjson"
)

type RunCommandOpt func(*cobra.Command)

func WithArgs(args ...string) RunCommandOpt {
	return func(cmd *cobra.Command) {
		cmd.SetArgs(args)
	}
}

func WithStdinFromReader(r io.Reader) RunCommandOpt {
	return func(cmd *cobra.Command) {
		cmd.SetIn(r)
	}
}

func WithStdinFromString(in string) RunCommandOpt {
	return WithStdinFromReader(strings.NewReader(in))
}

func RunCommand(cmd *cobra.Command, opts ...RunCommandOpt) (string, string, error) {
	stdoutSpy := new(strings.Builder)
	stderrSpy := new(strings.Builder)
	for _, opt := range opts {
		opt(cmd)
	}
	cmd.SetOut(stdoutSpy)
	cmd.SetErr(stderrSpy)
	err := cmd.Execute()
	return stdoutSpy.String(), stderrSpy.String(), err
}

func normalize(path sample.RelPath) string {
	return sample.MustNormalize("..", path)
}

func TestRootCmdShowsUsageOnEmptyArgs(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	cmd.InitRootCmd(rootCmd)
	stdout, stderr, err := RunCommand(rootCmd, WithArgs())
	require.NoError(t, err)
	require.Empty(t, stderr)
	require.Contains(t, stdout, rootCmd.Long)
	require.Contains(t, stdout, rootCmd.UsageString())
}

func TestRootCmdPetstoreReadOnlyConfig(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	cmd.InitRootCmd(rootCmd)
	stdout, stderr, err := RunCommand(rootCmd, WithArgs(normalize(sample.PETSTORE_TO_READONLY_STDOUT_PATH)))
	require.NoError(t, err)
	require.Empty(t, stderr)
	stdoutJSON, err := fastjson.Parse(stdout)
	require.NoError(t, err)
	title := stdoutJSON.Get("info", "title")
	require.NotNil(t, title)
	titleBytes, err := title.StringBytes()
	require.NoError(t, err)
	require.Equal(t, "Swagger Petstore", string(titleBytes))
}
