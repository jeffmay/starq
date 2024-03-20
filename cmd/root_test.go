package cmd_test

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jeffmay/starq/cmd"
	"github.com/jeffmay/starq/sample"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
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

// The following tests only check if the command runs without errors and outputs the correct format.
// If the CLI tool outputs to the wrong stream, it will generate invalid JSON or YAML, which these tests will catch.
// The actual output of these files are tested in file://./../internal/starq/runner_test.go

func TestRootCmdPetstoreJSONtoReadOnly(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	cmd.InitRootCmd(rootCmd)
	stdout, stderr, err := RunCommand(rootCmd, WithArgs(normalize(sample.CONFIG_PETSTORE_JSON_TO_READONLY_PATH)))
	require.NoError(t, err)
	require.Empty(t, stderr)
	require.Empty(t, stdout)
	out, err := os.Open(normalize(sample.OUTPUT_PETSTORE_OPENAPI_READONLY_JSON_PATH))
	require.NoError(t, err)
	defer out.Close()
	outBytes, err := io.ReadAll(out)
	require.NoError(t, err)
	var outJSON map[string]any
	err = json.Unmarshal(outBytes, &outJSON)
	require.NoError(t, err)
}

func TestRootCmdPetstoreJSONtoYAML(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	cmd.InitRootCmd(rootCmd)
	stdout, stderr, err := RunCommand(rootCmd, WithArgs(normalize(sample.CONFIG_PETSTORE_JSON_TO_YAML_PATH)))
	require.NoError(t, err)
	require.Empty(t, stderr)
	require.Empty(t, stdout)
	outFile, err := os.Open(normalize(sample.OUTPUT_PETSTORE_OPENAPI_READONLY_JSON_PATH))
	require.NoError(t, err)
	defer outFile.Close()
	outBytes, err := io.ReadAll(outFile)
	require.NoError(t, err)
	var outYAML map[string]any
	err = yaml.Unmarshal(outBytes, &outYAML)
	require.NoError(t, err)
}
