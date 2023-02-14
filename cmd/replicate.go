/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/cockroachdb/pebble"
	"github.com/fernandormoraes/glassof/entities"
	"github.com/jackc/pgx"
	"github.com/spf13/cobra"
)

// replicateCmd represents the replicate command
var replicateCmd = &cobra.Command{
	Use:   "replicate",
	Short: "Replicate changed data",
	Long: `Command takes data from logical slot and insert as a collection in target database.
	
	Take a ./glassof replicate [mongo_db_uri] to replicate data to a mongodb target.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "" {
			fmt.Printf("Target database arg is missing.\n")
			return
		}

		//tableList := listTables()

		//queriesList := listQueries(tableList)

		dataChanged := getDataPostgre()

		fmt.Print(dataChanged)

	},
}

func getDataPostgre() []entities.Slot {
	postgresUri := "postgresql://postgres:postgres@localhost:5432/postgres"

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

	rows, errSlot := conn.Query("SELECT data FROM pg_logical_slot_get_changes('glassof', NULL, NULL);")

	if errSlot != nil {
		log.Print(err)
		log.Printf("Error getting slot changes.")
	}

	var listRows []entities.Data

	for rows.Next() {
		var data entities.Data

		rows.Scan(&data.Data)

		listRows = append(listRows, data)
	}

	slot := rowsToSlotEntity(listRows)

	if err != nil {
		log.Print(err)
	}

	defer rows.Close()

	return slot
}

func rowsToSlotEntity(rows []entities.Data) []entities.Slot {
	var listSlots []entities.Slot

	for _, row := range rows {
		var slot entities.Slot

		json.Unmarshal([]byte(row.Data), &slot)

		listSlots = append(listSlots, slot)
	}

	return listSlots
}

func listQueries(tableList string) []string {

	tables := strings.Split(tableList, ";")

	db, err := pebble.Open("db", &pebble.Options{})

	if err != nil {
		log.Println("Error opening database")
		log.Fatal(err)
	}

	defer db.Close()

	var queries []string

	for counter, table := range tables {
		data, closer, err := db.Get([]byte(fmt.Sprintf("table.%s", table)))

		if err != nil {
			log.Println("Fail getting queries")
		}

		queries[counter] = string(data[:])

		if err := closer.Close(); err != nil {
			log.Println("Fail closing queries")
		}
	}

	return queries
}

func init() {
	rootCmd.AddCommand(replicateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// replicateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// replicateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
