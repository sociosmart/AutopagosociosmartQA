/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"smartgas-payment/cmd"
	"smartgas-payment/config"
	"smartgas-payment/internal/utils"
)

func main() {
	// Making sure the timezone is properly setted
	initialize()
	cmd.Execute()
}

func initialize() {

	cfg, err := config.NewConfig()

	if err != nil {
		panic(err)
	}
	utils.InitTimezone(cfg.Tz)

}
