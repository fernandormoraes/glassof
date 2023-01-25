/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Postgres database to prepare for replicate data",
	Long: `Glassof uses WAL replication configuration of Postgres to get updated data from database.

	--uri Postgres database URI - Default value: postgresql://postgres:postgres@localhost:5432/postgres
	
	Take a ./glassof init --uri postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]`,
	Run: func(cmd *cobra.Command, args []string) {
		postgresUri := "postgresql://postgres:postgres@localhost:5432/postgres"

		if len(args) >= 1 && args[0] != "" {
			postgresUri = args[0]
		}

		fmt.Printf("Your URI is: %s\n", postgresUri)

		configPgx, err := pgx.ParseConnectionString(postgresUri)

		if err != nil {
			fmt.Printf("Error parsing connection URI\n")
			panic(err)
		}

		conn, err := pgx.Connect(configPgx)

		if err != nil {
			panic(err)
		}

		defer conn.Close()

		err = conn.Ping(context.Background())

		if err != nil {
			panic(err)
		}

		fmt.Printf("Connection succeded\n")

		var dest interface{}

		errSlot := conn.QueryRow("SELECT * FROM pg_create_logical_replication_slot('glassof', 'wal2json', false);").Scan(dest)

		if errSlot != nil {
			fmt.Printf("Error creating logical replication\n")
			fmt.Printf("Maybe you forget to turn wal_level = logical in postgres.conf\n")
			fmt.Printf("Or maybe you forget to install wal2json plugin in postgresql instance?\n")

			if strings.Contains(errSlot.Error(), "already exists") {
				fmt.Printf("glassof replication slot already exists, proceeding...\n")
			} else {
				panic(errSlot)
			}

		}

		fmt.Printf("Created logical replication in database with success\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
