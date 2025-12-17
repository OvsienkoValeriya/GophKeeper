/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// unlockCmd represents the unlock command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "unlock secrets with master key",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if masterKeyStore.IsUnlocked() {
			fmt.Println("Already unlocked.")
			return
		}

		accessToken, _, err := tokenStore.LoadTokens()
		if err != nil {
			fmt.Println("Not logged in. Please login first.")
			return
		}

		salt, verifier, hasMasterKey, err := authClient.GetMasterKeyData(accessToken)
		if err != nil {
			fmt.Printf("Failed to get master key data: %v\n", err)
			return
		}

		if !hasMasterKey {
			fmt.Println("Master key not initialized. Run 'gophkeeper init' first.")
			return
		}

		masterPassword, err := promptPassword("Enter master key: ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := masterKeyStore.Unlock(masterPassword, salt, verifier); err != nil {
			fmt.Printf("✗ Invalid master key: %v\n", err)
			return
		}

		fmt.Println("✓ Secrets unlocked!")
	},
}

func init() {
	rootCmd.AddCommand(unlockCmd)
}
