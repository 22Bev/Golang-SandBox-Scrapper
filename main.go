package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Quote represents a single quote with its metadata
type Quote struct {
	Text   string   `json:"text"`
	Author string   `json:"author"`
	Tags   []string `json:"tags"`
}

// Scraper represents our web scraper
type Scraper struct {
	client  *http.Client
	baseURL string
}

// NewScraper creates a new scraper instance
func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "http://quotes.toscrape.com",
	}
}

// scrape a single page and return quotes
func (s *Scraper) scrapePage(pageURL string) ([]Quote, error) {
	// Add delay to be respectful to the server
	time.Sleep(1 * time.Second)

	res, err := s.client.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var quotes []Quote
	doc.Find(".quote").Each(func(i int, s *goquery.Selection) {
		quote := Quote{
			Text:   strings.TrimSpace(s.Find(".text").Text()),
			Author: strings.TrimSpace(s.Find(".author").Text()),
		}

		s.Find(".tags .tag").Each(func(i int, t *goquery.Selection) {
			quote.Tags = append(quote.Tags, strings.TrimSpace(t.Text()))
		})

		quotes = append(quotes, quote)
	})

	return quotes, nil
}

// ScrapeAllQuotes scrapes all pages and returns all quotes
func (s *Scraper) ScrapeAllQuotes() ([]Quote, error) {
	var allQuotes []Quote
	currentPage := s.baseURL

	for {
		quotes, err := s.scrapePage(currentPage)
		if err != nil {
			return nil, err
		}
		allQuotes = append(allQuotes, quotes...)

		// Check if there's a next page
		doc, err := goquery.NewDocument(currentPage)
		if err != nil {
			return nil, err
		}

		nextPage := doc.Find(".next > a").AttrOr("href", "")
		if nextPage == "" {
			break
		}
		currentPage = s.baseURL + nextPage
	}

	return allQuotes, nil
}

func main() {
	scraper := NewScraper()
	quotes, err := scraper.ScrapeAllQuotes()
	if err != nil {
		log.Fatalf("Error scraping quotes: %v", err)
	}

	// Save results to a JSON file
	file, err := os.Create("quotes.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(quotes); err != nil {
		log.Fatalf("Error encoding quotes: %v", err)
	}

	log.Printf("Successfully scraped %d quotes and saved to quotes.json", len(quotes))
}
