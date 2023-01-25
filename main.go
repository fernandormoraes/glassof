/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"encoding/json"

	"github.com/fernandormoraes/glassof/cmd"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle //2

func main() {
	bundle = i18n.NewBundle(language.English)            //4
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal) //5
	bundle.LoadMessageFile("resources/en.json")          //6
	bundle.LoadMessageFile("resources/br.json")          //7

	cmd.Execute()
}
