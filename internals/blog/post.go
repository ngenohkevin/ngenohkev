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
	"regexp"
	"strings"
	"sync"
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

// Define a regular expression for parsing date from frontmatter
var dateRegex = regexp.MustCompile(`(?m)^date:\s*(.+)$`)

// Add a cache to store posts
var (
	postCache    = make(map[string]*Post)
	postCacheMux sync.RWMutex
	cacheEnabled = true
)

func parseDate(content string) (time.Time, string) {
	// Default date is now
	date := time.Now()

	// Look for date in frontmatter if there is any
	if match := dateRegex.FindStringSubmatch(content); len(match) > 1 {
		// Try to parse the date in various formats
		formats := []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			"January 2, 2006",
			"Jan 2, 2006",
			time.RFC3339,
		}

		for _, format := range formats {
			if parsedDate, err := time.Parse(format, strings.TrimSpace(match[1])); err == nil {
				date = parsedDate
				// Remove the date line from content
				content = dateRegex.ReplaceAllString(content, "")
				break
			}
		}
	}

	// Clean up any leftover whitespace from removed frontmatter
	content = strings.TrimSpace(content)

	return date, content
}

func LoadPost(slug string) (*Post, error) {
	// Check cache first if enabled
	if cacheEnabled {
		postCacheMux.RLock()
		if post, found := postCache[slug]; found {
			postCacheMux.RUnlock()
			return post, nil
		}
		postCacheMux.RUnlock()
	}

	// Check if post exists
	postPath := filepath.Join("posts", slug+".md")
	content, err := os.ReadFile(postPath)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	contentStr := string(content)

	// Parse date from content
	date, contentAfterDateParsing := parseDate(contentStr)

	// Split the content into lines to extract title and modify content
	lines := strings.Split(contentAfterDateParsing, "\n")
	var title, summary, image string
	if len(lines) > 0 && strings.HasPrefix(lines[0], "#") {
		title = strings.TrimSpace(strings.TrimPrefix(lines[0], "#"))
		// Remove the title line for HTML conversion
		contentAfterDateParsing = strings.Join(lines[1:], "\n")
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
	if err := md.Convert([]byte(contentAfterDateParsing), &buf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}

	contentHTML := buf.String()

	// Extract summary from content
	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "!") && !strings.HasPrefix(trimmed, "#") {
			//limit the number of words to 35
			words := strings.Fields(trimmed)
			if len(words) > 35 {
				summary = strings.Join(words[:35], " ") + "....."
			} else {
				summary = trimmed
			}
			break
		}
	}

	post := &Post{
		Title:       title,
		Slug:        slug,
		Content:     contentAfterDateParsing,
		ContentHTML: contentHTML,
		Date:        date,
		Summary:     summary,
		Image:       image,
	}

	// Store in cache if enabled
	if cacheEnabled {
		postCacheMux.Lock()
		postCache[slug] = post
		postCacheMux.Unlock()
	}

	return post, nil
}

func ListPosts() ([]*Post, error) {
	var posts []*Post
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var globalErr error

	// Process files in parallel
	err := filepath.WalkDir("posts", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			// Extract the slug from the file name
			filename := filepath.Base(path)
			slug := strings.TrimSuffix(filename, ".md")

			wg.Add(1)
			go func(slug string) {
				defer wg.Done()

				post, err := LoadPost(slug)
				if err != nil {
					mutex.Lock()
					globalErr = err
					mutex.Unlock()
					return
				}

				mutex.Lock()
				posts = append(posts, post)
				mutex.Unlock()
			}(slug)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	wg.Wait()

	if globalErr != nil {
		return nil, globalErr
	}

	// Sort posts by date (most recent first)
	sortPostsByDate(posts)

	return posts, nil
}

// Helper function to sort posts by date (newest first)
func sortPostsByDate(posts []*Post) {
	for i := 0; i < len(posts)-1; i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].Date.Before(posts[j].Date) {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
}

// Functions to manage caching

// EnableCache turns on post caching
func EnableCache() {
	cacheEnabled = true
}

// DisableCache turns off post caching
func DisableCache() {
	cacheEnabled = false
}

// ClearCache removes all cached posts
func ClearCache() {
	postCacheMux.Lock()
	postCache = make(map[string]*Post)
	postCacheMux.Unlock()
}
