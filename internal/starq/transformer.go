package starq

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"starq/internal/iox"
)

// TODO: How should I handle closing streams here?
type Transformer interface {
	io.Closer
	Config() TransformerConfig
	InputReader(stdin io.Reader) (io.ReadCloser, error)
	OutputWriter(stdout io.Writer) (io.WriteCloser, error)
	ErrorsWriter(stderr io.Writer) (io.WriteCloser, error)
}

var _ Transformer = new(TransformerConfig)

func (config TransformerConfig) Config() TransformerConfig {
	return config
}

func (config *TransformerConfig) InputReader(stdin io.Reader) (io.ReadCloser, error) {
	if config == nil || config.Input == nil || config.Input.File == nil {
		return io.NopCloser(stdin), nil
	}
	if defnPath, ok := config.Source.AsFileName(); ok {
		defnDir := filepath.Dir(defnPath)
		absPath, err := filepath.Abs(filepath.Join(defnDir, *config.Input.File))
		if err != nil {
			return nil, fmt.Errorf("config .input.file not found '%s'", *config.Input.File)
		}
		file, err := os.Open(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config input.file at '%s': %w", absPath, err)
		}
		return file, nil
	}
	return nil, nil
}

func (config *TransformerConfig) OutputWriter(stdout io.Writer) (io.WriteCloser, error) {
	return iox.NopWriterCloser(stdout), nil
}

func (config *TransformerConfig) ErrorsWriter(stderr io.Writer) (io.WriteCloser, error) {
	return iox.NopWriterCloser(stderr), nil
}

func (config *TransformerConfig) Close() error {
	return nil
}
