/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "refresh the access token",
	Long:  `gophkeeper refresh`,
	Run: func(cmd *cobra.Command, args []string) {
		storedAccessToken, _, err := tokenStore.LoadTokens()
		if err != nil {
			fmt.Printf("Failed to load tokens: %v\n", err)
			return
		}

		respRefresh, err := authClient.RefreshToken(storedAccessToken)
		if err != nil {
			fmt.Printf("Failed to refresh token: %v\n", err)
			return
		}

		expiresAt := time.Now().Add(time.Hour * 1)
		userID, err := tokenStore.GetUserID()
		if err != nil {
			fmt.Printf("Failed to get user ID: %v\n", err)
			return
		}
		tokenStore.SaveTokensWithUserID(userID, respRefresh.GetAccessToken(), respRefresh.GetRefreshToken(), expiresAt)

		fmt.Println("✓ Tokens refreshed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
