package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

/*
* ScraperService handles web scraping operations
 */
type ScraperService struct {
	collector *colly.Collector
}

/*
* NewScraperService creates a new scraper service instance
 */
func NewScraperService() *ScraperService {
	c := colly.NewCollector(
		colly.UserAgent("Texly.Chat Bot/1.0 (+https://texly.chat)"),
		colly.AllowURLRevisit(),
	)

	// Set timeout
	c.SetRequestTimeout(30 * time.Second)

	// Set rate limiting (polite crawling)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})

	return &ScraperService{
		collector: c,
	}
}

/*
* FetchAndClean scrapes a URL and returns cleaned text content
 */
func (s *ScraperService) FetchAndClean(url string) (string, error) {
	var content strings.Builder
	var title string
	var scrapingError error

	// Create a new collector instance for this request
	c := s.collector.Clone()

	// Extract title
	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	// Extract main content
	// Priority: main, article, or body content
	c.OnHTML("main, article, [role=main]", func(e *colly.HTMLElement) {
		// Remove unwanted elements
		e.DOM.Find("script, style, nav, header, footer, aside, .navigation, .menu, .sidebar, .ad, .advertisement").Remove()
		
		text := strings.TrimSpace(e.Text)
		if text != "" {
			content.WriteString(text)
			content.WriteString("\n\n")
		}
	})

	// Fallback: if no main content found, use body
	c.OnHTML("body", func(e *colly.HTMLElement) {
		if content.Len() == 0 {
			// Remove unwanted elements
			e.DOM.Find("script, style, nav, header, footer, aside, .navigation, .menu, .sidebar, .ad, .advertisement").Remove()
			
			text := strings.TrimSpace(e.Text)
			if text != "" {
				content.WriteString(text)
			}
		}
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		scrapingError = fmt.Errorf("failed to fetch URL: %w", err)
	})

	// Visit the URL
	if err := c.Visit(url); err != nil {
		return "", fmt.Errorf("failed to visit URL: %w", err)
	}

	// Wait for all requests to finish
	c.Wait()

	// Check for errors
	if scrapingError != nil {
		return "", scrapingError
	}

	// Combine title and content
	var result strings.Builder
	if title != "" {
		result.WriteString("# ")
		result.WriteString(title)
		result.WriteString("\n\n")
	}
	result.WriteString(content.String())

	// Clean up whitespace
	cleaned := cleanWhitespace(result.String())

	if cleaned == "" {
		return "", fmt.Errorf("no content extracted from URL")
	}

	return cleaned, nil
}

/*
* cleanWhitespace removes excessive whitespace and normalizes text
 */
func cleanWhitespace(text string) string {
	// Replace multiple newlines with double newline
	lines := strings.Split(text, "\n")
	var cleaned []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	
	return strings.Join(cleaned, "\n")
}
