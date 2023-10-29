package commands

import (
	"ecoctl/pkg/client"
	"errors"
)

var clusterListOpt client.ClusterListOpt
var clusterCreateOpt client.ClusterCreateOpt
var poolCreateOpt client.PoolOpt

var pathToClusterCreationTemplate string
var pathToPoolCreationTemplate string

var pathToClusterConfig string

var clusterID uint

var ErrMarkRequired = errors.New("err during making required")

func initClusterListFlags() error {
	clusterListCmd.Flags().UintVarP(&clusterListOpt.RegionID, "region", "r", 0, "regionID")
	if err := clusterListCmd.MarkFlagRequired("region"); err != nil {
		return errors.Join(ErrMarkRequired, err)
	}
	clusterListCmd.Flags().UintVarP(&clusterListOpt.ProjectID, "project", "p", 0, "projectID")
	if err := clusterListCmd.MarkFlagRequired("project"); err != nil {
		return errors.Join(ErrMarkRequired, err)
	}
	return nil
}
func initClusterCreateFlags() error {
	clusterCreateCmd.Flags().StringVarP(&pathToClusterCreationTemplate, "template", "p", "",
		"template with the cluster configuration")
	if err := clusterCreateCmd.MarkFlagRequired("template"); err != nil {
		return errors.Join(ErrMarkRequired, err)
	}
	return nil
}
func initClusterConfigFlags() error {
	clusterConfigCmd.Flags().StringVarP(&pathToClusterConfig, "path", "p", "", "path to cluster config")
	if err := clusterConfigCmd.MarkFlagRequired("path"); err != nil {
		return errors.Join(ErrMarkRequired, err)
	}
	return nil
}

func initPoolCreateFlags() error {
	poolCreateCmd.Flags().UintVarP(&clusterID, "cluster", "c", 0, "clusterID")
	if err := poolCreateCmd.MarkFlagRequired("cluster"); err != nil {
		return errors.Join(ErrMarkRequired, err)
	}
	poolCreateCmd.Flags().StringVarP(&pathToPoolCreationTemplate, "template", "p", "",
		"template with the pool configuration")
	if err := poolCreateCmd.MarkFlagRequired("template"); err != nil {
		return errors.Join(ErrMarkRequired, err)
	}
	return nil
}
func initPoolDeleteFlags() error {
	poolDeleteCmd.Flags().UintVarP(&clusterID, "cluster", "c", 0, "clusterID")
	if err := poolDeleteCmd.MarkFlagRequired("cluster"); err != nil {
		return errors.Join(ErrMarkRequired, err)
	}
	return nil
}
func initPoolUpdateFlags() error {
	poolUpdateCmd.Flags().UintVarP(&clusterID, "cluster", "c", 0, "clusterID")
	if err := poolUpdateCmd.MarkFlagRequired("cluster"); err != nil {
		return errors.Join(ErrMarkRequired, err)
	}
	poolUpdateCmd.Flags().IntVarP(&poolCreateOpt.NodeCount, "count", "", 0, "node count")
	poolUpdateCmd.Flags().IntVarP(&poolCreateOpt.MinNodeCount, "min", "", 0, "max node count")
	poolUpdateCmd.Flags().IntVarP(&poolCreateOpt.MaxNodeCount, "max", "", 0, "min node count")
	poolUpdateCmd.Flags().IntVarP(&poolCreateOpt.VolumeSize, "volume-size", "", 0, "volume-size")
	poolUpdateCmd.Flags().StringVarP(&poolCreateOpt.Flavor, "flavor", "f", "", "flavor")
	poolUpdateCmd.Flags().StringVarP(&poolCreateOpt.VolumeType, "volume-type", "", "", "volume-type")
	poolUpdateCmd.Flags().StringVarP(&poolCreateOpt.Name, "name", "n", "", "name")
	return nil
}
