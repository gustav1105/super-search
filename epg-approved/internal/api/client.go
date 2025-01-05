package api

import (
	"net/http"
	"time"
)

type IPTVClient struct {
	HTTPClient *http.Client
	BaseURL    string
	Username   string
	Password   string
  SearchURL  string
}

// NewClient creates and returns a new API Client with authentication details
func NewIPTVClient(baseURL, username, password, searchURL string, timeoutSeconds int) *IPTVClient {
	return &IPTVClient{
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
		BaseURL:  baseURL,
		Username: username,
		Password: password,
    SearchURL: searchURL,
	}
}

