/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/cockroachdb/pebble"
	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Add a query to a table for replicate date",
	Long: `Query command takes a PostgreSQL query and will perform this query to replicate data to MongoDB.

	Take a ./glassof query add [table_name] [query_string]   -- for adding a query to a table
	
	Example: ./glassof query add products "SELECT ID, DESCRIPTION FROM PRODUCTS"  -- Added query "SELECT ID, DESCRIPTION FROM PRODUCTS" to products table`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Adding table to replication service")

		formattedKey := fmt.Sprintf("table.%s", args[1])

		db, err := pebble.Open("db", &pebble.Options{})

		if err != nil {
			log.Println("Error opening database")
			log.Fatal(err)
		}

		defer db.Close()

		_, closer, err := db.Get([]byte(formattedKey))

		if err != nil && err != pebble.ErrNotFound {
			log.Println("Error reading database")
			log.Fatal(err)
		}

		if err == pebble.ErrNotFound {
			log.Printf("Table isn't added to glassof.")

			return
		}

		if err := closer.Close(); err != nil {
			log.Println("Error closing closer")
			log.Fatal(err)
		}

		if errorInsert := db.Set([]byte(formattedKey), []byte(args[2]), pebble.Sync); errorInsert != nil {
			log.Println("Error inserting query in table")
			log.Fatal(errorInsert)
		}

		fmt.Printf("Added a glass of query.\n")
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
