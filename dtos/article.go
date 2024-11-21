package dtos

import (
	"time"

	"github.com/a-h/templ"
)

type Article struct {
	Title         string          `toml:"title"`
	Content       templ.Component `toml:"content"`
	Author        Author          `toml:"author"`
	Slug          string          `toml:"slug"`
	Date          time.Time       `toml:"date"`
	ContentString string
}

type ArticleStore struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
