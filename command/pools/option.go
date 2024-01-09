package pools

import (
	"ecoctl/command"
	"errors"
	"github.com/ipaqsa/ecogo"
	"github.com/spf13/viper"
	"log"
)

var pathToTemplate string

var opt ecogo.PoolOpt

func initFlags() error {
	Cmd.PersistentFlags().UintVarP(&command.GlobalOption.ClusterID, "cluster", "c", 0, "clusterID")
	if err := Cmd.MarkPersistentFlagRequired("cluster"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}

	createCmd.Flags().StringVarP(&pathToTemplate, "template", "t", "",
		"template with pool configuration")
	if err := createCmd.MarkFlagRequired("template"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}

	updateCmd.Flags().IntVarP(&opt.NodeCount, "count", "", 0, "node count")
	updateCmd.Flags().IntVarP(&opt.MinNodeCount, "min", "", 0, "max node count")
	updateCmd.Flags().IntVarP(&opt.MaxNodeCount, "max", "", 0, "min node count")
	updateCmd.Flags().IntVarP(&opt.VolumeSize, "volume-size", "", 0, "volume-size")
	updateCmd.Flags().StringVarP(&opt.Flavor, "flavor", "f", "", "flavor")
	updateCmd.Flags().StringVarP(&opt.VolumeType, "volume-type", "", "", "volume-type")
	updateCmd.Flags().StringVarP(&opt.Name, "name", "n", "", "name")

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
	if err := viper.Unmarshal(&opt); err != nil {
		log.Fatal(err)
	}
}
