/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"flag"
	"os"

	"github.com/aririfani/wallet-service/internal/app/bootstrap"
	"github.com/aririfani/wallet-service/internal/app/handler"
	"github.com/aririfani/wallet-service/internal/app/wallet"
	"github.com/spf13/cobra"
)

var (
	// deposit = flag.Bool("deposit", false, "send Deposit")
	broker = flag.String("broker", "localhost:9092", "boostrap kafka broker")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wallet-service",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		StartService()
	},
}

var runBroker = &cobra.Command{
	Use:   "broker:up",
	Short: "broker up",
	Long:  `Broker up command`,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.Run()
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wallet-service.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(runBroker)
}

func StartService() {
	flag.Parse()
	handler.Run([]string{*broker}, wallet.WalletTopic)
}
