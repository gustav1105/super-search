package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gustav1105/epg_approved/internal/api"
  "net/http"
  "bytes"
  "time"
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

func (h *IPTVHandler) UpdateVOD() error {
  fmt.Println("Starting VOD update process...")

  // Fetch VOD categories
  categories, err := h.fetchVODCategories()
  if err != nil {
    return fmt.Errorf("error fetching VOD categories: %w", err)
  }

  // Process each category sequentially
  for _, category := range categories {
    // Validate and extract category_id
    categoryID, ok := category["category_id"].(string)
    if !ok {
      fmt.Printf("Invalid category_id format; skipping category: %v\n", category)
      continue
    }

    // Fetch streams for the category
    streams, err := h.fetchVODStreams(categoryID)
    if err != nil {
      fmt.Printf("Error fetching streams for category ID %s: %v\n", categoryID, err)
      continue
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
        fmt.Printf("Invalid stream_id format; skipping stream: %v\n", stream)
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

    // Log category processing success
    fmt.Printf("Successfully processed category ID %s\n", categoryID)
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

func (h *IPTVHandler) fetchVODStreams(categoryID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s&category_id=%s&username=%s&password=%s",
		h.Client.BaseURL, h.VODStreams, categoryID, h.Client.Username, h.Client.Password)

	var streams []map[string]interface{}
	retryAttempts := 3
	retryDelay := 3 * time.Second // Configurable delay between retries

	for retries := 0; retries < retryAttempts; retries++ {
		// Attempt the HTTP GET request
		resp, err := h.Client.HTTPClient.Get(url)
		if err != nil {
			fmt.Printf("Attempt %d/%d: Failed to fetch streams for category ID %s: %v\n", retries+1, retryAttempts, categoryID, err)
			if retries < retryAttempts-1 {
				time.Sleep(retryDelay)
			}
			continue
		}

		// Ensure response body is closed
		defer resp.Body.Close()

		// Check HTTP status code
		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("Attempt %d/%d: Received HTTP %d for category ID %s: %s\n",
				retries+1, retryAttempts, resp.StatusCode, categoryID, string(body))
			if retries < retryAttempts-1 {
				time.Sleep(retryDelay)
			}
			continue
		}

		// Decode the response body
		if err := json.NewDecoder(resp.Body).Decode(&streams); err != nil {
			fmt.Printf("Attempt %d/%d: Failed to parse response for category ID %s: %v\n", retries+1, retryAttempts, categoryID, err)
			if retries < retryAttempts-1 {
				time.Sleep(retryDelay)
			}
			continue
		}

		// Success: return the parsed streams
		return streams, nil
	}

	// All retries failed
	return nil, fmt.Errorf("failed to fetch VOD streams for category ID %s after %d attempts", categoryID, retryAttempts)
}

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

  if resp.StatusCode != http.StatusOK {
    // Log response body to debug unexpected status codes
    fmt.Printf("Unexpected HTTP status for VOD ID %s: %d. Response: %s\n", vodID, resp.StatusCode, string(body))
    return nil, fmt.Errorf("unexpected HTTP status %d for VOD ID %s", resp.StatusCode, vodID)
  }

  var info map[string]interface{}
  if err := json.Unmarshal(body, &info); err != nil {
    // Log raw body for debugging invalid JSON
    fmt.Printf("Failed to parse VOD info for ID %s. Response: %s\n", vodID, string(body))
    return nil, fmt.Errorf("failed to parse VOD info: %w", err)
  }

  return info, nil
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
    movie_image, _ := info["movie_image"].(string)
    youtube_trailer, _ := info["youtube_trailer"].(string)
    tmdb_id, _ := info["tmdb_id"].(string)

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
      "movie_image":  movie_image,
      "youtube_trailer": youtube_trailer,
      "tmbd_id":       tmdb_id,
    })
  }

  // Create payload with metadata
  payload, err := json.Marshal(map[string]interface{}{
    "metadata": vodData,
  })
  if err != nil {
    return fmt.Errorf("failed to marshal request payload: %w", err)
  }

  //DEBUG:
  //fmt.Printf("Payload: %s\n", string(payload))

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

  //fmt.Println("Successfully sent metadata to embed API.")
  return nil
}

