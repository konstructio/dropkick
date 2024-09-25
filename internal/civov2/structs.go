package civov2

import "errors"

// ErrNotFound is returned when an item is not found.
var ErrNotFound = errors.New("not found")

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

// Firewall is a Civo firewall.
type Firewall struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	InstanceCount     int    `json:"instance_count"`
	ClusterCount      int    `json:"cluster_count"`
	LoadBalancerCount int    `json:"load_balancer_count"`
	NetworkID         string `json:"network_id"`
}

// Volume is a Civo volume.
type Volume struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	InstanceID string `json:"instance_id"`
	NetworkID  string `json:"network_id"`
	ClusterID  string `json:"cluster_id,omitempty"` // cluster_id is not available within a KubernetesCluster
	Status     string `json:"status"`
}

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

// Network is a Civo network.
type Network struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

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

// ObjectStoreCredential is a Civo object store credential.
type ObjectStoreCredential struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// SSHKey is a Civo SSH key.
type SSHKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}
