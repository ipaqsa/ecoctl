package commands

import (
	"context"
	"ecoctl/pkg/client"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"time"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "functions to interact with clusters",
}

var clusterListCmd = &cobra.Command{
	Use:     "list",
	Example: "ecoctl cluster list -p 27520 -r 8",
	Short:   "get the list of clusters by the region and project",
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, _, err := eco.Service.ClusterList(context.Background(), &clusterListOpt)
		if err != nil {
			log.Printf("err during getting the list of clusters")
			return err
		}
		fmt.Printf("There are clusters with %v IDs in the region %d in the project %d\n", ids, clusterListOpt.RegionID, clusterListOpt.ProjectID)
		return nil
	},
}

var clusterCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "create a new cluster",
	Example: "ecoctl cluster create -c cluster.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		if clusterCreateOpt.WorkerOpts == nil || clusterCreateOpt.MasterOpts.Flavor == "" {
			return errors.New("invalid template")
		}
		log.Printf("creating the cluster")
		id, _, err := eco.Service.CreateCluster(context.Background(), &clusterCreateOpt)
		if err != nil {
			log.Printf("err during creating the cluster")
			return err
		}
		fmt.Printf("cluster record with %d ID is created\n", id)
		return err
	},
}
var clusterGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "get the cluster by ID",
	Example: "ecoctl cluster get 3",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		clusterID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Printf("err during parsing the clusterID")
			return err
		}
		if clusterID <= 0 {
			return errors.New("clusterID must be set and greater than 0")
		}
		log.Printf("getting the cluster with %d ID", clusterID)
		cluster, _, err := eco.Service.GetCluster(context.Background(), uint(clusterID))
		if err != nil {
			log.Printf("err during getting the cluster")
			return err
		}
		var indent []byte
		if output == "yaml" {
			indent, err = yaml.Marshal(cluster)
			if err != nil {
				log.Printf("err during marshal")
				return err
			}
		} else {
			indent, err = json.MarshalIndent(cluster, "", "\t")
			if err != nil {
				log.Printf("err during marshal")
				return err
			}
		}
		fmt.Printf("%s\n", string(indent))
		return nil
	},
}
var clusterDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete the cluster by ID",
	Example: "ecoctl cluster delete 3",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		clusterID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Printf("err during parse clusterID")
			return err
		}
		if clusterID <= 0 {
			return errors.New("clusterID must be set and greater than 0")
		}
		log.Printf("deleting the cluster with %d ID", clusterID)
		_, err = eco.Service.DeleteCluster(context.Background(), uint(clusterID))
		if err != nil {
			log.Printf("err during deleting the cluster")
			return err
		}
		log.Printf("deleting cluster with %d ID started", clusterID)
		return nil
	},
}

var clusterWatchCmd = &cobra.Command{
	Use:     "watch",
	Short:   "watch the cluster`s state",
	Example: "ecoctl cluster watch 3",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		clusterID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Printf("err during parse clusterID")
			return err
		}
		if clusterID <= 0 {
			return errors.New("clusterID must be set and greater than 0")
		}
		var timer = time.NewTimer(time.Minute * 7)
		var ctx = context.Background()
		var cluster *client.Cluster
		for {
			cluster, _, err = eco.Service.GetCluster(ctx, uint(clusterID))
			if err != nil {
				log.Printf("err during getting the cluster")
				return err
			}
			if cluster.State == "cluster_updated" || cluster.State == "error" {
				break
			}
			log.Printf("watch the cluster %s with %d ID\nstate: %s\nstatus: %s\nexisted: %s\n", cluster.Name,
				clusterID, cluster.State, cluster.Status, cluster.Existed)
			time.Sleep(time.Second * 10)
			select {
			case <-timer.C:
				return errors.New(fmt.Sprintf("timeout occurred: now the cluster with %d ID is in %s state with %s status",
					clusterID, cluster.State, cluster.Status))
			default:
				continue
			}

		}
		fmt.Printf("the cluster with %d ID is in %s state with %s status\n", clusterID, cluster.State, cluster.Status)
		return nil
	},
}

var clusterConfigCmd = &cobra.Command{
	Use:     "config",
	Short:   "get the cluster`s config",
	Example: "ecoctl cluster config 3 -p ~/.kube/config",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		clusterID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Printf("err during parsing the clusterID")
			return err
		}
		if clusterID <= 0 {
			return errors.New("clusterID must be set and greater than 0")
		}
		log.Printf("getting the cluster config with %d ID", clusterID)
		conf, _, err := eco.Service.ClusterConfig(context.Background(), uint(clusterID))
		if err != nil {
			return err
		}
		file, err := os.OpenFile(pathToClusterConfig, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Printf("err during opening config`s file")
			return err
		}
		defer file.Close()
		_, err = file.WriteString(conf.Content)
		return err
	},
}
