package handlers

import (
	"net/http"

	"github.com/ngsalvo/roadmapsh-personal-blog/services"
)

type (
	GetAdmin interface {
		Handle(w http.ResponseWriter, r *http.Request)
	}

	getAdmin struct {
		fileReader services.FileReader
	}
)

func NewGetAdmin(fileReader services.FileReader) GetAdmin {
	return &getAdmin{
		fileReader: fileReader,
	}
}

func (h *getAdmin) Handle(w http.ResponseWriter, r *http.Request) {
	return
}
