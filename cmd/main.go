package main

import (
	"ecoctl/pkg/commands"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	cobra.OnInitialize()
}

func main() {
	if err := commands.Root.Execute(); err != nil {
		log.Fatal(err)
	}
}
