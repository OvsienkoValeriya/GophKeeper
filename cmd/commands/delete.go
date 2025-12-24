/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a resource by name",
	Long:  `gophkeeper delete <name>`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		resource, err := resourceClient.GetResourceByName(name)
		if err != nil {
			fmt.Printf("✗ Secret '%s' not found: %v\n", name, err)
			return
		}

		if err := resourceClient.DeleteResource(resource.GetId()); err != nil {
			fmt.Printf("✗ Failed to delete secret: %v\n", err)
			return
		}

		fmt.Printf("✓ Secret '%s' deleted successfully\n", name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

}
