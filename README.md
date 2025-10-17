A simple Go web scraper that collects quotes, authors, and tags from [quotes.toscrape.com](http://quotes.toscrape.com) and saves them to a JSON file.

Features
- Scrapes all quotes, authors, and tags from all pages
- Handles pagination automatically
- Saves results to `quotes.json` in structured JSON format
- Respects the website with rate limiting (1 second between requests)
- Uses Go best practices and error handling
