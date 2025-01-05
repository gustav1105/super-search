package cmd

import (
	"fmt"
	"github.com/gustav1105/epg_approved/internal/handlers"

	"github.com/spf13/cobra"
)

// updateVODCmd represents the "update vod" command
var updateVODCmd = &cobra.Command{
	Use:   "update vod",
	Short: "Update VOD categories, streams, and info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the VOD update process...")
		
		// Call the UpdateVOD method on the handler
		if err := handler.UpdateVOD(); err != nil {
			fmt.Println("Error during VOD update:", err)
		} else {
			fmt.Println("VOD update completed successfully.")
		}
	},
}

// SetHandler allows the main function to inject the handler
var handler *handlers.IPTVHandler
func SetHandler(h *handlers.IPTVHandler) {
	handler = h
}

