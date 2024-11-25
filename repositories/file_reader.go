package repositories

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/delaneyj/toolbelt"

	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
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

func (fr FileReader) Update(slug string, store *dtos.ArticleStore) error {
	data, err := fr.Read("static/blog/" + slug)

	if err != nil {
		return err
	}

	currentContent := strings.Split(data, "---")

	if len(currentContent) < 1 {
		return fmt.Errorf("%w: %s", errors.ErrorArticleNotFound, slug)
	}

	frontmatter := strings.Replace(currentContent[1], "toml", "", 1)
	md := strings.Join(currentContent[2:], "---")

	log.Printf("frontmatter: %s\nmd: %s", frontmatter, md)

	// 1. update frontmatter
	// 1.1 update title
	// 2. update article
	// 2.1 update content
	// 3 merge frontmatter and article
	// 4. write to file static/blog/{slug}.md and rename to title by kebab casing it

	// 1.1
	frontmatterProperties := strings.Split(frontmatter, "\n")

	for i, property := range frontmatterProperties {
		if strings.Contains(property, "title") {
			frontmatterProperties[i] = fmt.Sprintf("title = \"%s\"", store.Title)
		}
	}

	// 2.1
	md = store.Content

	// 3
	frontmatter = strings.Join(frontmatterProperties, "\n")
	updatedArticle := fmt.Sprintf("---toml\n%s\n---\n\n%s", frontmatter, md)

	//4
	err = os.WriteFile("static/blog/"+slug+".md", []byte(updatedArticle), 0644)
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrorWritingFile, err.Error())
	}

	os.Rename("static/blog/"+slug+".md", "static/blog/"+titleSlug(store.Title)+".md")

	return nil
}

func titleSlug(title string) string {
	illegalCharacters := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	sanitazed := illegalCharacters.ReplaceAllString(title, "")

	return toolbelt.Kebab(sanitazed)
}
