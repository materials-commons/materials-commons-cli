/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create the database that the mcc command uses.",
	Long:  `Create the database that the mcc command uses.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(config.GetProjectDBPath())
		if os.IsNotExist(err) {
			f, err := os.Create(config.GetProjectDBPath())
			if err != nil {
				log.Fatalf("Unable to create database file: %s", err)
			}
			_ = f.Close()
		}

		db := mcdb.MustConnectToDB()

		if err := mcdb.RunMigrations(db); err != nil {
			log.Fatalf("Unable to run migrations: %s", err)
		}

		fmt.Println("Database updated!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
