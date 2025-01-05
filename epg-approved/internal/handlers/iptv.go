package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gustav1105/epg_approved/internal/api"
  "net/http"
  "bytes"
  "sync"
)

// IPTVHandler handles operations related to IPTV
type IPTVHandler struct {
	Client        *api.IPTVClient
	XMLTVEndpoint string
	VODCategories string
	VODStreams    string
	VODInfo       string
  AddVod        string
  QueryVods     string
}
// CategoryResult holds the result or error for a category
type CategoryResult struct {
	CategoryID   string
	CategoryName string
	Error        error
}

type Category struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	ParentID     int    `json:"parent_id"`
}

// NewIPTVHandler creates and returns a new IPTVHandler with the provided client and endpoints
func NewIPTVHandler(client *api.IPTVClient, xmltv, vodCategories, vodStreams, vodInfo, addVod, queryVods string) *IPTVHandler {
	return &IPTVHandler{
		Client:        client,
		XMLTVEndpoint: xmltv,
		VODCategories: vodCategories,
		VODStreams:    vodStreams,
		VODInfo:       vodInfo,
    AddVod:        addVod,
    QueryVods:     queryVods,
	}
}

// UpdateVOD orchestrates the fetching of VOD data using channels
func (h *IPTVHandler) UpdateVOD() error {
	fmt.Println("Starting VOD update process...")

	// Fetch VOD categories
	categories, err := h.fetchVODCategories()
	if err != nil {
		return fmt.Errorf("error fetching VOD categories: %w", err)
	}

	// Create a channel for category processing
	categoryChan := make(chan CategoryResult)

	// Use a WaitGroup to manage parallel processing
	var wg sync.WaitGroup

	// Process each category
	for _, category := range categories {
		categoryID, ok := category["category_id"].(string)
		if !ok {
			fmt.Println("Invalid category_id format; skipping category")
			continue
		}
		categoryName, ok := category["category_name"].(string)
		if !ok {
			fmt.Println("Invalid category_name format; skipping category")
			continue
		}

		wg.Add(1)
		go func(categoryID, categoryName string) {
			defer wg.Done()

			// Fetch streams for the category
			streams, err := h.fetchVODStreams(categoryID)
			if err != nil {
				categoryChan <- CategoryResult{CategoryID: categoryID, CategoryName: categoryName, Error: err}
				return
			}

			// Process each stream
			for _, stream := range streams {
				var vodID string
				switch id := stream["stream_id"].(type) {
				case string:
					vodID = id
				case float64:
					vodID = fmt.Sprintf("%.0f", id)
				default:
					fmt.Println("Invalid stream_id format; skipping stream")
					continue
				}

				// Fetch VOD info
				info, err := h.fetchVODInfo(vodID)
				if err != nil {
					fmt.Printf("Error fetching info for VOD %s: %v\n", vodID, err)
					continue
				}

				// Send VOD info to embed API
				err = h.AddVODInfo([]map[string]interface{}{info})
				if err != nil {
					fmt.Printf("Error sending VOD Info for VOD %s: %v\n", vodID, err)
				} else {
					fmt.Printf("Successfully sent data for VOD %s to Search service.\n", vodID)
				}
			}

			// Send success result to the channel
			categoryChan <- CategoryResult{CategoryID: categoryID, CategoryName: categoryName, Error: nil}
		}(categoryID, categoryName)
	}

	// Close the channel after all goroutines complete
	go func() {
		wg.Wait()
		close(categoryChan)
	}()

	// Collect results
	for result := range categoryChan {
		if result.Error != nil {
			fmt.Printf("Error processing category %s (%s): %v\n", result.CategoryName, result.CategoryID, result.Error)
		} else {
			fmt.Printf("Successfully processed category %s (%s)\n", result.CategoryName, result.CategoryID)
		}
	}

	fmt.Println("VOD update process completed successfully.")
	return nil
}

// fetchVODCategories fetches VOD categories from the server
func (h *IPTVHandler) fetchVODCategories() ([]map[string]interface{}, error) {
  url := fmt.Sprintf("%s/%s&username=%s&password=%s", h.Client.BaseURL, h.VODCategories, h.Client.Username, h.Client.Password)
  resp, err := h.Client.HTTPClient.Get(url)
  if err != nil {
    return nil, fmt.Errorf("failed to fetch VOD categories: %w", err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  var categories []map[string]interface{}
  if err := json.Unmarshal(body, &categories); err != nil {
    return nil, fmt.Errorf("failed to parse categories: %w", err)
  }
  return categories, nil
}

// fetchVODStreams fetches VOD streams for a specific category
func (h *IPTVHandler) fetchVODStreams(categoryID string) ([]map[string]interface{}, error) {
  url := fmt.Sprintf("%s/%s&category_id=%s&username=%s&password=%s", h.Client.BaseURL, h.VODStreams, categoryID, h.Client.Username, h.Client.Password)

  resp, err := h.Client.HTTPClient.Get(url)
  if err != nil {
    return nil, fmt.Errorf("failed to fetch VOD streams for category ID %s: %w", categoryID, err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  var streams []map[string]interface{}
  if err := json.Unmarshal(body, &streams); err != nil {
    return nil, fmt.Errorf("failed to parse streams: %w", err)
  }

  return streams, nil
}

// fetchVODInfo fetches detailed VOD info for a specific stream
func (h *IPTVHandler) fetchVODInfo(vodID string) (map[string]interface{}, error) {
  url := fmt.Sprintf("%s/%s&vod_id=%s&username=%s&password=%s", h.Client.BaseURL, h.VODInfo, vodID, h.Client.Username, h.Client.Password)

  resp, err := h.Client.HTTPClient.Get(url)
  if err != nil {
    return nil, fmt.Errorf("failed to fetch VOD info for ID %s: %w", vodID, err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  var info map[string]interface{}
  if err := json.Unmarshal(body, &info); err != nil {
    return nil, fmt.Errorf("failed to parse VOD info: %w", err)
  }
  return info, nil
}

func (h *IPTVHandler) Embed(sentences []string) error {
  url := fmt.Sprintf("%s/%s", h.Client.SearchURL, h.AddVod)

  payload, err := json.Marshal(map[string]interface{}{
    "sentences": sentences,
  })
  if err != nil {
    return fmt.Errorf("failed to marshal request payload: %w", err)
  }

  resp, err := h.Client.HTTPClient.Post(url, "application/json", bytes.NewBuffer(payload))
  if err != nil {
    return fmt.Errorf("failed to send request to embed API: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
  }

  return nil
}

func (h *IPTVHandler) AddVODInfo(vodInfos []map[string]interface{}) error {
  var vodData []map[string]interface{}

  for _, vodInfo := range vodInfos {
    // Access the "info" section
    info, ok := vodInfo["info"].(map[string]interface{})
    if !ok {
      fmt.Println("Error: Missing or invalid 'info' field")
      continue
    }

    // Access the "movie_data" section
    movieData, ok := vodInfo["movie_data"].(map[string]interface{})
    if !ok {
      fmt.Println("Error: Missing or invalid 'movie_data' field")
      continue
    }

    // Safely extract fields with type checks
    streamID := fmt.Sprintf("%v", movieData["stream_id"]) // Ensure string conversion
    title, _ := movieData["name"].(string)
    plot, _ := info["plot"].(string)
    genre, _ := info["genre"].(string)
    releaseDate, _ := info["releasedate"].(string)
    rating := fmt.Sprintf("%v", info["rating"]) // Convert numeric ratings to string
    director, _ := info["director"].(string)
    cast, _ := info["cast"].(string)

    // Add metadata for retrieval
    vodData = append(vodData, map[string]interface{}{
      "stream_id":    streamID,
      "title":        title,
      "plot":         plot,
      "genre":        genre,
      "release_date": releaseDate,
      "rating":       rating,
      "director":     director,
      "cast":         cast,
    })
  }

  // Create payload with metadata
  payload, err := json.Marshal(map[string]interface{}{
    "metadata": vodData,
  })
  if err != nil {
    return fmt.Errorf("failed to marshal request payload: %w", err)
  }

  fmt.Printf("Payload: %s\n", string(payload))

  // Send to embed API
  url := fmt.Sprintf("%s/add", h.Client.SearchURL)
  resp, err := h.Client.HTTPClient.Post(url, "application/json", bytes.NewBuffer(payload))
  if err != nil {
    return fmt.Errorf("failed to send request to embed API: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    body, _ := ioutil.ReadAll(resp.Body)
    return fmt.Errorf("non-OK HTTP status: %s, response: %s", resp.Status, string(body))
  }

  fmt.Println("Successfully sent metadata to embed API.")
  return nil
}

