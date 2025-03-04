package blog

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Post struct {
	Title       string
	Slug        string
	Content     string
	ContentHTML string
	Date        time.Time
	Summary     string
}

func LoadPost(slug string) (*Post, error) {
	//Check if posts exists
	postPath := filepath.Join("posts", slug+".md")
	content, err := os.ReadFile(postPath)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	// Extract the first line as the title, removing the # if present
	lines := strings.Split(string(content), "\n")
	title := strings.TrimPrefix(lines[0], "# ")

	//convert markdown to HTML
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps(), html.WithXHTML(), html.WithUnsafe()),
	)
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}

	contentHTML := buf.String()

	// Creates a summary of the first paragraph of 100 chars
	//summary := contentHTML
	//Find the first meaningful paragraph for summary
	var summary string
	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "!") && !strings.HasPrefix(trimmed, "#") {
			summary = trimmed
			break
		}
	}

	return &Post{
		Title:       title,
		Slug:        slug,
		Content:     string(content),
		ContentHTML: contentHTML,
		Date:        time.Now(),
		Summary:     summary,
	}, nil
}

func ListPosts() ([]*Post, error) {
	var posts []*Post

	err := filepath.WalkDir("posts", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			// Extract the slug from the file name
			filename := filepath.Base(path)
			slug := strings.TrimSuffix(filename, ".md")

			post, err := LoadPost(slug)
			if err != nil {
				return err
			}
			posts = append(posts, post)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	return posts, nil
}
