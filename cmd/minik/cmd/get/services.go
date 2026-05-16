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

var services = &cobra.Command{
	Use:   "services",
	Short: "Get services from the minik cluster",
	Long:  `Get services from the minik cluster to see the services running in the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("http://localhost:8080/services")
		if err != nil {
			slog.Error("Failed to get services from minik cluster", "error", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result []store.Service
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				slog.Error("Failed to decode response", "error", err)
				return
			}

			if len(result) == 0 {
				fmt.Println("No services found.")
				return
			}

			const (
				reset = "\033[0m"
				bold  = "\033[1m"
				cyan  = "\033[36m"
				gray  = "\033[90m"
			)

			fmt.Println()
			fmt.Printf("  %s%-36s  %-20s  %-10s  %-40s%s\n",
				bold+cyan, "ID", "NAME", "PORT", "PODS", reset)
			fmt.Printf("  %s%s%s\n", gray, strings.Repeat("─", 112), reset)

			for _, svc := range result {
				pods := strings.Join(svc.Pods, ", ")
				if len(pods) == 0 {
					pods = gray + "none" + reset
				}
				fmt.Printf("  %-36s  %-20s  %-10s  %-40s\n",
					svc.ID,
					svc.Name,
					svc.Port,
					pods,
				)
			}
			fmt.Println()
		} else {
			slog.Error("Failed to get services from minik cluster", "status_code", resp.StatusCode)
		}
	},
}

func init() {
	minikcmd.Get.AddCommand(services)
}
