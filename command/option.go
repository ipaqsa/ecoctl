package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/ipaqsa/ecogo"
)

type Opt struct {
	ClusterID uint
	PoolID    uint

	PathToConfig string

	Namespace string

	AdminTTL string

	EcoURL    string
	UserID    string
	ProjectID string
	RegionID  string
	APIKey    string
}

var ErrMarkRequired = errors.New("err during making required")

var GlobalOption Opt

func NewClient() (*ecogo.Client, error) {
	c, err := ecogo.New(nil, ecogo.SetURL(GlobalOption.EcoURL), ecogo.SetRegionID(GlobalOption.RegionID),
		ecogo.SetProjectID(GlobalOption.ProjectID), ecogo.SetAPIKey(GlobalOption.APIKey),
		ecogo.SetRequestHeaders(map[string]string{"User-ID": GlobalOption.UserID}))
	if err != nil {
		fmt.Printf("err init eco client")
		return nil, err
	}
	fmt.Printf("server: %s\n", GlobalOption.EcoURL)
	fmt.Printf("getting eco version\n")
	v, err := c.Clusters.ServerVersion(context.Background())
	if err != nil {
		fmt.Printf("err during getting server version")
		return nil, err
	}
	fmt.Printf("version: %s\n\n", v)
	return c, nil
}
