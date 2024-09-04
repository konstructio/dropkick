package civo

import (
	"errors"
	"os"
	"testing"

	"github.com/civo/civogo"
)

func TestNukeKubernetesClusters(t *testing.T) {
	t.Run("successfully list and delete kubernetes clusters", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListKubernetesClusters: func() (*civogo.PaginatedKubernetesClusters, error) {
				return &civogo.PaginatedKubernetesClusters{
					Items: []civogo.KubernetesCluster{{ID: "cluster1", Name: "test-cluster1"}},
				}, nil
			},

			FnListVolumesForCluster: func(clusterID string) ([]civogo.Volume, error) {
				return []civogo.Volume{{ID: "volume1", Name: "test-volume1"}}, nil
			},

			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},

			FnDeleteKubernetesCluster: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},

			FnFindNetwork: func(search string) (*civogo.Network, error) {
				return &civogo.Network{ID: "network1", Name: "test-network1"}, nil
			},

			FnDeleteNetwork: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeKubernetesClusters()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("error listing kubernetes clusters", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListKubernetesClusters: func() (*civogo.PaginatedKubernetesClusters, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeKubernetesClusters()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error listing volumes for cluster", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListKubernetesClusters: func() (*civogo.PaginatedKubernetesClusters, error) {
				return &civogo.PaginatedKubernetesClusters{
					Items: []civogo.KubernetesCluster{{ID: "cluster1", Name: "test-cluster1"}},
				}, nil
			},

			FnListVolumesForCluster: func(clusterID string) ([]civogo.Volume, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeKubernetesClusters()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error deleting volume", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListKubernetesClusters: func() (*civogo.PaginatedKubernetesClusters, error) {
				return &civogo.PaginatedKubernetesClusters{
					Items: []civogo.KubernetesCluster{{ID: "cluster1", Name: "test-cluster1"}},
				}, nil
			},

			FnListVolumesForCluster: func(clusterID string) ([]civogo.Volume, error) {
				return []civogo.Volume{{ID: "volume1", Name: "test-volume1"}}, nil
			},

			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeKubernetesClusters()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error deleting kubernetes cluster", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListKubernetesClusters: func() (*civogo.PaginatedKubernetesClusters, error) {
				return &civogo.PaginatedKubernetesClusters{
					Items: []civogo.KubernetesCluster{{ID: "cluster1", Name: "test-cluster1"}},
				}, nil
			},

			FnListVolumesForCluster: func(clusterID string) ([]civogo.Volume, error) {
				return []civogo.Volume{{ID: "volume1", Name: "test-volume1"}}, nil
			},

			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},

			FnDeleteKubernetesCluster: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeKubernetesClusters()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error finding network", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListKubernetesClusters: func() (*civogo.PaginatedKubernetesClusters, error) {
				return &civogo.PaginatedKubernetesClusters{
					Items: []civogo.KubernetesCluster{{ID: "cluster1", Name: "test-cluster1"}},
				}, nil
			},

			FnListVolumesForCluster: func(clusterID string) ([]civogo.Volume, error) {
				return []civogo.Volume{{ID: "volume1", Name: "test-volume1"}}, nil
			},

			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},

			FnDeleteKubernetesCluster: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},

			FnFindNetwork: func(search string) (*civogo.Network, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeKubernetesClusters()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error deleting network", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListKubernetesClusters: func() (*civogo.PaginatedKubernetesClusters, error) {
				return &civogo.PaginatedKubernetesClusters{
					Items: []civogo.KubernetesCluster{{ID: "cluster1", Name: "test-cluster1"}},
				}, nil
			},

			FnListVolumesForCluster: func(clusterID string) ([]civogo.Volume, error) {
				return []civogo.Volume{{ID: "volume1", Name: "test-volume1"}}, nil
			},

			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},

			FnDeleteKubernetesCluster: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},

			FnFindNetwork: func(search string) (*civogo.Network, error) {
				return &civogo.Network{ID: "network1", Name: "test-network1"}, nil
			},

			FnDeleteNetwork: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeKubernetesClusters()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})
}
