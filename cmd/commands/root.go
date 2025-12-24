/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"time"

	"github.com/OvsienkoValeriya/GophKeeper/internal/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr  string
	tlsEnabled  bool
	tlsCertPath string

	grpcConn       *grpc.ClientConn
	authClient     *client.AuthClient
	tokenStore     *client.FileTokenStore
	masterKeyStore *client.MasterKeyStore
	resourceClient *client.ResourceClient
	clientConfig   *client.ClientConfig
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gophkeeper",
	Short: "GophKeeper - менеджер паролей и секретов",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		clientConfig = client.LoadConfig()

		if cmd.Flags().Changed("server") {
			clientConfig.ServerAddress = serverAddr
		} else {
			serverAddr = clientConfig.ServerAddress
		}
		if cmd.Flags().Changed("tls") {
			clientConfig.TLSEnabled = tlsEnabled
		}
		if cmd.Flags().Changed("tls-cert") {
			clientConfig.TLSCertPath = tlsCertPath
		}

		var opts []grpc.DialOption
		if clientConfig.TLSEnabled {
			creds, err := credentials.NewClientTLSFromFile(clientConfig.TLSCertPath, "")
			if err != nil {
				log.Fatalf("Failed to load TLS credentials: %v", err)
			}
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}

		grpcConn, err = grpc.NewClient(serverAddr, opts...)
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}

		authClient = client.NewAuthClient(grpcConn)
		tokenStore, err = client.NewFileTokenStore()
		if err != nil {
			log.Fatalf("Failed to create token store: %v", err)
		}
		resourceClient = client.NewResourceClient(grpcConn, tokenStore)

		masterKeyStore = client.NewMasterKeyStore()

		cmdName := cmd.Name()
		noAuthCommands := map[string]bool{
			"register": true,
			"login":    true,
			"help":     true,
		}

		softAuthCommands := map[string]bool{
			"logout":  true,
			"refresh": true,
		}

		if noAuthCommands[cmdName] {
			return
		}

		expired, err := tokenStore.IsAccessTokenExpired()
		if err != nil {
			if softAuthCommands[cmdName] {
				return
			}
			log.Fatalf("Not authenticated. Please login or register first.")
		}
		if expired {
			_, refreshToken, err := tokenStore.LoadTokens()
			if err != nil {
				if softAuthCommands[cmdName] {
					return
				}
				log.Fatalf("Failed to load tokens: %v", err)
			}
			respRefresh, err := authClient.RefreshToken(refreshToken)
			if err != nil {
				if softAuthCommands[cmdName] {
					return
				}
				log.Fatalf("Session expired. Please login again.")
			}
			userID, _ := tokenStore.GetUserID()
			tokenStore.SaveTokensWithUserID(userID, respRefresh.GetAccessToken(), respRefresh.GetRefreshToken(), time.Now().Add(time.Hour*1))
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if grpcConn != nil {
			grpcConn.Close()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cfg := client.LoadConfig()

	rootCmd.PersistentFlags().StringVarP(&serverAddr, "server", "s", cfg.ServerAddress, "gRPC server address")
	rootCmd.PersistentFlags().BoolVar(&tlsEnabled, "tls", cfg.TLSEnabled, "Enable TLS connection")
	rootCmd.PersistentFlags().StringVar(&tlsCertPath, "tls-cert", cfg.TLSCertPath, "Path to TLS certificate file")
}
