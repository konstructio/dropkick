package civo

import (
	"errors"
	"os"
	"testing"

	"github.com/civo/civogo"
)

func TestNukeVolumes(t *testing.T) {
	t.Run("successfully list and delete unattached volumes", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListVolumes: func() ([]civogo.Volume, error) {
				return []civogo.Volume{
					{ID: "volume1", Name: "test-volume1", Status: "available"},
					{ID: "volume2", Name: "test-volume2", Status: "available"},
				}, nil
			},
			FnDeleteVolume: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeVolumes()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
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
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeVolumes()
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
			FnListVolumes: func() ([]civogo.Volume, error) {
				return []civogo.Volume{
					{ID: "volume1", Name: "test-volume1", Status: "available"},
				}, nil
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

		err := c.NukeVolumes()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("skipping volume due to name filter", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListVolumes: func() ([]civogo.Volume, error) {
				return []civogo.Volume{
					{ID: "volume1", Name: "test-volume1", Status: "available"},
					{ID: "volume2", Name: "test-volume2", Status: "available"},
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "non-matching-filter",
			nuke:       true,
		}

		err := c.NukeVolumes()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("refuse to delete attached volume", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListVolumes: func() ([]civogo.Volume, error) {
				return []civogo.Volume{
					{ID: "volume1", Name: "attached-volume", Status: volumeStatusAttached, InstanceID: "instance1"},
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeVolumes()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("nuke is not enabled, refuse to delete volumes", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListVolumes: func() ([]civogo.Volume, error) {
				return []civogo.Volume{
					{ID: "volume1", Name: "test-volume1", Status: "available"},
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       false,
		}

		err := c.NukeVolumes()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})
}
