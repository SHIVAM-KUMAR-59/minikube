package cluster

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	minikcmd "github.com/SHIVAM-KUMAR-59/minikube/cmd/minik/cmd"
	"github.com/spf13/cobra"
)

var numberOfWorkers int

var startCluster = &cobra.Command{
	Use:   "start",
	Short: "Start the minik cluster",
	Long:  "Takes a number N and starts the server and N workers as background processes (default 2)",
	Run: func(cmd *cobra.Command, args []string) {
		pidDir := filepath.Join(os.Getenv("HOME"), ".minik")
		pidFile := filepath.Join(pidDir, "cluster.pid")

		// Create ~/.minik directory if it doesn't exist
		if err := os.MkdirAll(pidDir, 0755); err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to create config directory: %v\n", err)
			return
		}

		var pids []string

		fmt.Println()
		fmt.Printf("  \033[1m\033[36mStarting minik cluster...\033[0m\n")
		fmt.Printf("  \033[90m%s\033[0m\n", strings.Repeat("─", 40))

		// Get current executable path
		execPath, err := os.Executable()
		if err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to locate executable path: %v\n", err)
			return
		}

		// Directory where minik binary exists
		binDir := filepath.Dir(execPath)

		// Resolve sibling binaries
		serverBinary := filepath.Join(binDir, "minik-server")
		workerBinary := filepath.Join(binDir, "minik-worker")

		// Start the server
		server := exec.Command(serverBinary)
		server.Stdout = nil
		server.Stderr = nil
		server.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		if err := server.Start(); err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to start server: %v\n", err)
			return
		}
		pids = append(pids, strconv.Itoa(server.Process.Pid))
		fmt.Printf("  \033[32m✓\033[0m Server started        \033[90mpid=%d\033[0m\n", server.Process.Pid)

		// Start N workers
		for i := 1; i <= numberOfWorkers; i++ {
			nodeID := fmt.Sprintf("node%d", i)
			worker := exec.Command(workerBinary, "--node-id", nodeID)
			worker.Stdout = nil
			worker.Stderr = nil
			worker.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
			if err := worker.Start(); err != nil {
				fmt.Printf("  \033[31m✗\033[0m Failed to start worker \033[36m%s\033[0m: %v\n", nodeID, err)
				continue
			}
			pids = append(pids, strconv.Itoa(worker.Process.Pid))
			fmt.Printf("  \033[32m✓\033[0m Worker started        \033[90mpid=%-6d  node=%s\033[0m\n", worker.Process.Pid, nodeID)
		}

		// Save all PIDs to ~/.minik/cluster.pid
		if err := os.WriteFile(pidFile, []byte(strings.Join(pids, "\n")), 0644); err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to save PIDs: %v\n", err)
			return
		}

		fmt.Printf("  \033[90m%s\033[0m\n", strings.Repeat("─", 40))
		fmt.Printf("  \033[32m\033[1m✓ Cluster started\033[0m  \033[90m1 server, %d worker(s)\033[0m\n", numberOfWorkers)
		fmt.Println()
	},
}

func init() {
	startCluster.Flags().IntVarP(&numberOfWorkers, "workers", "w", 2, "Number of workers to start")
	minikcmd.Cluster.AddCommand(startCluster)
}
