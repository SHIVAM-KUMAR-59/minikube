package cmd

import "github.com/spf13/cobra"

var Cluster = &cobra.Command{
	Use:   "cluster",
	Short: "Start or Stop workers",
	Long:  "Start or stop the workers running in the background",
}

func init() {
	root.AddCommand(Cluster)
}
