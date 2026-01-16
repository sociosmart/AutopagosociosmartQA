/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"smartgas-payment/internal/injectors"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/utils"
	"strings"

	"github.com/spf13/cobra"
)

// createuserCmd represents the createuser command
var createuserCmd = &cobra.Command{
	Use:   "createuser",
	Short: "Create user interactive mode",
	Long:  `User will be promped for data in order to provide faster response`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			email       string
			password    string
			isAdmin     string
			isAdminBool bool
		)

		fmt.Print("Email: ")
		fmt.Scan(&email)
		fmt.Print("Password: ")
		fmt.Print("\033[8m")
		fmt.Scan(&password)
		fmt.Println("\033[28m")
		fmt.Print("Is admin? [Y/n]: ")
		fmt.Scan(&isAdmin)

		if strings.Contains(strings.ToLower(isAdmin), "y") {
			isAdminBool = true
		}

		repo, err := injectors.InitializeUserRepository()

		//defer database.CloseConnection(repo.DB)

		//if err != nil {
		//panic(err)
		//}

		user := &models.User{
			Email:    email,
			Password: password,
			IsAdmin:  utils.BoolAddr(isAdminBool),
		}

		err = repo.CreateUser(user)

		if err != nil {
			panic(err)
		}

		fmt.Println("User created successfuly")

	},
}

func init() {
	rootCmd.AddCommand(createuserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createuserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createuserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
