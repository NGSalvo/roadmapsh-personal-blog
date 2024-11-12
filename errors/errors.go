package errors

import (
	"errors"
)

var (
	ErrorOpeningFile      = errors.New("error opening file")
	ErrorReadingFile      = errors.New("error reading file")
	ErrorReadingDirectory = errors.New("error reading directory")
	ErrorArticleNotFound  = errors.New("article not found")
)

type ApplicationError struct {
	Message string `json:"message"`
}

func (e ApplicationError) Error() string {
	return e.Message
}
