package main

import (
	"ecoctl/command"
	"ecoctl/command/clusters"
	"ecoctl/command/pools"
	"ecoctl/command/users"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	cobra.OnInitialize(command.InitEnvs)
	rootCmd.AddCommand(clusters.Cmd)
	rootCmd.AddCommand(pools.Cmd)
	rootCmd.AddCommand(users.UserCmd)
	rootCmd.AddCommand(users.AdminCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ecoctl",
	Short: "ecoctl is tool to interact with eco-server",
}
