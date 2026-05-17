package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type ClusterHealthResponse struct {
	TotalPods     int    `json:"total_pods"`
	RunningPods   int    `json:"running_pods"`
	PendingPods   int    `json:"pending_pods"`
	TotalNodes    int    `json:"total_nodes"`
	ReadyNodes    int    `json:"ready_nodes"`
	TotalServices int    `json:"total_services"`
	ClusterHealth string `json:"cluster_health"`
}

var status = &cobra.Command{
	Use:   "status",
	Short: "Show overall cluster health",
	Long:  "Displays overall MiniK cluster health and statistics.",
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := http.Get(
			"http://localhost:8080/cluster/health",
		)

		if err != nil {
			slog.Error("Failed to fetch cluster health", "error", err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("\033[31m✗\033[0m Failed to fetch cluster health \033[90m(status=%d)\033[0m\n", resp.StatusCode)
			return
		}

		var result ClusterHealthResponse

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			slog.Error("Failed to decode response", "error", err)
			return
		}

		const (
			reset  = "\033[0m"
			bold   = "\033[1m"
			cyan   = "\033[36m"
			green  = "\033[32m"
			yellow = "\033[33m"
			red    = "\033[31m"
			gray   = "\033[90m"
		)

		healthColor := func(status string) string {

			switch status {

			case "HEALTHY":
				return green + bold + status + reset

			case "DEGRADED":
				return yellow + bold + status + reset

			default:
				return red + bold + status + reset
			}
		}

		fmt.Println()

		fmt.Printf("  %s%sMiniK Cluster Health%s\n", bold, cyan, reset)
		fmt.Printf("  %s%s%s\n", gray, strings.Repeat("─", 50), reset)
		fmt.Printf("  %sOverall Health:%s   %s\n", bold, reset, healthColor(result.ClusterHealth))

		fmt.Println()

		fmt.Printf("  %sPods%s\n", cyan, reset)
		fmt.Printf("    Total Pods:      %d\n", result.TotalPods)
		fmt.Printf("    Running Pods:    %d\n", result.RunningPods)
		fmt.Printf("    Pending Pods:    %d\n", result.PendingPods)

		fmt.Println()

		fmt.Printf("  %sNodes%s\n", cyan, reset)
		fmt.Printf("    Total Nodes:     %d\n", result.TotalNodes)
		fmt.Printf("    Ready Nodes:     %d\n", result.ReadyNodes)

		fmt.Println()

		fmt.Printf("  %sServices%s\n", cyan, reset)
		fmt.Printf("    Total Services:  %d\n", result.TotalServices)
		fmt.Printf("  %s%s%s\n", gray, strings.Repeat("─", 50), reset)

		fmt.Println()
	},
}

func init() {
	root.AddCommand(status)
}
