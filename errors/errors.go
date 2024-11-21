package errors

import (
	"errors"
)

var (
	ErrorOpeningFile      = errors.New("error opening file")
	ErrorReadingFile      = errors.New("error reading file")
	ErrorReadingDirectory = errors.New("error reading directory")
	ErrorWritingFile      = errors.New("error writing file")
	ErrorArticleNotFound  = errors.New("article not found")
)

type ApplicationError struct {
	Message string `json:"message"`
}

func (e ApplicationError) Error() string {
	return e.Message
}
