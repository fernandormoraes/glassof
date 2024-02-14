/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/fernandormoraes/glassof/entities"
	"github.com/fernandormoraes/glassof/pgoutput"
	"github.com/jackc/pgx"
	"github.com/spf13/cobra"
)

// replicateCmd represents the replicate command
var replicateCmd = &cobra.Command{
	Use:   "replicate",
	Short: "Replicate changed data",
	Long: `Command takes data from logical slot and insert as a collection in target database.
	
	Take a ./glassof replicate -- to start listening data.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "" {
			fmt.Printf("Target database arg is missing.\n")
			return
		}

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

	conn, err := pgx.ReplicationConnect(configPgx)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	err = conn.Ping(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Printf("Connection succeded\n")

	pluginArguments := []string{
		"proto_version '2'",
		"publication_names 'pub_glassof'",
		"messages 'false'",
		"streaming 'true'",
	}

	err = conn.StartReplication("glassof", 0, -1, pluginArguments...)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	set := pgoutput.NewRelationSet(nil)

	dump := func(relation uint32, row []pgoutput.Tuple) error {
		values, err := set.Values(relation, row)
		if err != nil {
			return fmt.Errorf("error parsing values: %s", err)
		}
		for name, value := range values {
			val := value.Get()
			log.Printf("%s (%T): %#v", name, val, val)
		}
		return nil
	}

	handler := func(m pgoutput.Message, walPos uint64) error {
		switch v := m.(type) {
		case pgoutput.Relation:
			log.Printf("RELATION")
			set.Add(v)
		case pgoutput.Insert:
			log.Printf("INSERT")
			return dump(v.RelationID, v.Row)
		case pgoutput.Update:
			log.Printf("UPDATE")
			return dump(v.RelationID, v.Row)
		case pgoutput.Delete:
			log.Printf("DELETE")
			return dump(v.RelationID, v.Row)
		}
		return nil
	}

	log.Print("Waiting for replication message\n")

	for {

		msg, err := conn.WaitForReplicationMessage(ctx)

		if err != nil {
			panic(err)
		}

		log.Print("Reading data\n")

		if msg.WalMessage != nil {
			logmsg, err := pgoutput.Parse(msg.WalMessage.WalData)

			if err != nil {
				panic(err)
			}

			log.Print("Handling data\n")

			handler(logmsg, msg.WalMessage.WalStart)

			if msg.ServerHeartbeat != nil {
				log.Printf("Got heartbeat: %s", msg.ServerHeartbeat)
			}
		}
	}
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
