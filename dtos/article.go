package dtos

import (
	"time"

	"github.com/a-h/templ"
)

type Article struct {
	Title   string          `toml:"title"`
	Content templ.Component `toml:"content"`
	Author  Author          `toml:"author"`
	Slug    string          `toml:"slug"`
	Date    time.Time       `toml:"date"`
}
