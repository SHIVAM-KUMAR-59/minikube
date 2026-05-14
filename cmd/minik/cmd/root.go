package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
  Use:   "minik",
  Short: "Minik is a mini kubernetes cluster for container orchestration",
  Long: `Minik is a mini kubernetes cluster for container orchestration. It is a tool that makes it easy to run Kubernetes locally. Minikube runs a single-node Kubernetes cluster inside a virtual machine on your laptop for users looking to try out Kubernetes or develop with it day-to-day.`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

func Execute() {
  if err := root.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
	slog.Error("Command execution failed", "error", err)
    os.Exit(1)
  }
}