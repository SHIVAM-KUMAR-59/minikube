package cmd

import "github.com/spf13/cobra"

var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources from the minik cluster",
	Long:  `Delete resources from the minik cluster such as pods, services, etc.`,
}

func init() {
	root.AddCommand(Delete)
}
