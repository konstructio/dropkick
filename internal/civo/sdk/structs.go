package sdk

import (
	"context"
	"errors"
)

// Civoer is the interface that wraps the Do method.
type Civoer interface {
	Do(ctx context.Context, location, method string, output interface{}, params map[string]string) error
	GetRegion() string
}

// ErrNotFound is returned when an item is not found.
var ErrNotFound = errors.New("not found")

// Resource represents any of the Civo resources returned
// by the Civo API.
type Resource interface {
	Instance | Firewall | Volume | KubernetesCluster | Network | ObjectStore | ObjectStoreCredential | SSHKey | LoadBalancer
	APIResource
}

// APIResource ensures all types from this SDK have the required
// methods to be addressed via API.
type APIResource interface {
	GetID() string
	GetName() string
	GetAPIEndpoint() string
	IsSinglePaged() bool
	GetResourceType() string
}

var (
	_ APIResource = &Instance{}
	_ APIResource = &Firewall{}
	_ APIResource = &Volume{}
	_ APIResource = &KubernetesCluster{}
	_ APIResource = &Network{}
	_ APIResource = &ObjectStore{}
	_ APIResource = &ObjectStoreCredential{}
	_ APIResource = &SSHKey{}
	_ APIResource = &LoadBalancer{}
)

// Instance is a Civo instance.
type Instance struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Hostname   string `json:"hostname"`
	FirewallID string `json:"firewall_id"`
	NetworkID  string `json:"network_id"`
	SSHKeyID   string `json:"ssh_key_id,omitempty"` // ssh_key_id is not available within a KubernetesCluster: they currently don't use SSH keys
	Status     string `json:"status"`
}

func (i Instance) GetID() string           { return i.ID }
func (i Instance) GetName() string         { return i.Name }
func (i Instance) GetAPIEndpoint() string  { return "/v2/instances" }
func (i Instance) IsSinglePaged() bool     { return false }
func (i Instance) GetResourceType() string { return "instance" }

// Firewall is a Civo firewall.
type Firewall struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	InstanceCount     int    `json:"instance_count"`
	ClusterCount      int    `json:"cluster_count"`
	LoadBalancerCount int    `json:"load_balancer_count"`
	NetworkID         string `json:"network_id"`
}

func (f Firewall) GetID() string           { return f.ID }
func (f Firewall) GetName() string         { return f.Name }
func (f Firewall) GetAPIEndpoint() string  { return "/v2/firewalls" }
func (f Firewall) IsSinglePaged() bool     { return true }
func (f Firewall) GetResourceType() string { return "firewall" }

// Volume is a Civo volume.
type Volume struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	InstanceID string `json:"instance_id"`
	NetworkID  string `json:"network_id"`
	ClusterID  string `json:"cluster_id,omitempty"` // cluster_id is not available within a KubernetesCluster
	Status     string `json:"status"`
}

func (v Volume) GetID() string           { return v.ID }
func (v Volume) GetName() string         { return v.Name }
func (v Volume) GetAPIEndpoint() string  { return "/v2/volumes" }
func (v Volume) IsSinglePaged() bool     { return true }
func (v Volume) GetResourceType() string { return "volume" }

// KubernetesCluster is a Civo Kubernetes cluster.
type KubernetesCluster struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	Ready      bool       `json:"ready"`
	FirewallID string     `json:"firewall_id"`
	NetworkID  string     `json:"network_id"`
	Volumes    []Volume   `json:"volumes"`
	Instances  []Instance `json:"instances"`
}

func (k KubernetesCluster) GetID() string           { return k.ID }
func (k KubernetesCluster) GetName() string         { return k.Name }
func (k KubernetesCluster) GetAPIEndpoint() string  { return "/v2/kubernetes/clusters" }
func (k KubernetesCluster) IsSinglePaged() bool     { return false }
func (k KubernetesCluster) GetResourceType() string { return "kubernetes cluster" }

// Network is a Civo network.
type Network struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (n Network) GetID() string           { return n.ID }
func (n Network) GetName() string         { return n.Name }
func (n Network) GetAPIEndpoint() string  { return "/v2/networks" }
func (n Network) IsSinglePaged() bool     { return true }
func (n Network) GetResourceType() string { return "network" }

// ObjectStore is a Civo object store.
type ObjectStore struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Credentials struct {
		ID   string `json:"credential_id"`
		Name string `json:"name"`
	} `json:"owner_info"`
	Status string `json:"status"`
}

func (o ObjectStore) GetID() string           { return o.ID }
func (o ObjectStore) GetName() string         { return o.Name }
func (o ObjectStore) GetAPIEndpoint() string  { return "/v2/objectstores" }
func (o ObjectStore) IsSinglePaged() bool     { return false }
func (o ObjectStore) GetResourceType() string { return "object store" }

// ObjectStoreCredential is a Civo object store credential.
type ObjectStoreCredential struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (o ObjectStoreCredential) GetID() string           { return o.ID }
func (o ObjectStoreCredential) GetName() string         { return o.Name }
func (o ObjectStoreCredential) GetAPIEndpoint() string  { return "/v2/objectstore/credentials" }
func (o ObjectStoreCredential) IsSinglePaged() bool     { return false }
func (o ObjectStoreCredential) GetResourceType() string { return "object store credential" }

// SSHKey is a Civo SSH key.
type SSHKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

func (s SSHKey) GetID() string           { return s.ID }
func (s SSHKey) GetName() string         { return s.Name }
func (s SSHKey) GetAPIEndpoint() string  { return "/v2/sshkeys" }
func (s SSHKey) IsSinglePaged() bool     { return true }
func (s SSHKey) GetResourceType() string { return "ssh key" }

// LoadBalancer is a Civo load balancer.
type LoadBalancer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FirewallID string `json:"firewall_id"`
	ClusterID  string `json:"cluster_id"`
}

func (l LoadBalancer) GetID() string           { return l.ID }
func (l LoadBalancer) GetName() string         { return l.Name }
func (l LoadBalancer) GetAPIEndpoint() string  { return "/v2/loadbalancers" }
func (l LoadBalancer) IsSinglePaged() bool     { return true }
func (l LoadBalancer) GetResourceType() string { return "load balancer" }
