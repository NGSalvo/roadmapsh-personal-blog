package dtos

import "github.com/a-h/templ"

type Article struct {
	Title   string          `toml:"title"`
	Content templ.Component `toml:"content"`
	Author  Author          `toml:"author"`
	Slug    string          `toml:"slug"`
}
