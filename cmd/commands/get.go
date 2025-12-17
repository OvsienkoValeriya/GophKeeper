/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get encrypted data from the storage",
	Long:  `gophkeeperk get <name>`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cryptoService, err := masterKeyStore.GetCryptoService()
		if err != nil {
			fmt.Println("✗ Secrets are locked. Run 'gophkeeper unlock' first.")
			return
		}

		response, err := resourceClient.GetResourceByName(name)
		if err != nil {
			fmt.Printf("✗ Failed to get secret: %v\n", err)
			return
		}

		decryptedData, err := cryptoService.DecryptData(response.GetData())
		if err != nil {
			fmt.Printf("✗ Decryption failed: %v\n", err)
			return
		}

		fmt.Printf("Name: %s\n", response.GetName())
		fmt.Printf("Type: %s\n", response.GetType())
		fmt.Printf("Value: %s\n", string(decryptedData))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
