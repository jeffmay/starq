package jsonx

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type jsonWriter struct {
	writer     io.Writer
	writerSpy  *strings.Builder
	parsed     *fastjson.Value
	parsedHash hashCode
	pretty     string
	prettyHash hashCode
}

var _ io.Writer = new(jsonWriter)

func NewJSONWriter() *jsonWriter {
	writerSpy := new(strings.Builder)
	return &jsonWriter{
		writer:    writerSpy,
		writerSpy: writerSpy,
	}
}

func NewJSONWriterSpy(writer io.Writer) *jsonWriter {
	if writer == nil {
		return nil
	}
	writerSpy := new(strings.Builder)
	teeWriter := io.MultiWriter(writer, writerSpy)
	return &jsonWriter{
		writer:    teeWriter,
		writerSpy: writerSpy,
	}
}

func (j jsonWriter) Write(p []byte) (n int, err error) {
	return j.writer.Write(p)
}

func (j jsonWriter) Output() string {
	return j.writerSpy.String()
}

func (j jsonWriter) OutputJSON() (*fastjson.Value, error) {
	out := j.Output()
	if j.parsedHash != hashString(out) {
		v, err := fastjson.Parse(out)
		if err != nil {
			return nil, fmt.Errorf("could not parse JSON: %w", err)
		}
		j.parsed = v
		j.parsedHash = hashString(out)
	}
	return j.parsed, nil
}

func (j jsonWriter) MustOutputJSON() *fastjson.Value {
	v, err := j.OutputJSON()
	if err != nil {
		panic(err)
	}
	return v
}

func (o jsonWriter) OutputJSONPretty() (string, error) {
	out := o.Output()
	outHash := hashString(out)
	if o.prettyHash != outHash {
		var obj interface{}
		err := json.Unmarshal([]byte(out), &obj)
		if err != nil {
			return "", fmt.Errorf("could not unmarshal invalid JSON: %w", err)
		}
		indented, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return "", fmt.Errorf("could not remarshal JSON: %w", err)
		}
		o.pretty = string(indented)
		o.prettyHash = hashString(o.pretty)
	}
	return o.pretty, nil
}

func (j jsonWriter) MustOutputJSONPretty() string {
	pretty, err := j.OutputJSONPretty()
	if err != nil {
		panic(err)
	}
	return pretty
}
