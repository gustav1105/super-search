package main

import (
	"fmt"
	"github.com/gustav1105/epg_approved/internal/api"
	"github.com/gustav1105/epg_approved/internal/config"
	"github.com/gustav1105/epg_approved/internal/handlers"
	"github.com/gustav1105/epg_approved/cmd"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	// Initialize the API client
	iptvClient := api.NewIPTVClient(
		cfg.Service.BaseURL,
		cfg.Service.Username,
		cfg.Service.Password,
    cfg.Service.SearchURL,
		cfg.Performance.TimeoutSeconds,
	)

	// Initialize the IPTV handler
	iptvHandler := handlers.NewIPTVHandler(
		iptvClient,
		cfg.Endpoints.XMLTV,
		cfg.Endpoints.VODCategories,
		cfg.Endpoints.VODStreams,
		cfg.Endpoints.VODInfo,
    cfg.Endpoints.AddVod,
    cfg.Endpoints.QueryVods,
	)

	// Inject the handler into commands
	cmd.SetHandler(iptvHandler)

	// Initialize and execute the root command
	rootCmd := cmd.InitializeRootCmd()
	rootCmd.Execute()
}

