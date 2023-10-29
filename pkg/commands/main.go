package commands

import (
	"context"
	"ecoctl/pkg/client"
	"ecoctl/pkg/version"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var url = ""
var token = ""
var eco *client.Client
var output string = "json"

func init() {
	cobra.OnInitialize(initConfiguration)

	Root.PersistentFlags().StringVarP(&url, "url", "u", "", "eco-server`s url")
	Root.PersistentFlags().StringVarP(&token, "token", "t", "", "apikey to interact with eco")
	Root.PersistentFlags().StringVarP(&output, "output", "o", "", "output format(yaml,json)")

	if err := initClusterListFlags(); err != nil {
		log.Fatal(err)
	}
	if err := initClusterCreateFlags(); err != nil {
		log.Fatal(err)
	}
	if err := initClusterConfigFlags(); err != nil {
		log.Fatal(err)
	}
	if err := initPoolCreateFlags(); err != nil {
		log.Fatal(err)
	}
	if err := initPoolDeleteFlags(); err != nil {
		log.Fatal(err)
	}
	if err := initPoolUpdateFlags(); err != nil {
		log.Fatal(err)
	}

	clusterCmd.AddCommand(clusterListCmd)
	clusterCmd.AddCommand(clusterCreateCmd)
	clusterCmd.AddCommand(clusterGetCmd)
	clusterCmd.AddCommand(clusterDeleteCmd)
	clusterCmd.AddCommand(clusterWatchCmd)
	clusterCmd.AddCommand(clusterConfigCmd)
	Root.AddCommand(clusterCmd)

	poolCmd.AddCommand(poolCreateCmd)
	poolCmd.AddCommand(poolUpdateCmd)
	poolCmd.AddCommand(poolDeleteCmd)
	Root.AddCommand(poolCmd)
}

func initConfiguration() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("err during load env")
		log.Fatal(err)
	}
	if url == "" {
		url = os.Getenv("ECO_URL")
		if url == "" {
			log.Fatal("eco-url must be set (use '-u' flag or ECO_URL env)")
		}
	}
	if token == "" {
		token = os.Getenv("ECO_TOKEN")
		if token == "" {
			log.Fatal("eco-token must be set (use '-t' flag or ECO_TOKEN env)")
		}
	}
	if pathToClusterCreationTemplate != "" {
		err = pathToTemplateIsValid(pathToClusterCreationTemplate)
		if err != nil {
			log.Fatal(err)
		}
		path, filename, ext := splitPath(pathToClusterCreationTemplate)
		if !exists(path) {
			log.Fatal("path doesnt exist")
		}
		viper.AddConfigPath(path)
		viper.SetConfigName(filename)
		viper.SetConfigType(ext[1:])
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
		err = viper.Unmarshal(&clusterCreateOpt)
		if err != nil {
			log.Fatal(err)
		}
	}
	if pathToPoolCreationTemplate != "" {
		err = pathToTemplateIsValid(pathToPoolCreationTemplate)
		if err != nil {
			log.Fatal(err)
		}
		path, filename, ext := splitPath(pathToPoolCreationTemplate)
		if !exists(path) {
			log.Fatal("path doesnt exist")
		}
		viper.AddConfigPath(path)
		viper.SetConfigName(filename)
		viper.SetConfigType(ext[1:])
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
		err = viper.Unmarshal(&poolCreateOpt)
		if err != nil {
			log.Fatal(err)
		}
	}
	eco, err = client.NewWithRetries(nil, client.SetAPIKey(token), client.SetBaseURL(url))
	if err != nil {
		log.Printf("err during creating client")
		log.Fatal(err)
	}
	v, _, err := eco.Service.Version(context.Background())
	if err != nil {
		log.Printf("err during getting version")
		log.Fatal(err)
	}
	fmt.Printf("ecoctl version: %s\neco-server url: %s\neco-server version: %v\n", version.Version, url, v)
	//version request
}

var Root = &cobra.Command{
	Use:   "ecoctl",
	Short: "ecoctl is tool for interact with eco-server",
}
