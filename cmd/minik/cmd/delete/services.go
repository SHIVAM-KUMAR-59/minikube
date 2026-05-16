package delete

import (
	"fmt"
	"log/slog"
	"net/http"

	minikcmd "github.com/SHIVAM-KUMAR-59/minikube/cmd/minik/cmd"
	"github.com/spf13/cobra"
)

var deleteService = &cobra.Command{
	Use:   "service <service-id>",
	Short: "Delete a service from the minik cluster",
	Long:  `Delete a service from the minik cluster by providing the service ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceID := args[0]

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/services/%s", serviceID), nil)
		if err != nil {
			slog.Error("Failed to create delete request", "error", err)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Failed to delete service", "error", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			fmt.Printf("\033[32m✓\033[0m Service \033[36m%s\033[0m deleted successfully.\n", serviceID)
		case http.StatusNotFound:
			fmt.Printf("\033[31m✗\033[0m Service \033[36m%s\033[0m not found.\n", serviceID)
		default:
			slog.Error("Failed to delete service", "status_code", resp.StatusCode)
		}
	},
}

func init() {
	minikcmd.Delete.AddCommand(deleteService)
}
