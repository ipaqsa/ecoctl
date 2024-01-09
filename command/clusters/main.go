package clusters

import (
	"context"
	"ecoctl/command"
	"errors"
	"fmt"
	"github.com/ipaqsa/ecogo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"strconv"
	"strings"
	"time"
)

func init() {
	if err := initFlags(); err != nil {
		log.Fatal(err)
	}
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(watchCmd)
}

var Cmd = &cobra.Command{
	Use:   "cluster",
	Short: "functions to interact with clusters",
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "get list of clusters by the region and project",
	Example: "ecoctl cluster list",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := command.NewClient()
		if err != nil {
			log.Printf("err during creating eco client")
			log.Fatal(err)
		}
		log.Println("getting list of clusters")
		clusters, err := c.Clusters.List(context.Background())
		if err != nil {
			if strings.Contains(err.Error(), "cluster not found") {
				log.Println("no clusters")
				return
			}
			log.Println("err during getting list of clusters")
			log.Fatal(err)
		}
		log.Printf("found %d clusters\n\n", len(clusters))
		for _, c := range clusters {
			indent, err := yaml.Marshal(c)
			if err != nil {
				log.Println("err during marshaling")
				log.Fatal(err)
			}
			fmt.Printf("%s\n", string(indent))
		}
	},
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "create a new cluster",
	Example: "ecoctl cluster create -t cluster.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		parseConfig()
		c, err := command.NewClient()
		if err != nil {
			log.Println("err during creating eco client")
			return err
		}
		log.Printf("creating cluster")
		id, err := c.Clusters.Create(context.Background(), &createOpt)
		if err != nil {
			log.Println("err during creating cluster")
			return err
		}
		log.Printf("cluster with %d ID is created\n", id)
		return nil
	},
}
var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "get cluster by ID",
	Example: "ecoctl cluster get 3",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		clusterID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Println("err during parsing clusterID")
			return err
		}
		if clusterID <= 0 {
			return errors.New("clusterID must be greater than 0")
		}
		c, err := command.NewClient()
		if err != nil {
			log.Println("err during creating eco client")
			return err
		}
		log.Println("getting cluster")
		cluster, err := c.Clusters.Request(context.Background(), uint(clusterID))
		if err != nil {
			log.Println("err during getting cluster")
			return err
		}
		indent, err := yaml.Marshal(cluster)
		if err != nil {
			log.Println("err during marshaling")
			return err
		}
		fmt.Printf("%s\n", string(indent))
		return nil
	},
}
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete cluster by ID",
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
			return errors.New("clusterID must be greater than 0")
		}
		c, err := command.NewClient()
		if err != nil {
			log.Println("err during creating eco client")
			return err
		}
		log.Println("deleting cluster")
		if err := c.Clusters.Delete(context.Background(), uint(clusterID)); err != nil {
			log.Println("err during deleting cluster")
			return err
		}
		log.Printf("deleting cluster with %d ID is started\n", clusterID)
		return nil
	},
}

var watchCmd = &cobra.Command{
	Use:     "watch",
	Short:   "watch cluster state",
	Example: "ecoctl cluster watch 3",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		clusterID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Println("err during parse clusterID")
			return err
		}
		if clusterID <= 0 {
			return errors.New("clusterID must be greater than 0")
		}
		c, err := command.NewClient()
		if err != nil {
			log.Println("err during creating eco client")
			return err
		}
		var timer = time.NewTimer(time.Minute * 8)
		var ctx = context.Background()
		var cluster *ecogo.Cluster
		for {
			cluster, err = c.Clusters.Request(ctx, uint(clusterID))
			if err != nil {
				log.Println("err during getting the cluster")
				return err
			}
			if cluster.State == "cluster_updated" || cluster.State == "error" {
				break
			}
			log.Printf("watch cluster %s with %d ID\nstate: %s\nstatus: %s\nexisted: %s\n", cluster.Name,
				clusterID, cluster.State, cluster.Status, cluster.Existed)
			time.Sleep(time.Second * 20)
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
