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

var nodes = &cobra.Command{
	Use:   "nodes",
	Short: "Get nodes from the minik cluster",
	Long:  `Get nodes from the minik cluster to see the status of the nodes running in the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("http://localhost:8080/nodes")
		if err != nil {
			slog.Error("Failed to get nodes from minik cluster", "error", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result []store.Node
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				slog.Error("Failed to decode response", "error", err)
				return
			}

			if len(result) == 0 {
				fmt.Println("No nodes found.")
				return
			}

			const (
				reset = "\033[0m"
				bold  = "\033[1m"
				cyan  = "\033[36m"
				green = "\033[32m"
				gray  = "\033[90m"
				red   = "\033[31m"
			)

			statusColor := func(status string) string {
				padded := fmt.Sprintf("%-12s", status)
				switch status {
				case store.NodeStatusReady:
					return green + padded + reset
				case store.NodeStatusNotReady:
					return red + padded + reset
				default:
					return gray + padded + reset
				}
			}

			fmt.Println()
			fmt.Printf("  %s%-36s  %-20s  %-28s  %-12s%s\n",
				bold+cyan, "ID", "NAME", "LAST HEARTBEAT", "STATUS", reset)
			fmt.Printf("  %s%s%s\n", gray, strings.Repeat("─", 102), reset)

			for _, node := range result {
				fmt.Printf("  %-36s  %-20s  %-28s  %s\n",
					node.ID,
					node.Name,
					node.LastHeartbeat.Format("2006-01-02 15:04:05"),
					statusColor(node.Status),
				)
			}
			fmt.Println()
		} else {
			slog.Error("Failed to get nodes from minik cluster", "status_code", resp.StatusCode)
		}
	},
}

func init() {
	minikcmd.Get.AddCommand(nodes)
}
