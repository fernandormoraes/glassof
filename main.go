/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"

	"github.com/cockroachdb/pebble"
	"github.com/fernandormoraes/glassof/cmd"
)

func main() {
	db, err := pebble.Open("db", &pebble.Options{})

	if err != nil {
		log.Fatal(err)
	}

	db.Close()

	cmd.Execute()
}
