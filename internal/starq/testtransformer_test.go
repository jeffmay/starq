package starq_test

import (
	"fmt"
	"io"
	"starq/internal/iox"
	"starq/internal/starq"
	"strings"
)

type TestTransformer struct {
	config               starq.TransformerConfig
	stubInputReader      io.ReadCloser
	stubInputReaderError error
	inputSpy             *strings.Builder
	outputSpy            *strings.Builder
	errorsSpy            *strings.Builder
}

var _ starq.Transformer = new(TestTransformer)

func (c TestTransformer) Config() starq.TransformerConfig {
	return c.config
}

func (c TestTransformer) InputReader(stdin io.Reader) (io.ReadCloser, error) {
	if c.stubInputReaderError != nil {
		// short-circuit with the stubbed error
		return nil, c.stubInputReaderError
	}
	reader := c.stubInputReader
	var err error = nil
	if c.stubInputReader == nil {
		reader, err = c.config.InputReader(stdin)
	}
	if reader == nil {
		return nil, err
	}
	// wrap the underlying reader with the inputSpy
	return iox.ProxyReadCloser(io.TeeReader(reader, c.inputSpy), reader.Close), err
}

func (c TestTransformer) OutputWriter(stdout io.Writer) (io.WriteCloser, error) {
	writer, err := c.config.OutputWriter(stdout)
	if err != nil {
		return nil, err
	}
	if writer == nil {
		return nil, err
	}
	return iox.ProxyWriteCloser(io.MultiWriter(c.outputSpy, writer), writer.Close), nil
}

func (c TestTransformer) ErrorsWriter(stderr io.Writer) (io.WriteCloser, error) {
	writer, err := c.config.ErrorsWriter(stderr)
	if err != nil {
		return nil, err
	}
	if writer == nil {
		return nil, err
	}
	return iox.ProxyWriteCloser(io.MultiWriter(c.errorsSpy, writer), writer.Close), nil
}

func (c TestTransformer) Close() error {
	return nil
}

// DefaultTestTransformer returns a [TestTransformer] with the [starq.DefaultTransformerConfig].
func DefaultTestTransformer() *TestTransformer {
	t := MakeTestTransformer()
	t.config = starq.DefaultTransformerConfig().Config()
	return t
}

func MakeTestTransformer() *TestTransformer {
	return &TestTransformer{
		config:               *starq.DefaultTransformerConfig(),
		stubInputReader:      nil,
		stubInputReaderError: nil,
		inputSpy:             new(strings.Builder),
		outputSpy:            new(strings.Builder),
		errorsSpy:            new(strings.Builder),
	}
}

func (c *TestTransformer) FromConfigFile(filename string) *TestTransformer {
	config, err := starq.ReadTransformerConfigFile(filename)
	if err != nil {
		panic(fmt.Errorf("could not parse config file '%s': %w", filename, err))
	}
	c.config = config
	if c.config.Input != nil && len(c.config.Input.File) > 0 {
		c.WithInputFile(c.config.Input.File)
	}
	return c
}

func (b *TestTransformer) WithRules(rules ...starq.Rule) *TestTransformer {
	b.config.Rules = rules
	return b
}

func (c *TestTransformer) WithInputFile(filename string) *TestTransformer {
	if c.config.Input == nil {
		c.config.Input = &starq.Input{}
	}
	c.config.Input.File = filename
	return c
}

// WithStubInputReader uses the given input reader instead of needing to load from a file
func (c *TestTransformer) WithStubInputReader(reader io.ReadCloser) *TestTransformer {
	c.stubInputReader = reader
	return c
}

func (c *TestTransformer) WithStubInputString(input string) *TestTransformer {
	return c.WithStubInputReader(io.NopCloser(strings.NewReader(input)))
}

func (c *TestTransformer) WithStubInputError(err error) *TestTransformer {
	c.stubInputReaderError = err
	return c
}

func (c TestTransformer) Output() string {
	return c.outputSpy.String()
}
