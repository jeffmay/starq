package jsonx_test

import (
	"encoding/json"
	"starq/internal/jsonx"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONWriterOutputJSON(t *testing.T) {
	writer := jsonx.NewJSONWriter()
	_, err := writer.Write([]byte(`{"hello": "world"}`))
	require.NoError(t, err)
	outJSON, err := writer.OutputJSON()
	require.NoError(t, err)
	require.NotNil(t, outJSON)
	helloBytes, err := outJSON.Get("hello").StringBytes()
	require.NoError(t, err)
	require.Equal(t, "world", string(helloBytes))
}

func TestJSONWriterOutputJSONInvalid(t *testing.T) {
	writer := jsonx.NewJSONWriter()
	_, err := writer.Write([]byte(`{"hello": invalid}`))
	require.NoError(t, err)
	empty, err := writer.OutputJSON()
	require.Error(t, err)
	require.Empty(t, empty)
}

func TestJSONWriterMustOutputJSON(t *testing.T) {
	writer := jsonx.NewJSONWriter()
	_, err := writer.Write([]byte(`{"hello": "world"}`))
	require.NoError(t, err)
	outJSON := writer.MustOutputJSON()
	require.NotNil(t, outJSON)
	helloBytes, err := outJSON.Get("hello").StringBytes()
	require.NoError(t, err)
	require.Equal(t, "world", string(helloBytes))
}

func TestJSONWriterMustOutputJSONInvalid(t *testing.T) {
	writer := jsonx.NewJSONWriter()
	_, err := writer.Write([]byte(`{"hello": invalid`))
	require.NoError(t, err)
	require.Panics(t, func() {
		writer.MustOutputJSON()
	})
}

func TestJSONWriterOutputJSONPretty(t *testing.T) {
	writer := jsonx.NewJSONWriter()
	sampleJSONBytes := []byte(`{"hello": "world"}`)
	var sampleJSON interface{}
	err := json.Unmarshal(sampleJSONBytes, &sampleJSON)
	require.NoError(t, err)
	expectedBytes, err := json.MarshalIndent(sampleJSON, "", "  ")
	require.NoError(t, err)
	_, err = writer.Write(sampleJSONBytes)
	require.NoError(t, err)
	prettyBytes, err := writer.OutputJSONPretty()
	require.NoError(t, err)
	require.Equal(t, string(expectedBytes), string(prettyBytes))
}

func TestJSONWriterOutputJSONPrettyUpdates(t *testing.T) {
	writer := jsonx.NewJSONWriter()
	invalidJSONStr := `{"hello": `
	_, err := writer.Write([]byte(invalidJSONStr))
	require.NoError(t, err)
	empty, err := writer.OutputJSONPretty()
	require.Error(t, err)
	require.Empty(t, empty)
	// finish writing the JSON
	remainingJSONStr := `"world"}`
	_, err = writer.Write([]byte(remainingJSONStr))
	require.NoError(t, err)
	updatedPrettyJSON, err := writer.OutputJSONPretty()
	// we should not receive the same error as before
	require.NoError(t, err)
	expectedJSONStr := invalidJSONStr + remainingJSONStr
	var expectedJSON interface{}
	err = json.Unmarshal([]byte(expectedJSONStr), &expectedJSON)
	require.NoError(t, err)
	expectedPrettyJSONBytes, err := json.MarshalIndent(expectedJSON, "", "  ")
	require.NoError(t, err)
	require.Equal(t, string(expectedPrettyJSONBytes), updatedPrettyJSON)
}

func TestJSONWriterMustOutputJSONPretty(t *testing.T) {
	writer := jsonx.NewJSONWriter()
	sampleJSONBytes := []byte(`{"hello": "world"}`)
	var sampleJSON interface{}
	err := json.Unmarshal(sampleJSONBytes, &sampleJSON)
	require.NoError(t, err)
	expectedBytes, err := json.MarshalIndent(sampleJSON, "", "  ")
	require.NoError(t, err)
	_, err = writer.Write(sampleJSONBytes)
	require.NoError(t, err)
	prettyBytes := writer.MustOutputJSONPretty()
	require.Equal(t, string(expectedBytes), string(prettyBytes))
}

func TestJSONWriterMustOutputJSONPrettyInvalid(t *testing.T) {
	writer := jsonx.NewJSONWriter()
	invalidJSONBytes := []byte(`{"hello": invalid`)
	_, err := writer.Write(invalidJSONBytes)
	require.NoError(t, err)
	require.Panics(t, func() {
		writer.MustOutputJSONPretty()
	})
}

func TestNewJSONWriterSpy(t *testing.T) {
	sampleJSONStr := `{"hello": "world"}`
	writer := new(strings.Builder)
	jsonWriter := jsonx.NewJSONWriterSpy(writer)
	_, err := jsonWriter.Write([]byte(sampleJSONStr))
	require.NoError(t, err)
	require.Equal(t, sampleJSONStr, writer.String())
	require.Equal(t, sampleJSONStr, jsonWriter.Output())
}

func TestNewJSONWriterSpyNil(t *testing.T) {
	writerSpy := jsonx.NewJSONWriterSpy(nil)
	require.Nil(t, writerSpy)
}
