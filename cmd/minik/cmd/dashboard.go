package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Start the MiniK dashboard",
	Long:  "Starts the frontend dashboard server and opens it automatically in the browser.",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println()
		fmt.Printf("  \033[1m\033[36mStarting MiniK Dashboard...\033[0m\n")
		fmt.Printf("  \033[90m%s\033[0m\n", strings.Repeat("─", 40))

		// Start frontend process
		frontend := exec.Command("npm", "run", "dev")
		frontend.Dir = "./dashboard"

		// Attach logs
		frontend.Stdout = os.Stdout
		frontend.Stderr = os.Stderr

		// Start frontend
		if err := frontend.Start(); err != nil {
			fmt.Printf("  \033[31m✗\033[0m Failed to start dashboard: %v\n", err)
			return
		}

		fmt.Printf(
			"  \033[32m✓\033[0m Frontend server started  \033[90mpid=%d\033[0m\n",
			frontend.Process.Pid,
		)

		fmt.Printf("  \033[33m⟳\033[0m Waiting for dashboard to boot...\n")

		// Wait for dev server
		time.Sleep(3 * time.Second)

		url := "http://localhost:3000"

		// Cross-platform browser opener
		var browserCmd *exec.Cmd

		switch runtime.GOOS {
		case "darwin":
			browserCmd = exec.Command("open", url)

		case "windows":
			browserCmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)

		case "linux":
			browserCmd = exec.Command("xdg-open", url)

		default:
			fmt.Printf(
				"  \033[33m!\033[0m Unsupported OS. Open manually: \033[36m%s\033[0m\n",
				url,
			)
		}

		// Open browser if supported
		if browserCmd != nil {
			if err := browserCmd.Start(); err != nil {
				fmt.Printf(
					"  \033[31m✗\033[0m Failed to open browser automatically\n",
				)
				fmt.Printf(
					"  \033[33m→\033[0m Open manually: \033[36m%s\033[0m\n",
					url,
				)
			} else {
				fmt.Printf(
					"  \033[32m✓\033[0m Browser opened at  \033[36m%s\033[0m\n",
					url,
				)
			}
		}

		fmt.Printf("  \033[90m%s\033[0m\n", strings.Repeat("─", 40))
		fmt.Printf("  \033[32m\033[1m✓ Dashboard ready\033[0m\n")
		fmt.Println()

		// Keep process alive
		if err := frontend.Wait(); err != nil {
			fmt.Printf(
				"  \033[31m✗\033[0m Dashboard process exited: %v\n",
				err,
			)
		}
	},
}

func init() {
	root.AddCommand(dashboardCmd)
}
