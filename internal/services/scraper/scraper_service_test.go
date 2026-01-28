package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScraperService_FetchAndClean_Success(t *testing.T) {
	// Create a test server with sample HTML
	testHTML := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
			<script>alert('test');</script>
			<style>.test { color: red; }</style>
		</head>
		<body>
			<header>Header content</header>
			<nav>Navigation</nav>
			<main>
				<h1>Main Content</h1>
				<p>This is the main content of the page.</p>
				<p>Another paragraph with information.</p>
			</main>
			<aside>Sidebar content</aside>
			<footer>Footer</footer>
			<script>console.log('test');</script>
		</body>
		</html>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	scraper := NewScraperService()
	content, err := scraper.FetchAndClean(server.URL)

	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	// Should include title
	assert.Contains(t, content, "Test Page")

	// Should include main content
	assert.Contains(t, content, "Main Content")
	assert.Contains(t, content, "main content of the page")

	// Should NOT include script or style content
	assert.NotContains(t, content, "alert('test')")
	assert.NotContains(t, content, "color: red")
	assert.NotContains(t, content, "console.log")
}

func TestScraperService_FetchAndClean_WithArticle(t *testing.T) {
	testHTML := `
		<!DOCTYPE html>
		<html>
		<head><title>Article Page</title></head>
		<body>
			<nav class="menu">Menu items</nav>
			<article>
				<h1>Article Title</h1>
				<p>First paragraph of the article.</p>
				<p>Second paragraph with details.</p>
			</article>
			<div class="ad">Advertisement</div>
		</body>
		</html>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	scraper := NewScraperService()
	content, err := scraper.FetchAndClean(server.URL)

	assert.NoError(t, err)
	assert.Contains(t, content, "Article Title")
	assert.Contains(t, content, "First paragraph")
	assert.Contains(t, content, "Second paragraph")
}

func TestScraperService_FetchAndClean_404Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	scraper := NewScraperService()
	content, err := scraper.FetchAndClean(server.URL)

	assert.Error(t, err)
	assert.Empty(t, content)
}

func TestScraperService_FetchAndClean_EmptyContent(t *testing.T) {
	testHTML := `
		<!DOCTYPE html>
		<html>
		<head><title>Empty Page</title></head>
		<body>
			<script>window.location='/';</script>
		</body>
		</html>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	scraper := NewScraperService()
	content, err := scraper.FetchAndClean(server.URL)

	// Script content is removed, but the page has title  
	// so it won't be empty
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "Empty Page")
}

func TestScraperService_FetchAndClean_MultipleMainElements(t *testing.T) {
	testHTML := `
		<!DOCTYPE html>
		<html>
		<head><title>Multi Main</title></head>
		<body>
			<main>
				<p>First main section.</p>
			</main>
			<main>
				<p>Second main section.</p>
			</main>
		</body>
		</html>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	scraper := NewScraperService()
	content, err := scraper.FetchAndClean(server.URL)

	assert.NoError(t, err)
	// Should extract content from both main elements
	assert.Contains(t, content, "First main section")
	assert.Contains(t, content, "Second main section")
}

func TestScraperService_FetchAndClean_FallbackToBody(t *testing.T) {
	// HTML without main/article tags
	testHTML := `
		<!DOCTYPE html>
		<html>
		<head><title>Simple Page</title></head>
		<body>
			<h1>Page Title</h1>
			<p>This is regular body content.</p>
			<div class="ad">Ad content</div>
		</body>
		</html>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	scraper := NewScraperService()
	content, err := scraper.FetchAndClean(server.URL)

	assert.NoError(t, err)
	assert.Contains(t, content, "Page Title")
	assert.Contains(t, content, "regular body content")
}

func TestScraperService_FetchAndClean_InvalidURL(t *testing.T) {
	scraper := NewScraperService()
	content, err := scraper.FetchAndClean("not-a-valid-url")

	assert.Error(t, err)
	assert.Empty(t, content)
}

func TestCleanWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Multiple newlines",
			input:    "Line 1\n\n\n\nLine 2\n\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "Trailing whitespace",
			input:    "  Text with spaces  \n  More text  ",
			expected: "Text with spaces\nMore text",
		},
		{
			name:     "Empty lines",
			input:    "Text\n\n\n\nMore text",
			expected: "Text\nMore text",
		},
		{
			name:     "Already clean",
			input:    "Clean text\nAnother line",
			expected: "Clean text\nAnother line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanWhitespace(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestScraperService_UserAgent(t *testing.T) {
	var capturedUserAgent string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserAgent = r.UserAgent()
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><main>Test content</main></body></html>`))
	}))
	defer server.Close()

	scraper := NewScraperService()
	_, err := scraper.FetchAndClean(server.URL)

	assert.NoError(t, err)
	assert.Contains(t, capturedUserAgent, "Texly.Chat Bot")
}
