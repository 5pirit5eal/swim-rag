package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

// ScrapeHandler handles the HTTP request for scraping a URL.
// It extracts the URL from the request, scrapes the data, and responds with a success message.
func (rs *RAGService) ScrapeHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Getting scraping request...")
	// Parse the URL from the request
	url := req.URL.Query().Get("url")
	if url == "" {
		logger.Error("Missing url parameter")
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}
	httplog.LogEntrySetField(req.Context(), "url", slog.StringValue(url))

	// Scrape the URL
	err := rs.db.ScrapeURL(req.Context(), url)
	if err != nil {
		logger.Error("Failed to scrape URL", httplog.ErrAttr(err))
		http.Error(w, fmt.Errorf("failed to scrape URL %s: %w", url, err).Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Scraping completed successfully"))
}
