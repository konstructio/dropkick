package civo

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/civo/civogo"
)

func TestDeleteVolumes(t *testing.T) {
	t.Run("successfully delete volumes", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		volumes := []civogo.Volume{{ID: "volume1", Name: "test-volume1"}}

		err := c.deleteVolumes(volumes)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("error deleting volume", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		volumes := []civogo.Volume{{ID: "volume1", Name: "test-volume1"}}

		err := c.deleteVolumes(volumes)
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("nuke not enabled", func(tt *testing.T) {
		mockClient := &mockCivoClient{}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   false,
		}

		volumes := []civogo.Volume{{ID: "volume1", Name: "test-volume1"}}

		err := c.deleteVolumes(volumes)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})
}

func TestDeleteSSHKeys(t *testing.T) {
	t.Run("successfully delete SSH keys", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnDeleteSSHKey: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		keys := []civogo.SSHKey{{ID: "key1", Name: "test-key1"}}

		err := c.deleteSSHKeys(keys)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("error deleting SSH key", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnDeleteSSHKey: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		keys := []civogo.SSHKey{{ID: "key1", Name: "test-key1"}}

		err := c.deleteSSHKeys(keys)
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("nuke not enabled", func(tt *testing.T) {
		mockClient := &mockCivoClient{}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   false,
		}

		keys := []civogo.SSHKey{{ID: "key1", Name: "test-key1"}}

		err := c.deleteSSHKeys(keys)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})
}

func TestDeleteNetworks(t *testing.T) {
	t.Run("successfully delete networks", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnDeleteNetwork: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		networks := []civogo.Network{{ID: "network1", Name: "test-network1"}}

		err := c.deleteNetworks(networks)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("error deleting network", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnDeleteNetwork: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		networks := []civogo.Network{{ID: "network1", Name: "test-network1"}}

		err := c.deleteNetworks(networks)
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("nuke not enabled", func(tt *testing.T) {
		mockClient := &mockCivoClient{}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   false,
		}

		networks := []civogo.Network{{ID: "network1", Name: "test-network1"}}

		err := c.deleteNetworks(networks)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})
}

func TestDeleteFirewalls(t *testing.T) {
	t.Run("successfully delete firewalls", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnDeleteFirewall: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		firewalls := []*civogo.Firewall{{ID: "firewall1", Name: "test-firewall1"}}

		err := c.deleteFirewalls(firewalls)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("error deleting firewall", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnDeleteFirewall: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		firewalls := []*civogo.Firewall{{ID: "firewall1", Name: "test-firewall1"}}

		err := c.deleteFirewalls(firewalls)
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("nuke not enabled", func(tt *testing.T) {
		mockClient := &mockCivoClient{}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   false,
		}

		firewalls := []*civogo.Firewall{{ID: "firewall1", Name: "test-firewall1"}}

		err := c.deleteFirewalls(firewalls)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})
}

func TestGetAllNodes(t *testing.T) {
	t.Run("successfully list all nodes from multiple pages", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				if page == 1 {
					return &civogo.PaginatedInstanceList{
						Items: []civogo.Instance{
							{ID: "node1", Hostname: "node-1"},
							{ID: "node2", Hostname: "node-2"},
						},
						Pages: 2,
					}, nil
				}
				if page == 2 {
					return &civogo.PaginatedInstanceList{
						Items: []civogo.Instance{
							{ID: "node3", Hostname: "node-3"},
						},
						Pages: 2,
					}, nil
				}
				return nil, errors.New("invalid page number")
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		nodes, err := c.getAllNodes()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(nodes) != 3 {
			tt.Errorf("expected 3 nodes, got %d", len(nodes))
		}

		if nodes[0].ID != "node1" || nodes[1].ID != "node2" || nodes[2].ID != "node3" {
			tt.Errorf("unexpected node IDs: %v", nodes)
		}
	})

	t.Run("error listing nodes on first page", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				if page == 1 {
					return nil, errTest
				}
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", Hostname: "node-1"},
					},
					Pages: 1,
				}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		_, err := c.getAllNodes()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("successfully list nodes when only one page exists", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				if page == 1 {
					return &civogo.PaginatedInstanceList{
						Items: []civogo.Instance{
							{ID: "node1", Hostname: "node-1"},
							{ID: "node2", Hostname: "node-2"},
						},
						Pages: 1,
					}, nil
				}
				return nil, errors.New("invalid page number")
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		nodes, err := c.getAllNodes()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(nodes) != 2 {
			tt.Errorf("expected 2 nodes, got %d", len(nodes))
		}
	})
}

func TestGetOrphanedVolumes(t *testing.T) {
	t.Run("successfully fetch orphaned volumes", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListVolumes: func() ([]civogo.Volume, error) {
				return []civogo.Volume{
					{ID: "volume1", Name: "orphan-volume1", Status: "available"},
					{ID: "volume2", Name: "attached-volume", Status: volumeStatusAttached, InstanceID: "node1"},
					{ID: "volume3", Name: "orphan-volume2", Status: "available"},
				}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		volumes, err := c.getOrphanedVolumes()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(volumes) != 2 {
			tt.Errorf("expected 2 orphaned volumes, got %d", len(volumes))
		}

		if volumes[0].ID != "volume1" || volumes[1].ID != "volume3" {
			tt.Errorf("unexpected orphaned volumes: %v", volumes)
		}
	})

	t.Run("error listing volumes", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListVolumes: func() ([]civogo.Volume, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		_, err := c.getOrphanedVolumes()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("no orphaned volumes", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListVolumes: func() ([]civogo.Volume, error) {
				return []civogo.Volume{
					{ID: "volume1", Name: "attached-volume1", Status: volumeStatusAttached, InstanceID: "node1"},
				}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		volumes, err := c.getOrphanedVolumes()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(volumes) != 0 {
			tt.Errorf("expected 0 orphaned volumes, got %d", len(volumes))
		}
	})
}

func TestGetOrphanedSSHKeys(t *testing.T) {
	t.Run("successfully fetch orphaned SSH keys", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return []civogo.SSHKey{
					{ID: "key1", Name: "key1"},
					{ID: "key2", Name: "key2"},
					{ID: "key3", Name: "key3"},
				}, nil
			},
		}

		nodes := []civogo.Instance{
			{ID: "node1", SSHKeyID: "key1"},
			{ID: "node2", SSHKeyID: "key2"},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		orphanedKeys, err := c.getOrphanedSSHKeys(nodes)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(orphanedKeys) != 1 {
			tt.Errorf("expected 1 orphaned SSH key, got %d", len(orphanedKeys))
		}

		if orphanedKeys[0].ID != "key3" {
			tt.Errorf("unexpected orphaned SSH key: %v", orphanedKeys[0])
		}
	})

	t.Run("error listing SSH keys", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return nil, errTest
			},
		}

		nodes := []civogo.Instance{
			{ID: "node1", SSHKeyID: "key1"},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		_, err := c.getOrphanedSSHKeys(nodes)
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("no orphaned SSH keys", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return []civogo.SSHKey{
					{ID: "key1", Name: "key1"},
					{ID: "key2", Name: "key2"},
				}, nil
			},
		}

		nodes := []civogo.Instance{
			{ID: "node1", SSHKeyID: "key1"},
			{ID: "node2", SSHKeyID: "key2"},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		orphanedKeys, err := c.getOrphanedSSHKeys(nodes)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(orphanedKeys) != 0 {
			tt.Errorf("expected 0 orphaned SSH keys, got %d", len(orphanedKeys))
		}
	})
}

func TestGetOrphanedNetworks(t *testing.T) {
	t.Run("successfully fetch orphaned networks", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListNetworks: func() ([]civogo.Network, error) {
				return []civogo.Network{
					{ID: "network1", Name: "network1"},
					{ID: "network2", Name: "network2"},
					{ID: "network3", Name: "network3"},
				}, nil
			},
		}

		nodes := []civogo.Instance{
			{ID: "node1", NetworkID: "network1"},
			{ID: "node2", NetworkID: "network2"},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		orphanedNetworks, err := c.getOrphanedNetworks(nodes)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(orphanedNetworks) != 1 {
			tt.Errorf("expected 1 orphaned network, got %d", len(orphanedNetworks))
		}

		if orphanedNetworks[0].ID != "network3" {
			tt.Errorf("unexpected orphaned network: %v", orphanedNetworks[0])
		}
	})

	t.Run("error listing networks", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListNetworks: func() ([]civogo.Network, error) {
				return nil, errTest
			},
		}

		nodes := []civogo.Instance{
			{ID: "node1", NetworkID: "network1"},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		_, err := c.getOrphanedNetworks(nodes)
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("no orphaned networks", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListNetworks: func() ([]civogo.Network, error) {
				return []civogo.Network{
					{ID: "network1", Name: "network1"},
					{ID: "network2", Name: "network2"},
				}, nil
			},
		}

		nodes := []civogo.Instance{
			{ID: "node1", NetworkID: "network1"},
			{ID: "node2", NetworkID: "network2"},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		orphanedNetworks, err := c.getOrphanedNetworks(nodes)
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(orphanedNetworks) != 0 {
			tt.Errorf("expected 0 orphaned networks, got %d", len(orphanedNetworks))
		}
	})
}

func TestGetOrphanedFirewalls(t *testing.T) {
	t.Run("successfully fetch orphaned firewalls", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) { return &civogo.PaginatedObjectstores{}, nil },
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{}, nil
			},
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return []civogo.Firewall{
					{ID: "firewall1", Name: "firewall1", ClusterCount: 0, InstanceCount: 0, LoadBalancerCount: 0},
					{ID: "firewall2", Name: "firewall2", ClusterCount: 1, InstanceCount: 0, LoadBalancerCount: 0},
					{ID: "firewall3", Name: "firewall3", ClusterCount: 0, InstanceCount: 0, LoadBalancerCount: 0},
				}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		orphanedFirewalls, err := c.getOrphanedFirewalls()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(orphanedFirewalls) != 2 {
			tt.Errorf("expected 2 orphaned firewalls, got %d", len(orphanedFirewalls))
		}

		if orphanedFirewalls[0].ID != "firewall1" || orphanedFirewalls[1].ID != "firewall3" {
			tt.Errorf("unexpected orphaned firewalls: %v", orphanedFirewalls)
		}
	})

	t.Run("error listing firewalls", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		_, err := c.getOrphanedFirewalls()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("no orphaned firewalls", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return []civogo.Firewall{
					{ID: "firewall1", Name: "firewall1", ClusterCount: 1, InstanceCount: 0, LoadBalancerCount: 0},
					{ID: "firewall2", Name: "firewall2", ClusterCount: 0, InstanceCount: 1, LoadBalancerCount: 0},
				}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		orphanedFirewalls, err := c.getOrphanedFirewalls()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}

		if len(orphanedFirewalls) != 0 {
			tt.Errorf("expected 0 orphaned firewalls, got %d", len(orphanedFirewalls))
		}
	})
}

// -----------------------------------------------------------------------------

func TestNukeOrphanedResources(t *testing.T) {
	// Test case for successfully nuking all orphaned resources
	t.Run("successfully nukes all orphaned resources", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", SSHKeyID: "key1", NetworkID: "network1"},
					},
					Pages: 1,
				}, nil
			},
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return &civogo.PaginatedObjectstores{
					Items: []civogo.ObjectStore{
						{ID: "objectstore1", Name: "objectstore1", OwnerInfo: civogo.BucketOwner{CredentialID: "credential1"}},
						{ID: "objectstore2", Name: "orphan-objectstore2", OwnerInfo: civogo.BucketOwner{CredentialID: "credential2"}},
					},
					Page: 1,
				}, nil
			},
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{
					Items: []civogo.ObjectStoreCredential{
						{ID: "credential1", Name: "credential1"},
						{ID: "credential2", Name: "orphan-credential2"},
						{ID: "credential3", Name: "credential3"},
					},
					Page: 1,
				}, nil
			},
			FnListVolumes: func() ([]civogo.Volume, error) {
				return []civogo.Volume{
					{ID: "volume1", Name: "orphan-volume1", Status: "available"},
				}, nil
			},
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return []civogo.SSHKey{
					{ID: "key1", Name: "key1"},
					{ID: "key2", Name: "orphan-key2"},
				}, nil
			},
			FnListNetworks: func() ([]civogo.Network, error) {
				return []civogo.Network{
					{ID: "network1", Name: "network1"},
					{ID: "network2", Name: "orphan-network2"},
				}, nil
			},
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return []civogo.Firewall{
					{ID: "firewall1", Name: "orphan-firewall1", ClusterCount: 0, InstanceCount: 0, LoadBalancerCount: 0},
				}, nil
			},
			FnDeleteObjectStoreCredential: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
			FnDeleteSSHKey: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
			FnDeleteNetwork: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
			FnDeleteFirewall: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		err := c.NukeOrphanedResources()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	// Test case where fetching nodes returns an error
	t.Run("error fetching nodes", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		err := c.NukeOrphanedResources()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	// Test case where fetching orphaned volumes returns an error
	t.Run("error fetching orphaned volumes", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", SSHKeyID: "key1", NetworkID: "network1"},
					},
					Pages: 1,
				}, nil
			},
			FnListSSHKeys:   func() ([]civogo.SSHKey, error) { return nil, nil },
			FnListNetworks:  func() ([]civogo.Network, error) { return nil, nil },
			FnListFirewalls: func() ([]civogo.Firewall, error) { return nil, nil },
			FnListVolumes: func() ([]civogo.Volume, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		err := c.NukeOrphanedResources()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	// Test case where deleting orphaned volumes returns an error
	t.Run("error deleting orphaned volumes", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", SSHKeyID: "key1", NetworkID: "network1"},
					},
					Pages: 1,
				}, nil
			},
			FnListVolumes: func() ([]civogo.Volume, error) {
				return []civogo.Volume{
					{ID: "volume1", Name: "orphan-volume1", Status: "available"},
				}, nil
			},
			FnListSSHKeys:   func() ([]civogo.SSHKey, error) { return nil, nil },
			FnListNetworks:  func() ([]civogo.Network, error) { return nil, nil },
			FnListFirewalls: func() ([]civogo.Firewall, error) { return nil, nil },
			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		err := c.NukeOrphanedResources()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	// Test case where fetching orphaned SSH keys returns an error
	t.Run("error fetching orphaned SSH keys", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", SSHKeyID: "key1", NetworkID: "network1"},
					},
					Pages: 1,
				}, nil
			},
			FnListVolumes:      func() ([]civogo.Volume, error) { return nil, nil },
			FnListNetworks:     func() ([]civogo.Network, error) { return nil, nil },
			FnListFirewalls:    func() ([]civogo.Firewall, error) { return nil, nil },
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) { return &civogo.PaginatedObjectstores{}, nil },
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{}, nil
			},
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		err := c.NukeOrphanedResources()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	// Test case where deleting orphaned SSH keys returns an error
	t.Run("error deleting orphaned SSH keys", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", SSHKeyID: "key1", NetworkID: "network1"},
					},
					Pages: 1,
				}, nil
			},
			FnListVolumes:      func() ([]civogo.Volume, error) { return nil, nil },
			FnListNetworks:     func() ([]civogo.Network, error) { return nil, nil },
			FnListFirewalls:    func() ([]civogo.Firewall, error) { return nil, nil },
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) { return &civogo.PaginatedObjectstores{}, nil },
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{}, nil
			},
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return []civogo.SSHKey{
					{ID: "key1", Name: "key1"},
					{ID: "key2", Name: "orphan-key2"},
				}, nil
			},
			FnDeleteSSHKey: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		err := c.NukeOrphanedResources()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	// Test case where fetching orphaned networks returns an error
	t.Run("error fetching orphaned networks", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", SSHKeyID: "key1", NetworkID: "network1"},
					},
					Pages: 1,
				}, nil
			},
			FnListVolumes:      func() ([]civogo.Volume, error) { return nil, nil },
			FnListSSHKeys:      func() ([]civogo.SSHKey, error) { return nil, nil },
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) { return &civogo.PaginatedObjectstores{}, nil },
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{}, nil
			},
			FnListNetworks: func() ([]civogo.Network, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		err := c.NukeOrphanedResources()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	// Test case where deleting orphaned networks returns an error
	t.Run("error deleting orphaned networks", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", SSHKeyID: "key1", NetworkID: "network1"},
					},
					Pages: 1,
				}, nil
			},
			FnListVolumes:      func() ([]civogo.Volume, error) { return nil, nil },
			FnListSSHKeys:      func() ([]civogo.SSHKey, error) { return nil, nil },
			FnListFirewalls:    func() ([]civogo.Firewall, error) { return nil, nil },
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) { return &civogo.PaginatedObjectstores{}, nil },
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{}, nil
			},
			FnListNetworks: func() ([]civogo.Network, error) {
				return []civogo.Network{
					{ID: "network1", Name: "network1"},
					{ID: "network2", Name: "orphan-network2"},
				}, nil
			},
			FnDeleteNetwork: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		err := c.NukeOrphanedResources()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	// Test case where fetching orphaned firewalls returns an error
	t.Run("error fetching orphaned firewalls", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", SSHKeyID: "key1", NetworkID: "network1"},
					},
					Pages: 1,
				}, nil
			},
			FnListVolumes:      func() ([]civogo.Volume, error) { return nil, nil },
			FnListSSHKeys:      func() ([]civogo.SSHKey, error) { return nil, nil },
			FnListNetworks:     func() ([]civogo.Network, error) { return nil, nil },
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) { return &civogo.PaginatedObjectstores{}, nil },
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{}, nil
			},
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
		}

		err := c.NukeOrphanedResources()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	// Test case where deleting orphaned firewalls returns an error
	t.Run("error deleting orphaned firewalls", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListInstances: func(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
				return &civogo.PaginatedInstanceList{
					Items: []civogo.Instance{
						{ID: "node1", SSHKeyID: "key1", NetworkID: "network1"},
					},
					Pages: 1,
				}, nil
			},
			FnListVolumes:      func() ([]civogo.Volume, error) { return nil, nil },
			FnListSSHKeys:      func() ([]civogo.SSHKey, error) { return nil, nil },
			FnListNetworks:     func() ([]civogo.Network, error) { return nil, nil },
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) { return &civogo.PaginatedObjectstores{}, nil },
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{}, nil
			},
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return []civogo.Firewall{
					{ID: "firewall1", Name: "orphan-firewall1", ClusterCount: 0, InstanceCount: 0, LoadBalancerCount: 0},
				}, nil
			},
			FnDeleteFirewall: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		err := c.NukeOrphanedResources()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})
}

func TestCivo_getOrphanedObjectStoreCredentials(t *testing.T) {
	t.Run("successfully fetch orphaned object store credentials", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return &civogo.PaginatedObjectstores{
					Items: []civogo.ObjectStore{
						{ID: "objectstore1", Name: "objectstore1", OwnerInfo: civogo.BucketOwner{CredentialID: "credential1"}},
						{ID: "objectstore2", Name: "objectstore2", OwnerInfo: civogo.BucketOwner{CredentialID: "credential2"}},
					},
					Pages: 1,
				}, nil
			},
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{
					Items: []civogo.ObjectStoreCredential{
						{ID: "credential1", Name: "credential1"},
						{ID: "credential2", Name: "credential2"},
						{ID: "credential3", Name: "credential3"},
						{ID: "credential4", Name: "credential4"},
					},
					Pages: 1,
				}, nil
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		got, err := c.getOrphanedObjectStoreCredentials()
		if err != nil {
			t.Fatalf("expected error to be nil, got %v", err)
		}

		expectedResponse := []civogo.ObjectStoreCredential{
			{ID: "credential3", Name: "credential3"},
			{ID: "credential4", Name: "credential4"},
		}

		if !reflect.DeepEqual(got, expectedResponse) {
			t.Errorf("expected response to be %#v, got %#v", expectedResponse, got)
		}
	})

	t.Run("error listing object stores", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		_, err := c.getOrphanedObjectStoreCredentials()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error listing object store credentials", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return &civogo.PaginatedObjectstores{
					Items: []civogo.ObjectStore{
						{ID: "objectstore1", Name: "objectstore1", OwnerInfo: civogo.BucketOwner{CredentialID: "credential1"}},
					},
					Pages: 1,
				}, nil
			},
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client: mockClient,
			logger: &mockLogger{os.Stderr},
			nuke:   true,
		}

		_, err := c.getOrphanedObjectStoreCredentials()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})
}
