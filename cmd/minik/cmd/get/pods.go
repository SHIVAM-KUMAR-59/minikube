package get

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	minikcmd "github.com/SHIVAM-KUMAR-59/minikube/cmd/minik/cmd"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	"github.com/spf13/cobra"
)

var pods = &cobra.Command{
	Use:   "pods",
	Short: "Get pods from the minik cluster",
	Long:  `Get pods from the minik cluster to see the status of the pods running in the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		// HTTP GET request to the minik cluster to get the pods
		resp, err := http.Get("http://localhost:8080/pods")
		if err != nil {
			slog.Error("Failed to get pods from minik cluster", "error", err)
			return
		}
		defer resp.Body.Close()

		// Check if the response status code is 200 OK
		if resp.StatusCode == http.StatusOK {
			var result []store.Pod
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				slog.Error("Failed to decode response", "error", err)
				return
			}

			if len(result) == 0 {
				fmt.Println("No pods found.")
				return
			}

			// Print the pods in a tabular format with colors based on the status
			const (
				reset  = "\033[0m"
				bold   = "\033[1m"
				cyan   = "\033[36m"
				green  = "\033[32m"
				yellow = "\033[33m"
				red    = "\033[31m"
				gray   = "\033[90m"
			)

			statusColor := func(status string) string {
				padded := fmt.Sprintf("%-12s", status)
				switch status {
				case store.StatusRunning:
					return green + padded + reset
				case store.StatusScheduled, store.StatusPending:
					return yellow + padded + reset
				default:
					return gray + padded + reset
				}
			}

			fmt.Println()
			fmt.Printf("  %s%-36s  %-20s  %-20s  %-12s  %-10s%s\n",
				bold+cyan, "ID", "NAME", "IMAGE", "STATUS", "NODE", reset)
			fmt.Printf("  %s%s%s\n", gray,
				strings.Repeat("─", 108), reset)

			for _, pod := range result {
				fmt.Printf("  %-36s  %-20s  %-20s  %s  %-10s\n",
					pod.ID,
					pod.Name,
					pod.Image,
					statusColor(pod.Status),
					pod.NodeID,
				)
			}
			fmt.Println()
		} else {
			slog.Error("Failed to get pods from minik cluster", "status_code", resp.StatusCode)
		}
	},
}

func init() {
	minikcmd.Get.AddCommand(pods)
}
