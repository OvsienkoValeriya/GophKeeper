/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Add encrypted data to the storage",
	Long: `Add encrypted data to the storage.

Examples:
  # Store a text value
  gophkeeper set -n "secret" -v "my-password" -t text

  # Store a file (for large data)
  gophkeeper set -n "bigfile" -f /path/to/file -t binary`,
	Run: func(cmd *cobra.Command, args []string) {

		cryptoService, err := masterKeyStore.GetCryptoService()
		if err != nil {
			fmt.Println("✗ Secrets are locked. Run 'gophkeeper unlock' first.")
			return
		}

		name, _ := cmd.Flags().GetString("name")
		value, _ := cmd.Flags().GetString("value")
		filePath, _ := cmd.Flags().GetString("file")
		secretType, _ := cmd.Flags().GetString("type")

		validTypes := map[string]bool{
			"credentials": true,
			"text":        true,
			"binary":      true,
			"card":        true,
		}
		if !validTypes[secretType] {
			fmt.Println("Invalid type. Use: credentials, text, binary, or card")
			return
		}

		if value == "" && filePath == "" {
			fmt.Println("✗ Either --value or --file must be provided")
			return
		}
		if value != "" && filePath != "" {
			fmt.Println("✗ Cannot use both --value and --file")
			return
		}

		var dataToEncrypt []byte

		if filePath != "" {
			dataToEncrypt, err = os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("✗ Failed to read file: %v\n", err)
				return
			}
			fmt.Printf("Read %d bytes from file\n", len(dataToEncrypt))
		} else {
			dataToEncrypt = []byte(value)
		}

		encryptedData, err := cryptoService.EncryptData(dataToEncrypt)
		if err != nil {
			fmt.Printf("✗ Encryption failed: %v\n", err)
			return
		}

		resourceID, err := resourceClient.CreateResource(name, secretType, encryptedData)
		if err != nil {
			fmt.Printf("✗ Failed to save secret: %v\n", err)
			return
		}

		fmt.Printf("✓ Secret '%s' saved (ID: %d, size: %d bytes)\n", name, resourceID, len(dataToEncrypt))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().StringP("name", "n", "", "Name of the secret")
	setCmd.Flags().StringP("value", "v", "", "Value to store (for small data)")
	setCmd.Flags().StringP("file", "f", "", "Path to file (for large data)")
	setCmd.Flags().StringP("type", "t", "", "Type: credentials | text | binary | card")
	setCmd.MarkFlagRequired("name")
	setCmd.MarkFlagRequired("type")
}
