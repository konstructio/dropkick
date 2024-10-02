package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/konstructio/dropkick/internal/civo/sdk/json"
)

// GetInstance gets an instance by ID.
func (c *Client) GetInstance(ctx context.Context, instanceID string) (*Instance, error) {
	ins := &Instance{ID: instanceID}
	err := getByID(ctx, c, ins)
	return ins, err
}

// GetFirewall gets a firewall by ID.
func (c *Client) GetFirewall(ctx context.Context, firewallID string) (*Firewall, error) {
	fw := &Firewall{ID: firewallID}
	err := getByID(ctx, c, fw)
	return fw, err
}

// GetVolume gets a volume by ID.
func (c *Client) GetVolume(ctx context.Context, volumeID string) (*Volume, error) {
	vol := &Volume{ID: volumeID}
	err := getByID(ctx, c, vol)
	return vol, err
}

// GetKubernetesCluster gets a Kubernetes cluster by ID.
func (c *Client) GetKubernetesCluster(ctx context.Context, clusterID string) (*KubernetesCluster, error) {
	cluster := &KubernetesCluster{ID: clusterID}
	err := getByID(ctx, c, cluster)
	return cluster, err
}

// GetNetwork gets a network by ID.
func (c *Client) GetNetwork(ctx context.Context, networkID string) (*Network, error) {
	net := &Network{ID: networkID}
	err := getByID(ctx, c, net)
	return net, err
}

// GetObjectStore gets an object store by ID.
func (c *Client) GetObjectStore(ctx context.Context, objstoreID string) (*ObjectStore, error) {
	objstore := &ObjectStore{ID: objstoreID}
	err := getByID(ctx, c, objstore)
	return objstore, err
}

// GetObjectStoreCredential gets an object store credential by ID.
func (c *Client) GetObjectStoreCredential(ctx context.Context, objstorecredID string) (*ObjectStoreCredential, error) {
	objstorecred := &ObjectStoreCredential{ID: objstorecredID}
	err := getByID(ctx, c, objstorecred)
	return objstorecred, err
}

// GetLoadBalancer gets a load balancer by ID.
func (c *Client) GetLoadBalancer(ctx context.Context, lbID string) (*LoadBalancer, error) {
	lb := &LoadBalancer{ID: lbID}
	err := getByID(ctx, c, lb)
	return lb, err
}

// GetSSHKey gets an SSH key by ID.
func (c *Client) GetSSHKey(ctx context.Context, sshkeyID string) (*SSHKey, error) {
	sshkey := &SSHKey{ID: sshkeyID}
	err := getByID(ctx, c, sshkey)
	return sshkey, err
}

// EmptyIDError is returned when the ID field in a resource is empty.
type EmptyIDError struct {
	ResourceType string
}

// Error returns the error message, by implementing the error interface.
func (e *EmptyIDError) Error() string {
	return fmt.Sprintf("the ID field in the resource %q is empty", e.ResourceType)
}

// Is checks if the target error is an EmptyIDError.
func (e *EmptyIDError) Is(target error) bool {
	_, ok := target.(*EmptyIDError)
	return ok
}

// GetByID gets a resource by ID as provided by the value in "resource".
// The ID value in resource must not be empty.
func getByID[T Resource](ctx context.Context, c Civoer, resource *T) error {
	params := map[string]string{"region": c.GetRegion()}

	endpoint := (*resource).GetAPIEndpoint()
	restype := (*resource).GetResourceType()
	id := (*resource).GetID()

	if id == "" {
		return &EmptyIDError{ResourceType: restype}
	}

	fullpath := path.Join(endpoint, id)
	if err := c.Do(ctx, fullpath, http.MethodGet, resource, params); err != nil {
		if errors.Is(err, &json.HTTPError{Code: http.StatusNotFound}) {
			return ErrNotFound
		}

		return fmt.Errorf("unable to get item: %w", err)
	}

	return nil
}
