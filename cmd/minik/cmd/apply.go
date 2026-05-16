package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type PodYamlInput struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

var filePath string

var apply = &cobra.Command{
	Use:   "apply",
	Short: "Apply the .yaml file",
	Long:  "Read the .yaml file from the command and create the pod using the data from the file",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to open file \033[36m%s\033[0m: %v\n", filePath, err)
			return
		}

		podYamlInput := PodYamlInput{}
		if err := yaml.Unmarshal(data, &podYamlInput); err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to parse yaml: %v\n", err)
			return
		}

		if podYamlInput.Name == "" || podYamlInput.Image == "" {
			fmt.Printf("\033[31m✗\033[0m Invalid yaml: \033[36mname\033[0m and \033[36mimage\033[0m are required\n")
			return
		}

		resp, err := http.Post("http://localhost:8080/pods", "application/json", strings.NewReader(fmt.Sprintf(`{"name": "%s", "image": "%s"}`, podYamlInput.Name, podYamlInput.Image)))
		if err != nil {
			fmt.Printf("\033[31m✗\033[0m Could not reach server at \033[36mlocalhost:8080\033[0m: %v\n", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			fmt.Printf("\033[32m✓\033[0m Pod \033[36m%s\033[0m created successfully.\n", podYamlInput.Name)
		default:
			fmt.Printf("\033[31m✗\033[0m Server returned unexpected status: \033[36m%d\033[0m\n", resp.StatusCode)
		}
	},
}

func init() {
	apply.Flags().StringVarP(&filePath, "file", "f", "pod.yaml", "Path to the yaml file")
	root.AddCommand(apply)
}
