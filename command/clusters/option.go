package clusters

import (
	"ecoctl/command"
	"errors"
	"github.com/ipaqsa/ecogo"
	"github.com/spf13/viper"
	"log"
)

var pathToTemplate string

var createOpt ecogo.ClusterOpt

func initFlags() error {
	createCmd.Flags().StringVarP(&pathToTemplate, "template", "t", "",
		"template with cluster configuration")
	if err := createCmd.MarkFlagRequired("template"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}
	return nil
}

func parseConfig() {
	if pathToTemplate == "" {
		log.Fatal("path must not be empty")
	}
	if err := command.PathIsValid(pathToTemplate); err != nil {
		log.Fatal(err)
	}
	path, filename, ext := command.SplitPath(pathToTemplate)
	if !command.Exists(path) {
		log.Fatal("path doesnt exist")
	}
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType(ext[1:])
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Unmarshal(&createOpt); err != nil {
		log.Fatal(err)
	}
}
