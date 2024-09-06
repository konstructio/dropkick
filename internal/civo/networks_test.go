package civo

import (
	"errors"
	"os"
	"testing"

	"github.com/civo/civogo"
)

func TestNukeNetworks(t *testing.T) {
	t.Run("successfully list and delete networks", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListNetworks: func() ([]civogo.Network, error) {
				return []civogo.Network{
					{ID: "network1", Name: "test-network1"},
				}, nil
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

		err := c.NukeNetworks()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("error listing networks", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListNetworks: func() ([]civogo.Network, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeNetworks()
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
			FnListNetworks: func() ([]civogo.Network, error) {
				return []civogo.Network{
					{ID: "network1", Name: "test-network1"},
				}, nil
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

		err := c.NukeNetworks()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("skipping network due to name filter", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListNetworks: func() ([]civogo.Network, error) {
				return []civogo.Network{
					{ID: "network1", Name: "test-network1"},
				}, nil
			},
			FnDeleteNetwork: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "non-matching-filter",
			nuke:       true,
		}

		err := c.NukeNetworks()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("nuke is not enabled, refuse to delete", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListNetworks: func() ([]civogo.Network, error) {
				return []civogo.Network{
					{ID: "network1", Name: "test-network1"},
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       false,
		}

		err := c.NukeNetworks()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})
}
