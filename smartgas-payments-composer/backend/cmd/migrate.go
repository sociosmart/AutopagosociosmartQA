/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"smartgas-payment/internal/database"
	"smartgas-payment/internal/injectors"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate models",
	Long:  `MIgrate models in database`,
	Run: func(cmd *cobra.Command, args []string) {

		db, err := injectors.InitializeDB()

		if err != nil {
			panic(err)
		}

		defer database.CloseConnection(db)

		database.RunMigrations(db)

	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
