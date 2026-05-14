package cmd

import (
	"encoding/json"
	"fmt"
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
			fmt.Println("Error pinging the minik cluster:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result map[string]string
			json.NewDecoder(resp.Body).Decode(&result)
			fmt.Println(result["message"])
		} else {
			fmt.Println("Minik cluster is not running. Status code:", resp.StatusCode)
		}
	},
}

func init() {
	root.AddCommand(ping)
}