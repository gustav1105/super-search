package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

// Config represents the structure of the TOML configuration
type Config struct {
	Service struct {
		BaseURL  string `toml:"base_url"`
		Username string `toml:"username"`
		Password string `toml:"password"`
    SearchURL string `toml:"search_url"`
	} `toml:"service"`

	Endpoints struct {
		XMLTV        string `toml:"xmltv"`
		VODCategories string `toml:"vod_categories"`
		VODStreams    string `toml:"vod_streams"`
		VODInfo       string `toml:"vod_info"`
    AddVod        string  `toml:"add_vod"`
    QueryVods     string  `toml:"query_vods"`
	} `toml:"endpoints"`

	Logging struct {
		Level  string `toml:"level"`
		Output string `toml:"output"`
	} `toml:"logging"`

	Performance struct {
		TimeoutSeconds int `toml:"timeout_seconds"`
	} `toml:"performance"`
}

// LoadConfig loads the TOML configuration from a file
func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

