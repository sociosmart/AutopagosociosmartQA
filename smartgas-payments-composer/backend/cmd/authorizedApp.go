/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"smartgas-payment/internal/injectors"
	"smartgas-payment/internal/models"

	"github.com/spf13/cobra"
)

// authorizedAppCmd represents the authorizedApp command
var authorizedAppCmd = &cobra.Command{
	Use:   "authorizedApp",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		appName := args[0]

		repository, err := injectors.InitializeSecurityRepository()

		if err != nil {
			panic(err)
		}

		authorizedApp := &models.AuthorizedApplication{
			ApplicationName: appName,
		}

		err = repository.Create(authorizedApp)

		if err != nil {
			panic(err)
		}

		fmt.Println("AppKey", authorizedApp.AppKey)
		fmt.Println("ApiKey", authorizedApp.ApiKey)
	},
}

func init() {
	rootCmd.AddCommand(authorizedAppCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//authorizedAppCmd.PersistentFlags().String("name", "", "The authorized application name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authorizedAppCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
