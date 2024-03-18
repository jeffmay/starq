package starq_test

import (
	"starq/internal/jsonx"
	"starq/internal/starq"
	"starq/sample"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fastjson"
)

// Tests the [starq.runner] functions using stubbed input and output.

func normalize(path sample.RelPath) string {
	return sample.MustNormalize("../..", path)
}

func TestNilStreamsNotAllowed(t *testing.T) {
	require.Panics(t, func() {
		starq.NewRunner(nil, new(strings.Builder), new(strings.Builder))
	})
	require.Panics(t, func() {
		starq.NewRunner(strings.NewReader(""), nil, new(strings.Builder))
	})
	require.Panics(t, func() {
		starq.NewRunner(strings.NewReader(""), new(strings.Builder), nil)
	})
}

func TestLoadSimpleConfig(t *testing.T) {
	opts := MakeTestOpts().WithTransformers(MakeTestTransformer().FromConfigFile(normalize(sample.CONFIG_PETSTORE_JSON_TO_READONLY_PATH)))
	transformers, err := opts.LoadTransformers()
	require.NoError(t, err)
	require.Len(t, transformers[0].Config().Rules, 1)
}

func TestInvalidJqRule(t *testing.T) {
	opts := MakeTestOpts().WithRulesAppended("invalid")
	stdout := jsonx.NewJSONWriter()
	stderr := new(strings.Builder)
	runner := starq.NewRunner(strings.NewReader(""), stdout, stderr)
	err := runner.RunAllTransformers(opts)
	require.Error(t, err)
	require.Equal(t, err.Error(), "failed to run jq: exit status 3")
	require.Contains(t, stderr.String(), "jq: 1 compile error")
}

// TODO: This is outputing YAML as expected (since no file name is provided), but we should make validating YAML a little easier.
func TestPetstoreTitle(t *testing.T) {
	opts := MakeTestOpts().SetPrependRules(".info.title")
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	runner := starq.NewRunner(strings.NewReader(sample.PETSTORE_OPENAPI_JSON), stdout, stderr)
	err := runner.RunAllTransformers(opts)
	require.Empty(t, stderr.String())
	require.NoError(t, err, stderr.String())
	require.Equal(t, "Swagger Petstore\n", stdout.String())
}

func TestPetstoreJSONtoReadOnly(t *testing.T) {
	transformer := MakeTestTransformer().FromConfigFile(normalize(sample.CONFIG_PETSTORE_JSON_TO_READONLY_PATH))
	opts := MakeTestOpts().WithTransformers(transformer)
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	runner := starq.NewRunner(strings.NewReader(""), stdout, stderr)
	err := runner.RunAllTransformers(opts)
	require.NoError(t, err, stderr.String())
	require.Empty(t, stderr.String())
	require.Empty(t, stdout.String())
	out := transformer.Output()
	debugPretty, err := jsonx.Pretty(out)
	require.NoError(t, err)
	outJSON, err := fastjson.Parse(out)
	require.NoError(t, err)
	paths := outJSON.GetObject("paths")
	require.NotNilf(t, paths, ".paths should be defined in:\n%s", debugPretty)
	require.Equalf(t, paths.Len(), 2, ".paths should have length 2 in:\n%s", debugPretty)
	paths.Visit(func(key []byte, path *fastjson.Value) {
		getRoute := path.GetObject("get")
		require.NotNilf(t, getRoute, ".paths[\"%s\"].get should be defined in:\n%s", string(key), debugPretty)
		postRoute := path.GetObject("post")
		require.Nilf(t, postRoute, ".paths[\"%s\"].post should NOT be defined in:\n%s", string(key), debugPretty)
	})
	components := outJSON.GetObject("components")
	require.NotNilf(t, components, ".components should be defined in:\n%s", debugPretty)
}

func TestPetstoreJSONtoYAML(t *testing.T) {
	transformer := MakeTestTransformer().FromConfigFile(normalize(sample.CONFIG_PETSTORE_JSON_TO_YAML_PATH))
	opts := MakeTestOpts().WithTransformers(transformer)
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	runner := starq.NewRunner(strings.NewReader(""), stdout, stderr)
	err := runner.RunAllTransformers(opts)
	require.NoError(t, err, stderr.String())
	require.Empty(t, stderr.String())
	require.Empty(t, stdout.String())
	out := transformer.Output()
	require.NotEmpty(t, out)
	path, err := yaml.PathString("$.info.title")
	require.NoError(t, err)
	titleNode, err := path.ReadNode(strings.NewReader(out))
	require.NoError(t, err)
	require.Equal(t, "Swagger Petstore", titleNode.String())
}
