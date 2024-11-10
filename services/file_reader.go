package services

import (
	"io"
	"os"
	"strings"
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

func (fr FileReader) GetFileNames(path string) ([]string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, file := range dir {
		fileName := strings.Split(file.Name(), ".")[0]
		files = append(files, fileName)
	}

	return files, nil
}
