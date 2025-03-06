package blog

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
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
	Image       string
}

func LoadPost(slug string) (*Post, error) {
	//Check if posts exists
	postPath := filepath.Join("posts", slug+".md")
	content, err := os.ReadFile(postPath)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	// Split the content into lines to extract title and modify content
	lines := strings.Split(string(content), "\n")
	var title, summary, image string
	if len(lines) > 0 && strings.HasPrefix(lines[0], "#") {
		title = strings.TrimSpace(strings.TrimPrefix(lines[0], "#"))
		// Remove the title line for HTML conversion
		content = []byte(strings.Join(lines[1:], "\n"))
	}

	// Extract the first image
	for _, line := range lines {
		if strings.HasPrefix(line, "![") && strings.Contains(line, "](") {
			start := strings.Index(line, "](") + 2
			end := strings.Index(line, ")")
			if start > 1 && end > start {
				image = line[start:end]
				break
			}
		}
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("rrt"),
			),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(), // Allow raw HTML
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}

	contentHTML := buf.String()

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
		Image:       image,
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
