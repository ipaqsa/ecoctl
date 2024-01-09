package pools

import (
	"context"
	"ecoctl/command"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

func init() {
	if err := initFlags(); err != nil {
		log.Fatal(err)
	}
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)
}

var Cmd = &cobra.Command{
	Use:   "pool",
	Short: "functions to interact with pools",
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "create a new pool in the cluster",
	Example: "ecoctl pool create -c 1 -t pool.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		parseConfig()
		c, err := command.NewClient()
		if err != nil {
			log.Println("err during creating eco client")
			return err
		}
		log.Println("creating pool")
		if err := c.Pools.Create(context.Background(), command.GlobalOption.ClusterID, &opt); err != nil {
			log.Println("err during creating pool")
			return err
		}
		fmt.Println("creating pool is started")
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "update the pool in the cluster",
	Example: "ecoctl pool update 2 -c 1 --count 2",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		poolID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Println("err during parsing the poolID")
			return err
		}
		if poolID <= 0 {
			return errors.New("poolID must greater than 0")
		}
		c, err := command.NewClient()
		if err != nil {
			log.Println("err during creating eco client")
			return err
		}
		log.Println("updating pool")
		if err := c.Pools.Update(context.Background(), command.GlobalOption.ClusterID, uint(poolID), &opt); err != nil {
			log.Println("err during updating pool")
			return err
		}
		fmt.Printf("updating pool %d ID in cluster %d is started\n", poolID, command.GlobalOption.ClusterID)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete the pool in the cluster by ID",
	Example: "ecoctl pool delete 2 -c 1",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		poolID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Println("err during parsing the poolID")
			return err
		}
		if poolID <= 0 {
			return errors.New("poolID must be greater than 0")
		}
		c, err := command.NewClient()
		if err != nil {
			log.Println("err during creating eco client")
			return err
		}
		log.Println("deleting pool")
		if err = c.Pools.Delete(context.Background(), command.GlobalOption.ClusterID, uint(poolID)); err != nil {
			log.Println("err during deleting pool")
			return err
		}
		log.Printf("deleting pool with %d ID in the cluster with %d is started\n", poolID, command.GlobalOption.ClusterID)
		return nil
	},
}
