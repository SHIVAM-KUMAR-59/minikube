package delete

import (
	"fmt"
	"log/slog"
	"net/http"

	minikcmd "github.com/SHIVAM-KUMAR-59/minikube/cmd/minik/cmd"
	"github.com/spf13/cobra"
)

var deletePod = &cobra.Command{
	Use:   "pod <pod-id>",
	Short: "Delete a pod from the minik cluster",
	Long:  `Delete a pod from the minik cluster by providing the pod ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podID := args[0]

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/pods/%s", podID), nil)
		if err != nil {
			slog.Error("Failed to create delete request", "error", err)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Failed to delete pod", "error", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			fmt.Printf("\033[32m✓\033[0m Pod \033[36m%s\033[0m deleted successfully.\n", podID)
		case http.StatusNotFound:
			fmt.Printf("\033[31m✗\033[0m Pod \033[36m%s\033[0m not found.\n", podID)
		default:
			slog.Error("Failed to delete pod", "status_code", resp.StatusCode)
		}
	},
}

func init() {
	minikcmd.Delete.AddCommand(deletePod)
}
