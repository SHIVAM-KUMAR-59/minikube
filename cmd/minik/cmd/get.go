package cmd

import "github.com/spf13/cobra"

var Get = &cobra.Command{
	Use:   "get",
	Short: "Get resources from the minik cluster",
	Long:  `Get resources from the minik cluster such as pods, services, etc.`,
}

func init() {
	root.AddCommand(Get)
}