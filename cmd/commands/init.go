/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the master key",
	Long: `Important: Remember your master key - it cannot be recovered!
	Example:
	gophkeeper init`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := tokenStore.GetUserID()
		if err != nil {
			fmt.Println("Error: You must be logged in to initialize master key")
			fmt.Println("Please run 'login' or 'register' first")
			return
		}

		hasMasterKey, _ := tokenStore.HasMasterKey()
		if hasMasterKey {
			fmt.Println("Master key is already initialized.")
			fmt.Println("If you want to reset it, please contact support.")
			return
		}

		masterPassword, err := promptPassword("Enter master key: ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		confirmPassword, err := promptPassword("Confirm master key: ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if masterPassword != confirmPassword {
			fmt.Println("Error: Master keys do not match")
			return
		}

		if len(masterPassword) < 8 {
			fmt.Println("Error: Master key must be at least 8 characters")
			return
		}

		salt, verifier, err := masterKeyStore.SetupAndUnlock(masterPassword)
		if err != nil {
			fmt.Printf("Error setting up master key: %v\n", err)
			return
		}

		accessToken, _, err := tokenStore.LoadTokens()
		if err != nil {
			fmt.Printf("Error loading tokens: %v\n", err)
			return
		}

		if _, err := authClient.SetMasterKey(accessToken, salt, verifier); err != nil {
			fmt.Printf("Error saving master key to server: %v\n", err)
			masterKeyStore.Lock()
			return
		}

		if err := tokenStore.SetHasMasterKey(true); err != nil {
			fmt.Printf("Error saving master key status: %v\n", err)
			return
		}

		fmt.Println("✓ Master key initialized successfully!")
		fmt.Println("⚠ IMPORTANT: Remember your master key - it cannot be recovered!")
	},
}

// promptPassword prompts the user for a password without displaying the input
func promptPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(line), nil
	}

	fmt.Println()
	return string(password), nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
