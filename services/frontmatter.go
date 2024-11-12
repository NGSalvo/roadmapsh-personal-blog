package services

import (
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
)

func Parse[T any](element string) (*T, error) {
	var data T

	_, err := frontmatter.Parse(strings.NewReader(element), &data)

	if err != nil {
		return nil, fmt.Errorf("error parsing frontmatter: %w", err)
	}

	return &data, nil
}
