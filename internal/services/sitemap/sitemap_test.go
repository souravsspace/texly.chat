package sitemap

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSitemapParser(t *testing.T) {
	parser := NewSitemapParser(100)
	if parser == nil {
		t.Fatal("Expected parser to be created")
	}
	if parser.maxURLs != 100 {
		t.Errorf("Expected maxURLs to be 100, got %d", parser.maxURLs)
	}

	// Test default maxURLs
	parser2 := NewSitemapParser(0)
	if parser2.maxURLs != 1000 {
		t.Errorf("Expected default maxURLs to be 1000, got %d", parser2.maxURLs)
	}
}

func TestParseSitemap_BasicSitemap(t *testing.T) {
	// Create test server with a simple sitemap
	sitemapXML := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url>
		<loc>https://example.com/page1</loc>
		<lastmod>2024-01-01</lastmod>
	</url>
	<url>
		<loc>https://example.com/page2</loc>
	</url>
	<url>
		<loc>https://example.com/page3</loc>
	</url>
</urlset>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sitemap.xml" {
			w.Header().Set("Content-Type", "application/xml")
			w.Write([]byte(sitemapXML))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	parser := NewSitemapParser(100)
	urls, err := parser.ParseSitemap(server.URL + "/sitemap.xml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(urls) != 3 {
		t.Errorf("Expected 3 URLs, got %d", len(urls))
	}

	expectedURLs := []string{
		"https://example.com/page1",
		"https://example.com/page2",
		"https://example.com/page3",
	}

	for i, url := range urls {
		if url != expectedURLs[i] {
			t.Errorf("Expected URL %s, got %s", expectedURLs[i], url)
		}
	}
}

func TestParseSitemap_SitemapIndex(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")

		if r.URL.Path == "/sitemap-index.xml" {
			// Sitemap index
			indexXML := `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<sitemap>
		<loc>` + server.URL + `/sitemap1.xml</loc>
	</sitemap>
	<sitemap>
		<loc>` + server.URL + `/sitemap2.xml</loc>
	</sitemap>
</sitemapindex>`
			w.Write([]byte(indexXML))
		} else if r.URL.Path == "/sitemap1.xml" {
			// First child sitemap
			sitemap1 := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url><loc>https://example.com/page1</loc></url>
	<url><loc>https://example.com/page2</loc></url>
</urlset>`
			w.Write([]byte(sitemap1))
		} else if r.URL.Path == "/sitemap2.xml" {
			// Second child sitemap
			sitemap2 := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url><loc>https://example.com/page3</loc></url>
	<url><loc>https://example.com/page4</loc></url>
</urlset>`
			w.Write([]byte(sitemap2))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	parser := NewSitemapParser(100)
	urls, err := parser.ParseSitemap(server.URL + "/sitemap-index.xml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(urls) != 4 {
		t.Errorf("Expected 4 URLs from index, got %d", len(urls))
	}
}

func TestFilterURLs(t *testing.T) {
	parser := NewSitemapParser(100)

	testURLs := []string{
		"https://example.com/page1",
		"https://example.com/page2",
		"https://example.com/image.jpg",
		"https://example.com/document.pdf",
		"https://example.com/page3",
		"https://example.com/style.css",
		"https://example.com/script.js",
		"https://example.com/page4.html",
	}

	filtered := parser.filterURLs(testURLs)

	expected := 3 // page1, page2, page3, page4.html
	if len(filtered) != 4 {
		t.Errorf("Expected %d filtered URLs, got %d", expected, len(filtered))
	}

	// Verify no image or PDF URLs
	for _, url := range filtered {
		if url == "https://example.com/image.jpg" || url == "https://example.com/document.pdf" {
			t.Errorf("URL should have been filtered: %s", url)
		}
	}
}

func TestDiscoverSitemapURL(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			robotsTxt := "User-agent: *\nDisallow: /admin\nSitemap: " + server.URL + "/custom-sitemap.xml"
			w.Write([]byte(robotsTxt))
		} else if r.URL.Path == "/custom-sitemap.xml" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/sitemap.xml" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	parser := NewSitemapParser(100)

	// Test robots.txt discovery
	discoveredURL, err := parser.discoverSitemapURL(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if discoveredURL != server.URL+"/custom-sitemap.xml" {
		t.Errorf("Expected custom sitemap from robots.txt, got %s", discoveredURL)
	}
}

func TestParseSitemap_MaxURLsLimit(t *testing.T) {
	// Create sitemap with many URLs
	var urlSet URLSet
	urlSet.XMLName = xml.Name{Local: "urlset"}
	for i := 0; i < 150; i++ {
		urlSet.URLs = append(urlSet.URLs, URL{
			Loc: "https://example.com/page" + string(rune(i)),
		})
	}

	sitemapXML, _ := xml.MarshalIndent(urlSet, "", "  ")
	xmlWithHeader := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(sitemapXML)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(xmlWithHeader))
	}))
	defer server.Close()

	// Set max URLs to 50
	parser := NewSitemapParser(50)
	urls, err := parser.ParseSitemap(server.URL + "/sitemap.xml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(urls) != 50 {
		t.Errorf("Expected 50 URLs (limited), got %d", len(urls))
	}
}

func TestParseSitemap_InvalidXML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is not XML"))
	}))
	defer server.Close()

	parser := NewSitemapParser(100)
	_, err := parser.ParseSitemap(server.URL + "/sitemap.xml")
	if err == nil {
		t.Error("Expected error for invalid XML")
	}
}

func TestParseSitemap_404Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	parser := NewSitemapParser(100)
	_, err := parser.ParseSitemap(server.URL)
	if err == nil {
		t.Error("Expected error when sitemap not found")
	}
}
