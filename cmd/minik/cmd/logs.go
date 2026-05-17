package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var logs = &cobra.Command{
	Use:   "logs",
	Short: "Fetch the logs",
	Long:  "Fetch the docker logs and streams back to the client",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("\033[31m✗\033[0m Pod ID is required\n")
			fmt.Printf("  \033[90mUsage: minik logs <pod-id>\033[0m\n")
			return
		}

		podName := args[0]

		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/pods/%s/logs", podName))
		if err != nil {
			fmt.Printf("\033[31m✗\033[0m Could not reach server at \033[36mlocalhost:8080\033[0m: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("\033[31m✗\033[0m Failed to fetch logs (HTTP \033[36m%d\033[0m)\n", resp.StatusCode)
			return
		}

		fmt.Println()
		fmt.Printf("  \033[1m\033[36mLogs\033[0m  \033[90mpod/%s\033[0m\n", podName)
		fmt.Printf("  \033[90m%s\033[0m\n\n", strings.Repeat("─", 60))

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			upper := strings.ToUpper(line)

			switch {
			case strings.Contains(upper, "[ERROR]") ||
				strings.Contains(upper, "[FATAL]") ||
				strings.Contains(upper, "ERROR") ||
				strings.Contains(upper, "FATAL"):
				fmt.Printf("  \033[31m●\033[0m \033[31m%s\033[0m\n", line)

			case strings.Contains(upper, "[WARN]") ||
				strings.Contains(upper, "WARN"):
				fmt.Printf("  \033[33m●\033[0m \033[33m%s\033[0m\n", line)

			case strings.Contains(upper, "[NOTICE]") ||
				strings.Contains(upper, "NOTICE"):
				fmt.Printf("  \033[36m●\033[0m \033[36m%s\033[0m\n", line)

			case strings.Contains(upper, "[INFO]") ||
				strings.Contains(upper, "INFO"):
				fmt.Printf("  \033[32m●\033[0m \033[32m%s\033[0m\n", line)

			case strings.Contains(line, "docker-entrypoint"):
				fmt.Printf("  \033[90m  %s\033[0m\n", line)

			default:
				fmt.Printf("  \033[37m  %s\033[0m\n", line)
			}
		}
		fmt.Println()
	},
}

func init() {
	root.AddCommand(logs)
}
