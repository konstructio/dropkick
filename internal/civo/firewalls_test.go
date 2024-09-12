package civo

import (
	"errors"
	"os"
	"testing"

	"github.com/civo/civogo"
)

func TestNukeFirewalls(t *testing.T) {
	t.Run("successfully list and delete firewalls", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return []civogo.Firewall{
					{ID: "firewall1", Name: "test-firewall1"},
				}, nil
			},
			FnDeleteFirewall: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeFirewalls()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
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
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeFirewalls()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error deleting firewall", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return []civogo.Firewall{
					{ID: "firewall1", Name: "test-firewall1"},
				}, nil
			},
			FnDeleteFirewall: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeFirewalls()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("skipping firewall due to name filter", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return []civogo.Firewall{
					{ID: "firewall1", Name: "test-firewall1"},
				}, nil
			},
			FnDeleteFirewall: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "non-matching-filter",
			nuke:       true,
		}

		err := c.NukeFirewalls()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("nuke is not enabled, refuse to delete", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListFirewalls: func() ([]civogo.Firewall, error) {
				return []civogo.Firewall{
					{ID: "firewall1", Name: "test-firewall1"},
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       false,
		}

		err := c.NukeFirewalls()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})
}
