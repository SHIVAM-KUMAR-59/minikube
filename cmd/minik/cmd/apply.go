package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ClusterYaml struct {
	Cluster  ClusterConfig   `yaml:"cluster"`
	Pods     []PodConfig     `yaml:"pods"`
	Services []ServiceConfig `yaml:"services"`
}

type ClusterConfig struct {
	Workers int `yaml:"workers"`
}

type PodConfig struct {
	Name     string `yaml:"name"`
	Image    string `yaml:"image"`
	Replicas int    `yaml:"replicas"`
}

type ServiceConfig struct {
	Name string   `yaml:"name"`
	Port string   `yaml:"port"`
	Pods []string `yaml:"pods"`
}

type PodYamlInput struct {
	Name     string `yaml:"name"`
	Image    string `yaml:"image"`
	Replicas int    `yaml:"replicas"`
}

var filePath string

func handleClusterYaml(clusterYaml ClusterYaml) {
	// Check if cluster is running
	resp, err := http.Get("http://localhost:8080/ping")
	clusterRunning := err == nil && resp.StatusCode == http.StatusOK
	if resp != nil {
		resp.Body.Close()
	}

	if !clusterRunning && clusterYaml.Cluster.Workers > 0 {
		fmt.Printf("\033[36m→\033[0m Starting cluster with \033[36m%d\033[0m workers...\n", clusterYaml.Cluster.Workers)

		execPath, err := os.Executable()
		if err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to locate executable: %v\n", err)
			return
		}

		cmd := exec.Command(execPath, "cluster", "start", "--workers", fmt.Sprintf("%d", clusterYaml.Cluster.Workers))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to start cluster: %v\n", err)
			return
		}
	} else if clusterRunning {
		fmt.Printf("\033[32m✓\033[0m Cluster already running\n")
	}

	// Wait for workers to be ready
	fmt.Printf("\033[36m→\033[0m Waiting for \033[36m%d\033[0m node(s) to be ready\n", clusterYaml.Cluster.Workers)

	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinIdx := 0
	maxRetries := 30

	for i := 0; i < maxRetries; i++ {

		// Smooth spinner animation
		for j := 0; j < 10; j++ {
			fmt.Printf("\r  \033[36m%s\033[0m Waiting for nodes...", spinner[spinIdx])
			spinIdx = (spinIdx + 1) % len(spinner)
			time.Sleep(80 * time.Millisecond)
		}

		resp, err := http.Get("http://localhost:8080/nodes")
		if err != nil {
			continue
		}

		var nodes []struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
			resp.Body.Close()
			continue
		}

		resp.Body.Close()

		readyCount := 0
		for _, n := range nodes {
			if n.Status == "READY" {
				readyCount++
			}
		}

		if readyCount >= clusterYaml.Cluster.Workers {
			fmt.Printf("\r  \033[32m✓\033[0m %d node(s) ready          \n", readyCount)
			break
		}

		if i == maxRetries-1 {
			fmt.Printf("\r  \033[31m✗\033[0m Timed out waiting for nodes\n")
			return
		}
	}

	// Create pods
	podsCreated := 0
	if len(clusterYaml.Pods) > 0 {
		fmt.Println()
		fmt.Printf("\033[1mCreating pods...\033[0m\n")

		for _, pod := range clusterYaml.Pods {
			if pod.Name == "" || pod.Image == "" {
				fmt.Printf("\033[31m✗\033[0m Skipping pod — name and image are required\n")
				continue
			}

			replicas := pod.Replicas
			if replicas <= 0 {
				replicas = 1
			}

			body := fmt.Sprintf(`{"name": "%s", "image": "%s", "replicas": %d}`, pod.Name, pod.Image, replicas)
			resp, err := http.Post("http://localhost:8080/pods", "application/json", strings.NewReader(body))
			if err != nil {
				fmt.Printf("\033[31m✗\033[0m Could not reach server: %v\n", err)
				return
			}
			resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK, http.StatusCreated:
				fmt.Printf("  \033[32m✓\033[0m Pod \033[36m%s\033[0m created (%d replica(s))\n", pod.Name, replicas)
				podsCreated += replicas
			default:
				fmt.Printf("  \033[31m✗\033[0m Failed to create pod \033[36m%s\033[0m (HTTP %d)\n", pod.Name, resp.StatusCode)
			}
		}
	}

	// Create services
	servicesCreated := 0
	if len(clusterYaml.Services) > 0 {
		fmt.Println()
		fmt.Printf("\033[1mCreating services...\033[0m\n")

		for _, svc := range clusterYaml.Services {
			if svc.Name == "" || svc.Port == "" {
				fmt.Printf("\033[31m✗\033[0m Skipping service — name and port are required\n")
				continue
			}

			podsJSON, _ := json.Marshal(svc.Pods)
			body := fmt.Sprintf(`{"name": "%s", "port": "%s", "pods": %s}`, svc.Name, svc.Port, string(podsJSON))
			resp, err := http.Post("http://localhost:8080/services", "application/json", strings.NewReader(body))
			if err != nil {
				fmt.Printf("\033[31m✗\033[0m Could not reach server: %v\n", err)
				return
			}
			resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK, http.StatusCreated:
				fmt.Printf("  \033[32m✓\033[0m Service \033[36m%s\033[0m created on port \033[36m%s\033[0m\n", svc.Name, svc.Port)
				servicesCreated++
			default:
				fmt.Printf("  \033[31m✗\033[0m Failed to create service \033[36m%s\033[0m (HTTP %d)\n", svc.Name, resp.StatusCode)
			}
		}
	}

	// Summary
	fmt.Println()
	fmt.Printf("\033[90m────────────────────────────────────────\033[0m\n")
	fmt.Printf("\033[32m\033[1m✓ Applied successfully\033[0m\n")
	fmt.Printf("  \033[90m%d pod(s) created, %d service(s) created\033[0m\n", podsCreated, servicesCreated)
	fmt.Println()
}

func handlePodYaml(data []byte) {
	podYamlInput := PodYamlInput{}
	if err := yaml.Unmarshal(data, &podYamlInput); err != nil {
		fmt.Printf("\033[31m✗\033[0m Failed to parse yaml: %v\n", err)
		return
	}

	if podYamlInput.Name == "" || podYamlInput.Image == "" {
		fmt.Printf("\033[31m✗\033[0m Invalid yaml: \033[36mname\033[0m and \033[36mimage\033[0m are required\n")
		return
	}

	replicas := podYamlInput.Replicas
	if replicas <= 0 {
		replicas = 1
	}

	body := fmt.Sprintf(`{"name": "%s", "image": "%s", "replicas": %d}`, podYamlInput.Name, podYamlInput.Image, replicas)
	resp, err := http.Post("http://localhost:8080/pods", "application/json", strings.NewReader(body))
	if err != nil {
		fmt.Printf("\033[31m✗\033[0m Could not reach server at \033[36mlocalhost:8080\033[0m: %v\n", err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		fmt.Printf("\033[32m✓\033[0m Pod \033[36m%s\033[0m created successfully (%d replica(s)).\n", podYamlInput.Name, replicas)
	default:
		fmt.Printf("\033[31m✗\033[0m Server returned unexpected status: \033[36m%d\033[0m\n", resp.StatusCode)
	}
}

var apply = &cobra.Command{
	Use:   "apply",
	Short: "Apply the .yaml file",
	Long:  "Read the .yaml file from the command and create resources defined in the file",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("\033[31m✗\033[0m Failed to open file \033[36m%s\033[0m: %v\n", filePath, err)
			return
		}

		var clusterYaml ClusterYaml
		yaml.Unmarshal(data, &clusterYaml)

		if len(clusterYaml.Pods) > 0 || clusterYaml.Cluster.Workers > 0 || len(clusterYaml.Services) > 0 {
			handleClusterYaml(clusterYaml)
		} else {
			handlePodYaml(data)
		}
	},
}

func init() {
	apply.Flags().StringVarP(&filePath, "file", "f", "pod.yaml", "Path to the yaml file")
	root.AddCommand(apply)
}
