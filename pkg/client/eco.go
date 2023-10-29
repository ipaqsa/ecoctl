package client

import (
	"context"
	"net/http"
	"path"
	"strconv"
	"time"
)

const (
	clusterBasePathV1  = "/v1/cluster"
	clustersBasePathV1 = "/v1/clusters"
	versionBasePathV1  = "/version"
)

var _ CloudOrchestrator = &CloudOrchestratorOp{}

type CloudOrchestrator interface {
	Version(ctx context.Context) (string, *Response, error)

	ClusterList(ctx context.Context, opt *ClusterListOpt) ([]uint, *Response, error)

	//ClusterHealth(ctx context.Context, clusterID uint) (*Response, error)
	ClusterConfig(ctx context.Context, clusterID uint) (*ClusterConfig, *Response, error)
	//ClusterToHA(ctx context.Context, clusterID uint) (*Response, error)

	GetCluster(ctx context.Context, clusterID uint) (*Cluster, *Response, error)
	CreateCluster(ctx context.Context, opt *ClusterCreateOpt) (uint, *Response, error)
	DeleteCluster(ctx context.Context, clusterID uint) (*Response, error)

	CreatePool(ctx context.Context, clusterID uint, opt *PoolOpt) (*Response, error)
	UpdatePool(ctx context.Context, clusterID, poolID uint, opt *PoolOpt) (*Response, error)
	DeletePool(ctx context.Context, clusterID, poolID uint) (*Response, error)
}

type CloudOrchestratorOp struct {
	client *Client
}

type ClusterListOpt struct {
	ProjectID uint `url:"project,omitempty"`
	RegionID  uint `url:"region,omitempty"`
}

type ClusterCreateOpt struct {
	ProjectID   int       `json:"projectID" yaml:"projectID"`
	RegionID    int       `json:"regionID" yaml:"regionID"`
	HA          bool      `json:"ha" yaml:"ha"`
	InternalLB  bool      `json:"internalLB" yaml:"internalLB"`
	Name        string    `json:"name" yaml:"name"`
	Version     string    `json:"version" yaml:"version"`
	APILbFlavor string    `json:"apiLbFlavor" yaml:"APILbFlavor"`
	NetworkID   string    `json:"networkID" yaml:"networkID"`
	SubnetID    string    `json:"subnetID" yaml:"subnetID"`
	MasterOpts  PoolOpt   `json:"masterOpts" yaml:"masterOpts"`
	WorkerOpts  []PoolOpt `json:"workerOpts" yaml:"workerOpts"`
}
type PoolOpt struct {
	NodeCount      int    `json:"nodeCount" yaml:"nodeCount"`
	MaxNodeCount   int    `json:"maxNodeCount" yaml:"maxNodeCount"`
	MinNodeCount   int    `json:"minNodeCount" yaml:"minNodeCount"`
	EtcdVolumeSize int    `json:"etcdVolumeSize" yaml:"etcdVolumeSize"`
	VolumeSize     int    `json:"volumeSize" yaml:"volumeSize"`
	SetTaint       bool   `json:"set-taint" yaml:"setTaint"`
	K8sRole        string `json:"k8s-role" yaml:"k8sRole"`
	VolumeType     string `json:"volumeType" yaml:"volumeType"`
	Name           string `json:"name" yaml:"name"`
	Flavor         string `json:"flavor" yaml:"flavor"`
}

type ClusterConfig struct {
	Content string `json:"content"`
}

type Cluster struct {
	ID             uint      `json:"id"`
	ProjectID      int       `json:"projectID"`
	RegionID       int       `json:"regionID"`
	Retry          int       `json:"retry"`
	Processing     bool      `json:"processing"`
	HA             bool      `json:"ha"`
	InternalLB     bool      `json:"internalLB"`
	Version        string    `json:"version"`
	Name           string    `json:"name"`
	NetworkID      string    `json:"networkID"`
	SubnetID       string    `json:"subnetID"`
	State          string    `json:"state"`
	Status         string    `json:"status"`
	LoadBalancerID string    `json:"loadBalancerID"`
	Created        time.Time `json:"created"`
	Existed        string    `json:"existed"`
	MastersPool    Pool      `json:"mastersPool"`
	WorkersPools   []Pool    `json:"workersPools"`
}
type Pool struct {
	ID             uint   `json:"id"`
	NodeCount      int    `json:"nodeCount"`
	MaxNodeCount   int    `json:"maxNodeCount"`
	MinNodeCount   int    `json:"minNodeCount"`
	EtcdVolumeSize int    `json:"etcdVolumeSize"`
	VolumeSize     int    `json:"volumeSize"`
	VolumeType     string `json:"volumeType"`
	Name           string `json:"name"`
	Flavor         string `json:"flavor"`
	Role           string `json:"k8s-role"`
	SetTaint       bool   `json:"set-k8s-taint"`
	State          string `json:"state"`
	Status         string `json:"status"`
}

func (ec *CloudOrchestratorOp) Version(ctx context.Context) (string, *Response, error) {
	req, err := ec.client.NewRequest(ctx, http.MethodGet, versionBasePathV1, nil)
	if err != nil {
		return "", nil, err
	}
	var version string
	resp, err := ec.client.Do(ctx, req, &version)
	if err != nil {
		return "", resp, err
	}
	return version, resp, nil
}

func (ec *CloudOrchestratorOp) ClusterList(ctx context.Context, opt *ClusterListOpt) ([]uint, *Response, error) {
	project := strconv.FormatUint(uint64(opt.ProjectID), 10)
	region := strconv.FormatUint(uint64(opt.RegionID), 10)
	requestPath := path.Join(clustersBasePathV1, project, region)

	req, err := ec.client.NewRequest(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return nil, nil, err
	}

	var clusterIDs []uint
	resp, err := ec.client.Do(ctx, req, &clusterIDs)
	if err != nil {
		return nil, resp, err
	}
	return clusterIDs, resp, nil
}

func (ec *CloudOrchestratorOp) CreateCluster(ctx context.Context, opt *ClusterCreateOpt) (uint, *Response, error) {
	req, err := ec.client.NewRequest(ctx, http.MethodPost, clusterBasePathV1, opt)
	if err != nil {
		return 0, nil, err
	}

	var id uint
	resp, err := ec.client.Do(ctx, req, &id)
	if err != nil {
		return 0, resp, err
	}
	return id, resp, err
}
func (ec *CloudOrchestratorOp) GetCluster(ctx context.Context, clusterID uint) (*Cluster, *Response, error) {
	id := strconv.FormatUint(uint64(clusterID), 10)
	requestPath := path.Join(clusterBasePathV1, id)

	req, err := ec.client.NewRequest(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return nil, nil, err
	}

	var cluster Cluster
	resp, err := ec.client.Do(ctx, req, &cluster)
	if err != nil {
		return nil, resp, err
	}
	return &cluster, resp, nil
}
func (ec *CloudOrchestratorOp) DeleteCluster(ctx context.Context, clusterID uint) (*Response, error) {
	id := strconv.FormatUint(uint64(clusterID), 10)
	requestPath := path.Join(clusterBasePathV1, id)

	req, err := ec.client.NewRequest(ctx, http.MethodDelete, requestPath, nil)
	if err != nil {
		return nil, err
	}
	var status string
	resp, err := ec.client.Do(ctx, req, &status)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (ec *CloudOrchestratorOp) ClusterConfig(ctx context.Context, clusterID uint) (*ClusterConfig, *Response, error) {
	id := strconv.FormatUint(uint64(clusterID), 10)
	requestPath := path.Join(clusterBasePathV1, id, "config")

	req, err := ec.client.NewRequest(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return nil, nil, err
	}

	var config ClusterConfig
	resp, err := ec.client.Do(ctx, req, &config)
	if err != nil {
		return nil, resp, err
	}
	return &config, resp, nil
}

func (ec *CloudOrchestratorOp) CreatePool(ctx context.Context, clusterID uint, opt *PoolOpt) (*Response, error) {
	id := strconv.FormatUint(uint64(clusterID), 10)
	requestPath := path.Join(clusterBasePathV1, id, "pool")

	req, err := ec.client.NewRequest(ctx, http.MethodPost, requestPath, opt)
	if err != nil {
		return nil, err
	}

	var status string
	resp, err := ec.client.Do(ctx, req, &status)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (ec *CloudOrchestratorOp) UpdatePool(ctx context.Context, clusterID, poolID uint, opt *PoolOpt) (*Response, error) {
	cid := strconv.FormatUint(uint64(clusterID), 10)
	pid := strconv.FormatUint(uint64(poolID), 10)
	requestPath := path.Join(clusterBasePathV1, cid, "pool", pid)

	req, err := ec.client.NewRequest(ctx, http.MethodPatch, requestPath, opt)
	if err != nil {
		return nil, err
	}

	var status string
	resp, err := ec.client.Do(ctx, req, &status)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (ec *CloudOrchestratorOp) DeletePool(ctx context.Context, clusterID, poolID uint) (*Response, error) {
	cid := strconv.FormatUint(uint64(clusterID), 10)
	pid := strconv.FormatUint(uint64(poolID), 10)
	requestPath := path.Join(clusterBasePathV1, cid, "pool", pid)

	req, err := ec.client.NewRequest(ctx, http.MethodDelete, requestPath, nil)
	if err != nil {
		return nil, err
	}

	var status string
	resp, err := ec.client.Do(ctx, req, &status)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
