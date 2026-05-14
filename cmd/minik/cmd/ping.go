package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"
)

var ping = &cobra.Command{
	Use:   "ping",
	Short: "Ping the minik cluster",
	Long:  `Ping the minik cluster to check if it is running`,
	Run: func(cmd *cobra.Command, args []string) {
		// HTTP GET request to the minik cluster to check if it is running
		resp, err := http.Get("http://localhost:8080/ping")
		if err != nil {
			slog.Error("Failed to ping minik cluster", "error", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result map[string]string
			json.NewDecoder(resp.Body).Decode(&result)
			slog.Info(result["message"])
		} else {
			slog.Error("Minik cluster is not running", "status_code", resp.StatusCode)
		}
	},
}

func init() {
	root.AddCommand(ping)
}