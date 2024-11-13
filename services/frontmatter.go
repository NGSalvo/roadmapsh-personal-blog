package services

import (
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
)

type Frontmatter[T any] struct {
	Frontmatter T
	RemaingData []byte
}

func Parse[T any](element string) (*Frontmatter[T], error) {
	var data Frontmatter[T]

	remainingData, err := frontmatter.Parse(strings.NewReader(element), &data.Frontmatter)

	if err != nil {
		return nil, fmt.Errorf("error parsing frontmatter: %w", err)
	}

	data.RemaingData = remainingData

	return &data, nil
}
