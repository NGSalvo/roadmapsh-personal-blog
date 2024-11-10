package services

import (
	"io"
	"os"
)

type FileReader struct{}

func NewFileReader() FileReader {
	return FileReader{}
}

func (fr FileReader) Read(slug string) (string, error) {
	f, err := os.Open(slug + ".md")
	if err != nil {
		return "", nil
	}

	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}
