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
	if config == nil || config.Input == nil || len(config.Input.File) == 0 {
		return io.NopCloser(stdin), nil
	}
	if defnPath, ok := config.Source.AsFileName(); ok {
		defnDir := filepath.Dir(defnPath)
		absPath, err := filepath.Abs(filepath.Join(defnDir, config.Input.File))
		if err != nil {
			return nil, fmt.Errorf("invalid path (%s) from config at .input.file", config.Input.File)
		}
		file, err := os.Open(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open path (%s) from config at .input.file: %w", absPath, err)
		}
		return file, nil
	}
	return nil, nil
}

func (config *TransformerConfig) OutputWriter(stdout io.Writer) (io.WriteCloser, error) {
	if config == nil || config.Output == nil || len(config.Output.File) == 0 {
		return iox.NopWriterCloser(stdout), nil
	}
	if defnPath, ok := config.Source.AsFileName(); ok {
		defnDir := filepath.Dir(defnPath)
		outPath, err := filepath.Abs(filepath.Join(defnDir, config.Output.File))
		if err != nil {
			return nil, fmt.Errorf("invalid path (%s) from config at .output.file", config.Output.File)
		}
		// create the directory if it doesn't exist
		outDir := filepath.Dir(outPath)
		err = os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed to create directory (%s) from config at .output.file: %w", outDir, err)
		}
		// open the file and overwrite its contents
		outFile, err := os.Create(outPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file (%s) from config at .output.file: %w", outPath, err)
		}
		return outFile, nil
	}
	return nil, nil
}

func (config *TransformerConfig) ErrorsWriter(stderr io.Writer) (io.WriteCloser, error) {
	return iox.NopWriterCloser(stderr), nil
}

func (config *TransformerConfig) Close() error {
	return nil
}
