/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the system",
	Long:  `gophkeeper login -u <username> -p <password>`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		resp, err := authClient.Login(username, password)
		if err != nil {
			fmt.Printf("Login failed: %v\n", err)
			return
		}

		userID, _ := strconv.ParseUint(resp.GetUserId(), 10, 64)

		expiresAt := time.Now().Add(time.Hour * 1)
		err = tokenStore.SaveTokensWithUserID(uint(userID), resp.GetAccessToken(), resp.GetRefreshToken(), expiresAt)
		if err != nil {
			fmt.Printf("Failed to save tokens: %v\n", err)
			return
		}
		fmt.Printf("Login successful! User ID: %s\n", resp.GetUserId())

		respMK, err := authClient.GetMasterKeyData(resp.GetAccessToken())
		if err != nil {
			fmt.Printf("Warning: failed to get master key data: %v\n", err)
			return
		}

		if !respMK.GetHasMasterKey() {
			fmt.Println("\n⚠ Master key not initialized. Run 'gophkeeper init' to set it up.")
			return
		}

		fmt.Println("\nEnter master key to unlock your secrets:")
		masterPassword, err := promptPassword("Master key: ")
		if err != nil {
			fmt.Printf("Error reading password: %v\n", err)
			return
		}

		if err := masterKeyStore.Unlock(masterPassword, respMK.GetSalt(), respMK.GetVerifier()); err != nil {
			fmt.Printf("✗ Failed to unlock: %v\n", err)
			fmt.Println("Your secrets remain locked. Use 'gophkeeper unlock' to try again.")
			return
		}

		tokenStore.SetHasMasterKey(true)
		fmt.Println("✓ Secrets unlocked!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("password", "p", "", "Password")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}
