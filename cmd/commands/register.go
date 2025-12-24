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

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Long:  `gophkeeper register -u <username> -p <password>`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		respRegister, err := authClient.Register(username, password)
		if err != nil {
			fmt.Printf("Registration failed: %v\n", err)
			return
		}

		userID, _ := strconv.ParseUint(respRegister.GetUserId(), 10, 64)

		expiresAt := time.Now().Add(time.Hour)
		if err := tokenStore.SaveTokensWithUserID(uint(userID), respRegister.GetAccessToken(), respRegister.GetRefreshToken(), expiresAt); err != nil {
			fmt.Printf("Failed to save tokens: %v\n", err)
			return
		}
		fmt.Printf("Registration successful! User ID: %s\n", respRegister.GetUserId())
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringP("username", "u", "", "Username")
	registerCmd.Flags().StringP("password", "p", "", "Password")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")

}
