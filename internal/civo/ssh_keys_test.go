package civo

import (
	"errors"
	"os"
	"testing"

	"github.com/civo/civogo"
)

func TestNukeSSHKeys(t *testing.T) {
	t.Run("successfully list and delete SSH keys", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return []civogo.SSHKey{
					{ID: "key1", Name: "test-key1"},
					{ID: "key2", Name: "test-key2"},
				}, nil
			},
			FnDeleteSSHKey: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeSSHKeys()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("error listing SSH keys", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeSSHKeys()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error deleting SSH key", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return []civogo.SSHKey{
					{ID: "key1", Name: "test-key1"},
				}, nil
			},
			FnDeleteSSHKey: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeSSHKeys()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("skipping SSH key due to name filter", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return []civogo.SSHKey{
					{ID: "key1", Name: "test-key1"},
					{ID: "key2", Name: "test-key2"},
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "non-matching-filter",
			nuke:       true,
		}

		err := c.NukeSSHKeys()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("nuke is not enabled, refuse to delete SSH key", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListSSHKeys: func() ([]civogo.SSHKey, error) {
				return []civogo.SSHKey{
					{ID: "key1", Name: "test-key1"},
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       false,
		}

		err := c.NukeSSHKeys()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})
}
