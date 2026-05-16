package cluster

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	minikcmd "github.com/SHIVAM-KUMAR-59/minikube/cmd/minik/cmd"
	"github.com/spf13/cobra"
)

var stopCluster = &cobra.Command{
	Use:   "stop",
	Short: "Stop the minik cluster",
	Long:  "Reads the PID file and stops all running server and worker processes.",
	Run: func(cmd *cobra.Command, args []string) {
		pidFile := filepath.Join(os.Getenv("HOME"), ".minik", "cluster.pid")

		// Read PID file
		data, err := os.ReadFile(pidFile)
		if err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to read PID file \033[36m%s\033[0m: %v\n", pidFile, err)
			fmt.Printf("  \033[90mIs the cluster running?\033[0m\n")
			return
		}

		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		if len(lines) == 0 {
			fmt.Printf("\033[31m✗\033[0m No processes found in PID file.\n")
			return
		}

		fmt.Println()
		fmt.Printf("  \033[1m\033[36mStopping minik cluster...\033[0m\n")
		fmt.Printf("  \033[90m%s\033[0m\n", strings.Repeat("─", 40))

		killed := 0
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			pid, err := strconv.Atoi(line)
			if err != nil {
				fmt.Printf("  \033[31m✗\033[0m Invalid PID \033[36m%s\033[0m: %v\n", line, err)
				continue
			}

			label := fmt.Sprintf("node%d", i)
			if i == 0 {
				label = "server"
			}

			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				fmt.Printf("  \033[31m✗\033[0m Failed to stop \033[36m%s\033[0m  \033[90mpid=%d: %v\033[0m\n", label, pid, err)
				continue
			}

			fmt.Printf("  \033[32m✓\033[0m Stopped %-10s  \033[90mpid=%d\033[0m\n", label, pid)
			killed++
		}

		// Delete PID file
		if err := os.Remove(pidFile); err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to delete PID file: %v\n", err)
			return
		}

		fmt.Printf("  \033[90m%s\033[0m\n", strings.Repeat("─", 40))
		fmt.Printf("  \033[32m\033[1m✓ Cluster stopped\033[0m  \033[90m%d process(es) terminated\033[0m\n", killed)
		fmt.Println()
	},
}

func init() {
	minikcmd.Cluster.AddCommand(stopCluster)
}
