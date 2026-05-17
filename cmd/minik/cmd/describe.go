package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	"github.com/spf13/cobra"
)

var describe = &cobra.Command{
	Use:   "describe",
	Short: "Get information about the pod",
	Long:  "Get detailed information about a pod by its name",
	Run: func(cmd *cobra.Command, args []string) {

		// Extract the pod name
		if len(args) == 0 {
			fmt.Printf("\033[31m✗\033[0m Pod name is required\n")
			fmt.Printf("  \033[90mUsage: minik describe <pod-name>\033[0m\n")
			return
		}

		podName := args[0]

		// Make HTTP request
		url := fmt.Sprintf(
			"http://localhost:8080/pods/%s",
			podName,
		)

		resp, err := http.Get(url)
		if err != nil {
			slog.Error("Failed to fetch pod details", "pod_name", podName, "error", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("\033[31m✗\033[0m Failed to fetch pod details \033[90m(status=%d)\033[0m\n", resp.StatusCode)
			return
		}

		var pod store.Pod

		if err := json.NewDecoder(resp.Body).Decode(&pod); err != nil {
			slog.Error("Failed to decode pod response", "error", err)
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

		statusColor := func(status string) string {
			switch status {
			case store.StatusRunning:
				return green + status + reset

			case store.StatusScheduled, store.StatusPending:
				return yellow + status + reset

			default:
				return red + status + reset
			}
		}

		fmt.Println()

		fmt.Printf("  %s%sPod Description%s\n", bold, cyan, reset)
		fmt.Printf("  %s%s%s\n", gray, strings.Repeat("─", 50), reset)
		fmt.Printf("  %sName:%s       %s\n", bold, reset, pod.Name)
		fmt.Printf("  %sID:%s         %s\n", bold, reset, pod.ID)
		fmt.Printf("  %sImage:%s      %s\n", bold, reset, pod.Image)
		fmt.Printf("  %sStatus:%s     %s\n", bold, reset, statusColor(pod.Status))
		fmt.Printf("  %sNode:%s       %s\n", bold, reset, pod.NodeID)
		fmt.Printf("  %sReplicas:%s   %d\n", bold, reset, pod.Replicas)
		fmt.Printf("  %s%s%s\n", gray, strings.Repeat("─", 50), reset)

		fmt.Println()
	},
}

func init() {
	root.AddCommand(describe)
}
