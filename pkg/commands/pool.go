package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var poolCmd = &cobra.Command{
	Use:   "pool",
	Short: "functions to interact with pools",
}

var poolCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "create a new pool in the cluster",
	Example: "ecoctl pool create -c 1 -p pool.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := eco.Service.CreatePool(context.Background(), clusterID, &poolCreateOpt)
		if err != nil {
			return err
		}
		fmt.Printf("pool record is created\n")
		return nil
	},
}

var poolUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "update the pool in the cluster",
	Example: "ecoctl pool update 2 -c 1 --count 2",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		poolID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Printf("err during parsing the poolID")
			return err
		}
		if poolID <= 0 {
			return errors.New("poolID must be set and greater than 0")
		}
		_, err = eco.Service.UpdatePool(context.Background(), clusterID, uint(poolID), &poolCreateOpt)
		if err != nil {
			return err
		}
		fmt.Printf("pool record with %d ID in cluster %d is updated\n", poolID, clusterID)
		return nil
	},
}

var poolDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete the pool in the cluster by ID",
	Example: "ecoctl pool delete 2 -c 1",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("err during parse args")
		}
		poolID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Printf("err during parsing the poolID")
			return err
		}
		if poolID <= 0 {
			return errors.New("poolID must be set and greater than 0")
		}
		_, err = eco.Service.DeletePool(context.Background(), clusterID, uint(poolID))
		if err != nil {
			return err
		}
		log.Printf("deleting the pool with %d ID in the cluster with %d started", poolID, clusterID)
		return nil
	},
}
