/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"smartgas-payment/internal/injectors"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Long: `
  
  `,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		syncTask, err := injectors.InitializeSynchronizationTask()
		if err != nil {
			panic(err)
		}

		if len(args) == 0 {
			err := syncTask.SyncGasStations()
			if err != nil {
				log.Fatalln(err)
			}
			err = syncTask.SyncGasPumps()

			if err != nil {
				log.Fatalln(err)
			}
			os.Exit(0)
		} else if args[0] == "gas-stations" {
			syncTask.SyncGasStations()
		} else if args[0] == "gas-pumps" {
			syncTask.SyncGasPumps()
		} else if args[0] == "customer-levels" {
			syncTask.GenerateElegibilityCustomers()
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("gas-stations", "", "Synchronize gas-stations only")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
