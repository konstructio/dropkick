package sdk

import (
	"errors"
	"strings"
)

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

// Compile-time assertions for each type implementing the APIResource interface.
// This ensures that the types are correctly implemented.
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

func (i Instance) GetID() string               { return i.ID }            // GetID returns the ID of the instance.
func (i Instance) GetName() string             { return i.Name }          // GetName returns the name of the instance.
func (i Instance) GetAPIEndpoint() string      { return "/v2/instances" } // GetAPIEndpoint returns the API endpoint for instances.
func (i Instance) IsSinglePaged() bool         { return false }           // IsSinglePaged returns whether the resource is single paged.
func (i Instance) GetResourceType() string     { return "instance" }      // GetResourceType returns the type of the resource.
func (i Instance) ConsumeOtherResources() bool { return false }           // ConsumeOtherResources returns whether the resource blocks deletion of others.

// Firewall is a Civo firewall.
type Firewall struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	InstanceCount     int    `json:"instance_count"`
	ClusterCount      int    `json:"cluster_count"`
	LoadBalancerCount int    `json:"load_balancer_count"`
	NetworkID         string `json:"network_id"`
}

func (f Firewall) GetID() string           { return f.ID }            // GetID returns the ID of the firewall.
func (f Firewall) GetName() string         { return f.Name }          // GetName returns the name of the firewall.
func (f Firewall) GetAPIEndpoint() string  { return "/v2/firewalls" } // GetAPIEndpoint returns the API endpoint for firewalls.
func (f Firewall) IsSinglePaged() bool     { return true }            // IsSinglePaged returns whether the resource is single paged.
func (f Firewall) GetResourceType() string { return "firewall" }      // GetResourceType returns the type of the resource.

// Volume is a Civo volume.
type Volume struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	InstanceID string `json:"instance_id"`
	NetworkID  string `json:"network_id"`
	ClusterID  string `json:"cluster_id,omitempty"` // cluster_id is not available within a KubernetesCluster
	Status     string `json:"status"`
}

func (v Volume) GetID() string           { return v.ID }          // GetID returns the ID of the volume.
func (v Volume) GetName() string         { return v.Name }        // GetName returns the name of the volume.
func (v Volume) GetAPIEndpoint() string  { return "/v2/volumes" } // GetAPIEndpoint returns the API endpoint for volumes.
func (v Volume) IsSinglePaged() bool     { return true }          // IsSinglePaged returns whether the resource is single paged.
func (v Volume) GetResourceType() string { return "volume" }      // GetResourceType returns the type of the resource.

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

func (k KubernetesCluster) GetID() string           { return k.ID }                      // GetID returns the ID of the Kubernetes cluster.
func (k KubernetesCluster) GetName() string         { return k.Name }                    // GetName returns the name of the Kubernetes cluster.
func (k KubernetesCluster) GetAPIEndpoint() string  { return "/v2/kubernetes/clusters" } // GetAPIEndpoint returns the API endpoint for Kubernetes clusters.
func (k KubernetesCluster) IsSinglePaged() bool     { return false }                     // IsSinglePaged returns whether the resource is single paged.
func (k KubernetesCluster) GetResourceType() string { return "kubernetes cluster" }      // GetResourceType returns the type of the resource.

// Network is a Civo network.
type Network struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (n Network) GetID() string           { return n.ID }           // GetID returns the ID of the network.
func (n Network) GetName() string         { return n.Name }         // GetName returns the name of the network.
func (n Network) GetAPIEndpoint() string  { return "/v2/networks" } // GetAPIEndpoint returns the API endpoint for networks.
func (n Network) IsSinglePaged() bool     { return true }           // IsSinglePaged returns whether the resource is single paged.
func (n Network) GetResourceType() string { return "network" }      // GetResourceType returns the type of the resource.

// ObjectStore is a Civo object store.
type ObjectStore struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Credentials ObjectStoreCredential `json:"owner_info"`
	Status      string                `json:"status"`
}

func (o ObjectStore) GetID() string           { return o.ID }               // GetID returns the ID of the object store.
func (o ObjectStore) GetName() string         { return o.Name }             // GetName returns the name of the object store.
func (o ObjectStore) GetAPIEndpoint() string  { return "/v2/objectstores" } // GetAPIEndpoint returns the API endpoint for object stores.
func (o ObjectStore) IsSinglePaged() bool     { return false }              // IsSinglePaged returns whether the resource is single paged.
func (o ObjectStore) GetResourceType() string { return "object store" }     // GetResourceType returns the type of the resource.

// ObjectStoreCredential is a Civo object store credential.
type ObjectStoreCredential struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (o ObjectStoreCredential) GetID() string           { return o.ID }                          // GetID returns the ID of the object store credential.
func (o ObjectStoreCredential) GetName() string         { return o.Name }                        // GetName returns the name of the object store credential.
func (o ObjectStoreCredential) GetAPIEndpoint() string  { return "/v2/objectstore/credentials" } // GetAPIEndpoint returns the API endpoint for object store credentials.
func (o ObjectStoreCredential) IsSinglePaged() bool     { return false }                         // IsSinglePaged returns whether the resource is single paged.
func (o ObjectStoreCredential) GetResourceType() string { return "object store credential" }     // GetResourceType returns the type of the resource.

// SSHKey is a Civo SSH key.
type SSHKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

func (s SSHKey) GetID() string           { return s.ID }          // GetID returns the ID of the SSH key.
func (s SSHKey) GetName() string         { return s.Name }        // GetName returns the name of the SSH key.
func (s SSHKey) GetAPIEndpoint() string  { return "/v2/sshkeys" } // GetAPIEndpoint returns the API endpoint for SSH keys.
func (s SSHKey) IsSinglePaged() bool     { return true }          // IsSinglePaged returns whether the resource is single paged.
func (s SSHKey) GetResourceType() string { return "ssh key" }     // GetResourceType returns the type of the resource.

// LoadBalancer is a Civo load balancer.
type LoadBalancer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FirewallID string `json:"firewall_id"`
	ClusterID  string `json:"cluster_id"`
}

func (l LoadBalancer) GetID() string           { return l.ID }                // GetID returns the ID of the load balancer.
func (l LoadBalancer) GetName() string         { return l.Name }              // GetName returns the name of the load balancer.
func (l LoadBalancer) GetAPIEndpoint() string  { return "/v2/loadbalancers" } // GetAPIEndpoint returns the API endpoint for load balancers.
func (l LoadBalancer) IsSinglePaged() bool     { return true }                // IsSinglePaged returns whether the resource is single paged.
func (l LoadBalancer) GetResourceType() string { return "load balancer" }     // GetResourceType returns the type of the resource.

func IsPaginatedResource(endpoint string) (bool, error) {
	resources := [...]APIResource{
		&Instance{},
		&Firewall{},
		&Volume{},
		&KubernetesCluster{},
		&Network{},
		&ObjectStore{},
		&ObjectStoreCredential{},
		&SSHKey{},
		&LoadBalancer{},
	}

	for _, resource := range resources {
		if strings.HasPrefix(endpoint, resource.GetAPIEndpoint()) {
			return !resource.IsSinglePaged(), nil
		}
	}

	return false, ErrNotFound
}
