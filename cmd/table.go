/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/cockroachdb/pebble"
	"github.com/jackc/pgx"
	"github.com/spf13/cobra"
)

// tableCmd represents the table command
var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Command to list and add tables",
	Long: `Adding a table to replication service only means that service will listen to changes in these tables.

	Take a ./glassof table add table_name     -- for adding tables
	Take a ./glassof table list               -- for listing tables added
	Take a ./glassof table rm table_name      -- for removing tables
	
	It is necessary to take a ./glassof primary   who will be responsible for identifying which row changed
	It is necessary to take a ./glassof query     who will be responsible for get data in a certain way and replicate`,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "add" {
			addTable(args)
		} else if args[0] == "list" {
			listTables()
		} else if args[0] == "rm" {
			removeTable(args)
		}

	},
}

func removeTable(args []string) {
	fmt.Println("Tables before delete: ")
	tablesList := listTables()

	if strings.Contains(tablesList, args[1]) {
		after := strings.Replace(tablesList, args[1]+";", "", -1)

		db, err := pebble.Open("db", &pebble.Options{})

		if err != nil {
			log.Println("Error opening database")
			log.Fatal(err)
		}

		if after != "" && after[len(after)-1] != ';' {
			after = after + ";"
		}

		if errorInsert := db.Set([]byte("tables"), []byte(after), pebble.Sync); errorInsert != nil {
			log.Println("Error inserting in tables.")
			log.Fatal(errorInsert)
		}

		if errorDelete := db.Delete([]byte(fmt.Sprintf("table.%s", args[1])), pebble.Sync); errorDelete != nil {
			log.Println("Error deleting tables.")
			log.Fatal(errorDelete)
		}

		fmt.Println("Tables after delete: ")

		fmt.Println(after)

	} else {
		fmt.Println("Table was not found in Glassof.")
	}
}

func listTables() string {
	db, err := pebble.Open("db", &pebble.Options{})

	if err != nil {
		log.Println("Error opening database")
		log.Fatal(err)
	}

	defer db.Close()

	res, closer, err := db.Get([]byte("tables"))

	if err != nil {
		fmt.Println("Error reading tables")
		log.Fatal(err)
	}

	defer closer.Close()

	fmt.Println("List of tables: ")
	fmt.Println(string(res[:]))

	return string(res[:])
}

func addTable(args []string) {
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

	if err != pebble.ErrNotFound {
		log.Printf("Table is already in Glassof.")

		if err := closer.Close(); err != nil {
			log.Println("Error closing closer")
			log.Fatal(err)
		}

		return
	}

	if errorInsert := db.Set([]byte(formattedKey), []byte(""), pebble.Sync); errorInsert != nil {
		log.Println("Error inserting table in database")
		log.Fatal(errorInsert)
	}

	res, closer, err := db.Get([]byte("tables"))

	if err != nil && err != pebble.ErrNotFound {
		log.Println("Error reading database")
		log.Fatal(err)
	}

	resTables := string(res[:])

	if err == nil {
		if err := closer.Close(); err != nil {
			log.Println("Error closing closer")
			log.Fatal(err)
		}
	}

	if errorInsert := db.Set([]byte("tables"), []byte(args[1]+";"+resTables), pebble.Sync); errorInsert != nil {
		log.Println("Error inserting in tables.")
		log.Fatal(errorInsert)
	}

	postgresUri, closer, err := db.Get([]byte("connectionUri"))

	if err == nil {
		if err := closer.Close(); err != nil {
			log.Println("Error closing closer")
			log.Fatal(err)
		}
	}

	configPgx, err := pgx.ParseConnectionString(string(postgresUri[:]))

	if err != nil {
		fmt.Printf("Error parsing connection URI\n")
		panic(err)
	}

	conn, err := pgx.Connect(configPgx)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	_, err = conn.Exec("ALTER TABLE " + args[1] + " REPLICA IDENTITY FULL;")

	if err != nil {
		panic(err)
	}

	fmt.Printf("Added a glass of table %s.\n", args[1])
}

func init() {
	rootCmd.AddCommand(tableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
