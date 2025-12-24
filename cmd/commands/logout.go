/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout from the system",
	Long:  `gophkeeper logout`,
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			accessToken, _, err := tokenStore.LoadTokens()
			if err != nil {
				fmt.Println("You are not logged in.")
				return
			}
			token = accessToken
		}

		_, err := authClient.Logout(token)
		if err != nil {
			fmt.Printf("Failed to logout: %v\n", err)
			return
		}
		tokenStore.ClearTokens()
		fmt.Printf("Logged out successfully!\n")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
	logoutCmd.Flags().StringP("token", "t", "", "Access token (if not specified, loads from token store)")
}
