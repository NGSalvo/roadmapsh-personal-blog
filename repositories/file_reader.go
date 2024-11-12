package repositories

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ngsalvo/roadmapsh-personal-blog/errors"
)

type FileReader struct{}

func NewFileReader() FileReader {
	return FileReader{}
}

func (fr FileReader) Read(slug string) (string, error) {
	f, err := os.Open(slug + ".md")
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrorOpeningFile, err.Error())
	}

	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrorReadingFile, err.Error())
	}

	if len(bytes) == 0 {
		return "", fmt.Errorf("%w: %s", errors.ErrorArticleNotFound, slug)
	}

	return string(bytes), nil
}

func (fr FileReader) GetFileNames(path string) ([]string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrorReadingDirectory, err.Error())
	}

	var files []string
	for _, file := range dir {
		fileName := strings.Split(file.Name(), ".")[0]
		files = append(files, fileName)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("%w: %s", errors.ErrorArticleNotFound, path)
	}

	return files, nil
}
