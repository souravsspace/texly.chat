package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/*
 * URLSet represents the root element of a sitemap XML
 */
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

/*
 * URL represents a single URL entry in a sitemap
 */
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

/*
 * SitemapIndex represents a sitemap index XML (points to multiple sitemaps)
 */
type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

/*
 * Sitemap represents a single sitemap entry in a sitemap index
 */
type Sitemap struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

/*
 * SitemapParser handles sitemap parsing operations
 */
type SitemapParser struct {
	client  *http.Client
	maxURLs int // Maximum number of URLs to extract
}

/*
 * NewSitemapParser creates a new sitemap parser instance
 */
func NewSitemapParser(maxURLs int) *SitemapParser {
	if maxURLs <= 0 {
		maxURLs = 1000 // Default maximum
	}

	return &SitemapParser{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxURLs: maxURLs,
	}
}

/*
 * ParseSitemap fetches and parses a sitemap from a given URL or domain
 * It automatically discovers sitemap.xml and handles sitemap indexes
 */
func (p *SitemapParser) ParseSitemap(baseURL string) ([]string, error) {
	// Normalize URL
	sitemapURL, err := p.discoverSitemapURL(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover sitemap: %w", err)
	}

	fmt.Printf("Discovered sitemap URL: %s\n", sitemapURL)

	// Fetch and parse sitemap
	urls, err := p.fetchAndParse(sitemapURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sitemap: %w", err)
	}

	// Filter and limit URLs
	filtered := p.filterURLs(urls)
	if len(filtered) > p.maxURLs {
		filtered = filtered[:p.maxURLs]
	}

	return filtered, nil
}

/*
 * discoverSitemapURL tries to find the sitemap URL from various common locations
 */
func (p *SitemapParser) discoverSitemapURL(baseURL string) (string, error) {
	// Parse base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Ensure scheme is present
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	// If the URL already points to a sitemap file, use it directly
	if strings.HasSuffix(strings.ToLower(parsedURL.Path), ".xml") {
		return parsedURL.String(), nil
	}

	// Try common sitemap locations
	sitemapURLs := []string{
		fmt.Sprintf("%s://%s/sitemap.xml", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/sitemap_index.xml", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/sitemap-index.xml", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/sitemap/sitemap.xml", parsedURL.Scheme, parsedURL.Host),
	}

	// Try to fetch robots.txt first
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)
	if sitemapFromRobots, err := p.findSitemapInRobots(robotsURL); err == nil && sitemapFromRobots != "" {
		return sitemapFromRobots, nil
	}

	// Try each common sitemap location
	for _, sitemapURL := range sitemapURLs {
		if p.urlExists(sitemapURL) {
			return sitemapURL, nil
		}
	}

	return "", fmt.Errorf("no sitemap found for domain: %s", parsedURL.Host)
}

/*
 * findSitemapInRobots tries to find sitemap URL in robots.txt
 */
func (p *SitemapParser) findSitemapInRobots(robotsURL string) (string, error) {
	resp, err := p.client.Get(robotsURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("robots.txt not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse robots.txt for Sitemap directive
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "sitemap:") {
			sitemapURL := strings.TrimSpace(line[8:]) // Remove "Sitemap:" prefix
			return sitemapURL, nil
		}
	}

	return "", fmt.Errorf("no sitemap found in robots.txt")
}

/*
 * urlExists checks if a URL returns a 200 status code
 */
func (p *SitemapParser) urlExists(url string) bool {
	resp, err := p.client.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

/*
 * fetchAndParse fetches a sitemap URL and parses all URLs (handles sitemap indexes)
 */
func (p *SitemapParser) fetchAndParse(sitemapURL string) ([]string, error) {
	resp, err := p.client.Get(sitemapURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sitemap: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sitemap returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read sitemap body: %w", err)
	}

	// Try parsing as sitemap index first
	var sitemapIndex SitemapIndex
	if err := xml.Unmarshal(body, &sitemapIndex); err == nil && len(sitemapIndex.Sitemaps) > 0 {
		fmt.Printf("Found sitemap index with %d sitemaps\n", len(sitemapIndex.Sitemaps))
		return p.parseSitemapIndex(sitemapIndex)
	}

	// Parse as regular sitemap
	var urlSet URLSet
	if err := xml.Unmarshal(body, &urlSet); err != nil {
		return nil, fmt.Errorf("failed to parse sitemap XML: %w", err)
	}

	urls := make([]string, 0, len(urlSet.URLs))
	for _, u := range urlSet.URLs {
		if u.Loc != "" {
			urls = append(urls, u.Loc)
		}
	}

	return urls, nil
}

/*
 * parseSitemapIndex parses a sitemap index and fetches all child sitemaps
 */
func (p *SitemapParser) parseSitemapIndex(index SitemapIndex) ([]string, error) {
	var allURLs []string

	for _, sitemap := range index.Sitemaps {
		if len(allURLs) >= p.maxURLs {
			break
		}

		fmt.Printf("Fetching child sitemap: %s\n", sitemap.Loc)
		urls, err := p.fetchAndParse(sitemap.Loc)
		if err != nil {
			// Log error but continue with other sitemaps
			fmt.Printf("Warning: Failed to parse child sitemap %s: %v\n", sitemap.Loc, err)
			continue
		}

		allURLs = append(allURLs, urls...)
	}

	return allURLs, nil
}

/*
 * filterURLs filters out non-content URLs (images, PDFs, etc.)
 */
func (p *SitemapParser) filterURLs(urls []string) []string {
	// File extensions to exclude
	excludedExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".ico",
		".pdf", ".doc", ".docx", ".xls", ".xlsx",
		".zip", ".tar", ".gz", ".rar",
		".mp3", ".mp4", ".avi", ".mov", ".wmv",
		".css", ".js", ".json", ".xml",
	}

	filtered := make([]string, 0, len(urls))

	for _, u := range urls {
		// Parse URL
		parsedURL, err := url.Parse(u)
		if err != nil {
			continue
		}

		// Check if URL ends with excluded extension
		isExcluded := false
		lowerPath := strings.ToLower(parsedURL.Path)
		for _, ext := range excludedExtensions {
			if strings.HasSuffix(lowerPath, ext) {
				isExcluded = true
				break
			}
		}

		if !isExcluded {
			filtered = append(filtered, u)
		}
	}

	return filtered
}
