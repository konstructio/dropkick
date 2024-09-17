package civo

import (
	"errors"
	"os"
	"testing"

	"github.com/civo/civogo"
)

func TestNukeObjectStores(t *testing.T) {
	t.Run("successfully list and delete object stores", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return &civogo.PaginatedObjectstores{
					Items: []civogo.ObjectStore{{ID: "objStore1", Name: "test-store1"}},
					Pages: 1,
				}, nil
			},
			FnDeleteObjectStore: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
			FnGetObjectStoreCredential: func(id string) (*civogo.ObjectStoreCredential, error) {
				return &civogo.ObjectStoreCredential{ID: "cred1", Name: "test-cred1"}, nil
			},
			FnDeleteObjectStoreCredential: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeObjectStores()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
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
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeObjectStores()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error deleting object store", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return &civogo.PaginatedObjectstores{
					Items: []civogo.ObjectStore{{ID: "objStore1", Name: "test-store1"}},
					Pages: 1,
				}, nil
			},
			FnGetObjectStoreCredential: func(id string) (*civogo.ObjectStoreCredential, error) {
				return nil, civogo.ZeroMatchesError
			},
			FnDeleteObjectStore: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeObjectStores()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("skipping object store due to name filter", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return &civogo.PaginatedObjectstores{
					Items: []civogo.ObjectStore{{ID: "objStore1", Name: "test-store1"}},
					Pages: 1,
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "non-matching-filter",
			nuke:       true,
		}

		err := c.NukeObjectStores()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("nuke is not enabled, refuse to delete object store", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return &civogo.PaginatedObjectstores{
					Items: []civogo.ObjectStore{{ID: "objStore1", Name: "test-store1"}},
					Pages: 1,
				}, nil
			},
			FnGetObjectStoreCredential: func(id string) (*civogo.ObjectStoreCredential, error) {
				return nil, civogo.ZeroMatchesError
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       false,
		}

		err := c.NukeObjectStores()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("delete object store but no credential found", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return &civogo.PaginatedObjectstores{
					Items: []civogo.ObjectStore{{ID: "objStore1", Name: "test-store1"}},
					Pages: 1,
				}, nil
			},
			FnGetObjectStoreCredential: func(id string) (*civogo.ObjectStoreCredential, error) {
				return nil, civogo.ZeroMatchesError
			},
			FnDeleteObjectStore: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeObjectStores()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("error finding object store credential", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListObjectStores: func() (*civogo.PaginatedObjectstores, error) {
				return &civogo.PaginatedObjectstores{
					Items: []civogo.ObjectStore{{ID: "objStore1", Name: "test-store1"}},
					Pages: 1,
				}, nil
			},
			FnGetObjectStoreCredential: func(id string) (*civogo.ObjectStoreCredential, error) {
				return nil, errTest
			},
			FnDeleteObjectStore: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
			FnDeleteObjectStoreCredential: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeObjectStores()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})
}

func TestNukeObjectStoreCredentials(t *testing.T) {
	t.Run("successfully list and delete object store credentials", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{
					Items: []civogo.ObjectStoreCredential{{ID: "cred1", Name: "test-cred1"}},
					Pages: 1,
				}, nil
			},
			FnDeleteObjectStoreCredential: func(id string) (*civogo.SimpleResponse, error) {
				return &civogo.SimpleResponse{ErrorCode: "200"}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeObjectStoreCredentials()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("error listing object store credentials", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeObjectStoreCredentials()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("error deleting object store credential", func(tt *testing.T) {
		errTest := errors.New("test error")

		mockClient := &mockCivoClient{
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{
					Items: []civogo.ObjectStoreCredential{{ID: "cred1", Name: "test-cred1"}},
					Pages: 1,
				}, nil
			},
			FnDeleteObjectStoreCredential: func(id string) (*civogo.SimpleResponse, error) {
				return nil, errTest
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       true,
		}

		err := c.NukeObjectStoreCredentials()
		if err == nil {
			tt.Errorf("expected error to be %v, got nil", errTest)
		}

		if !errors.Is(err, errTest) {
			tt.Errorf("expected error to be %v, got %v", errTest, err)
		}
	})

	t.Run("skipping object store credential due to name filter", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{
					Items: []civogo.ObjectStoreCredential{{ID: "cred1", Name: "test-cred1"}},
					Pages: 1,
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "non-matching-filter",
			nuke:       true,
		}

		err := c.NukeObjectStoreCredentials()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})

	t.Run("nuke is not enabled, refuse to delete object store credential", func(tt *testing.T) {
		mockClient := &mockCivoClient{
			FnListObjectStoreCredentials: func() (*civogo.PaginatedObjectStoreCredentials, error) {
				return &civogo.PaginatedObjectStoreCredentials{
					Items: []civogo.ObjectStoreCredential{{ID: "cred1", Name: "test-cred1"}},
					Pages: 1,
				}, nil
			},
		}

		c := &Civo{
			client:     mockClient,
			logger:     &mockLogger{os.Stderr},
			nameFilter: "",
			nuke:       false,
		}

		err := c.NukeObjectStoreCredentials()
		if err != nil {
			tt.Errorf("expected error to be nil, got %v", err)
		}
	})
}
