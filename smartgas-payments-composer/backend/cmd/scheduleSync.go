/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"smartgas-payment/internal/injectors"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
)

// scheduleSyncCmd represents the scheduleSync command
var scheduleSyncCmd = &cobra.Command{
	Use:   "scheduleSync",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		syncTask, err := injectors.InitializeSynchronizationTask()
		if err != nil {
			panic(err)
		}
		loc, err := time.LoadLocation("America/Mazatlan")
		if err != nil {
			panic(err)
		}

		log.Println("Init schedule synchronization for gas stations")

		s := gocron.NewScheduler(loc)

		//s.Every("10s").Do(func() {
		//fmt.Println("Called")
		//})

		s.Every(1).Day().At("07:00;15:00").Do(func() {
			syncTask.SyncGasStations()
			syncTask.SyncGasPumps()
		})

		log.Println("Init schedule for customer levels")
		s.Every(1).Month(1).At("00:01").Do(func() {
			syncTask.GenerateElegibilityCustomers()
		})

		s.StartBlocking()
	},
}

func init() {
	rootCmd.AddCommand(scheduleSyncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scheduleSyncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scheduleSyncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
