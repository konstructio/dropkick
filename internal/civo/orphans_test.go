package civo

import (
	"context"
	"os"
	"testing"

	"github.com/konstructio/dropkick/internal/civo/sdk"
	"github.com/konstructio/dropkick/internal/civo/sdk/testutils"
	"github.com/konstructio/dropkick/internal/logger"
)

func Test_NukeOrphanedResources(t *testing.T) {
	instances := []sdk.Instance{{
		ID:        "1",
		Name:      "test-instance-1",
		SSHKeyID:  "1", // uses an existing SSH key
		NetworkID: "1", // uses an existing network
	}, {
		ID:   "2",
		Name: "test-instance-2",
	}}

	volumes := []sdk.Volume{{
		ID:        "1",
		Name:      "test-volume-1",
		Status:    "attached", // attached to an existing instance
		NetworkID: "1",        // uses an existing network
	}, {
		ID:   "2",
		Name: "test-volume-2",
	}}

	loadbalancers := []sdk.LoadBalancer{{
		ID:         "1",
		Name:       "test-loadbalancer-1",
		ClusterID:  "1", // uses an existing kubernetes cluster
		FirewallID: "1", // uses an existing firewall
	}, {
		ID:   "2",
		Name: "test-loadbalancer-2",
	}}

	objectstores := []sdk.ObjectStore{{
		ID:          "1",
		Name:        "test-objectstore-1",
		Credentials: sdk.ObjectStoreCredential{ID: "1"}, // uses an existing object store credential
	}, {
		ID:   "2",
		Name: "test-objectstore-2",
	}}

	objectstorecreds := []sdk.ObjectStoreCredential{{
		ID:   "1",
		Name: "test-objectstore-credential-1",
	}, {
		ID:   "2",
		Name: "test-objectstore-credential-2",
	}}

	sshkeys := []sdk.SSHKey{{
		ID:   "1",
		Name: "test-sshkey-1",
	}, {
		ID:   "2",
		Name: "test-sshkey-2",
	}}

	networks := []sdk.Network{{
		ID:    "1",
		Name:  "cust-test-network-1-219873hudsf",
		Label: "test-network-1",
	}, {
		ID:    "2",
		Name:  "cust-test-network-2-894359435jk",
		Label: "test-network-2",
	}}

	firewalls := []sdk.Firewall{{
		ID:                "1",
		Name:              "test-firewall-1",
		ClusterCount:      1,   // attached to an existing kubernetes cluster
		InstanceCount:     6,   // attached to existing instances
		LoadBalancerCount: 2,   // attached to existing load balancers
		NetworkID:         "1", // uses an existing network
	}, {
		ID:   "2",
		Name: "test-firewall-2",
	}}

	mock := &mockClient{
		fnGetInstances:              func(ctx context.Context) ([]sdk.Instance, error) { return instances, nil },
		fnGetVolumes:                func(ctx context.Context) ([]sdk.Volume, error) { return volumes, nil },
		fnGetLoadBalancers:          func(ctx context.Context) ([]sdk.LoadBalancer, error) { return loadbalancers, nil },
		fnGetObjectStores:           func(ctx context.Context) ([]sdk.ObjectStore, error) { return objectstores, nil },
		fnGetObjectStoreCredentials: func(ctx context.Context) ([]sdk.ObjectStoreCredential, error) { return objectstorecreds, nil },
		fnGetSSHKeys:                func(ctx context.Context) ([]sdk.SSHKey, error) { return sshkeys, nil },
		fnGetNetworks:               func(ctx context.Context) ([]sdk.Network, error) { return networks, nil },
		fnGetFirewalls:              func(ctx context.Context) ([]sdk.Firewall, error) { return firewalls, nil },
	}

	civo := &Civo{
		client: mock,
		logger: logger.New(os.Stderr),
	}

	err := civo.NukeOrphanedResources(context.Background())
	testutils.AssertNoErrorf(t, err, "expected no error when calling NukeOrphanedResources, got %v", err)
}
