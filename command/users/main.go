package users

import (
	"context"
	"ecoctl/command"
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

func init() {
	if err := initFlags(); err != nil {
		log.Fatal(err)
	}
	UserCmd.AddCommand(requestUserConfigCmd)
	AdminCmd.AddCommand(requestAdminConfigCmd)
}

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "functions to interact with users",
}
var requestUserConfigCmd = &cobra.Command{
	Use:     "request",
	Example: "ecoctl user request 2 --role somerole -u somename -n default -s 120 -o cluster.conf",
	Short:   "request user config",
	RunE: func(cmd *cobra.Command, args []string) error {
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
			log.Printf("err during creating eco client")
			return err
		}
		log.Println("requesting config")
		conf, err := c.Users.RequestUserConfig(context.Background(), uint(clusterID), &userOpt)
		if err != nil {
			log.Println("err during get user config")
			return err
		}
		file, err := os.OpenFile(command.GlobalOption.PathToConfig, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Printf("err during opening config file")
			return err
		}
		defer file.Close()
		_, err = file.WriteString(conf.Content)
		return err
	},
}

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "functions to interact with admin",
}
var requestAdminConfigCmd = &cobra.Command{
	Use:     "request",
	Example: "ecoctl admin request 2 --ttl 1h -o cluster.conf",
	Short:   "request admin config",
	RunE: func(cmd *cobra.Command, args []string) error {
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
			log.Printf("err during creating eco client")
			return err
		}
		log.Println("requesting admin config")
		conf, err := c.Users.RequestAdminConfig(context.Background(), uint(clusterID), command.GlobalOption.AdminTTL)
		if err != nil {
			log.Println("err during getting admin config")
			return err
		}
		file, err := os.OpenFile(command.GlobalOption.PathToConfig, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Printf("err during opening config file")
			return err
		}
		defer file.Close()
		_, err = file.WriteString(conf.Content)
		return err
	},
}
