package civo

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/civo/civogo"
	"github.com/konstructio/dropkick/internal/logger"
)

type errNotImplemented struct {
	funcName string
}

func (e *errNotImplemented) Error() string {
	return fmt.Sprintf("%q is not implemented", e.funcName)
}

type mockCivoClient struct {
	FnListInstances               func(page int, perPage int) (*civogo.PaginatedInstanceList, error)
	FnListSSHKeys                 func() ([]civogo.SSHKey, error)
	FnDeleteSSHKey                func(id string) (*civogo.SimpleResponse, error)
	FnListVolumes                 func() ([]civogo.Volume, error)
	FnDeleteVolume                func(id string) (*civogo.SimpleResponse, error)
	FnListKubernetesClusters      func() (*civogo.PaginatedKubernetesClusters, error)
	FnDeleteKubernetesCluster     func(id string) (*civogo.SimpleResponse, error)
	FnListVolumesForCluster       func(clusterID string) ([]civogo.Volume, error)
	FnListNetworks                func() ([]civogo.Network, error)
	FnFindNetwork                 func(search string) (*civogo.Network, error)
	FnDeleteNetwork               func(id string) (*civogo.SimpleResponse, error)
	FnListObjectStoreCredentials  func() (*civogo.PaginatedObjectStoreCredentials, error)
	FnGetObjectStoreCredential    func(id string) (*civogo.ObjectStoreCredential, error)
	FnDeleteObjectStoreCredential func(id string) (*civogo.SimpleResponse, error)
	FnListObjectStores            func() (*civogo.PaginatedObjectstores, error)
	FnDeleteObjectStore           func(id string) (*civogo.SimpleResponse, error)
	FnListFirewalls               func() ([]civogo.Firewall, error)
	FnDeleteFirewall              func(id string) (*civogo.SimpleResponse, error)
}

func (m *mockCivoClient) ListInstances(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
	if m.FnListInstances == nil {
		return nil, &errNotImplemented{funcName: "ListInstances"}
	}

	return m.FnListInstances(page, perPage)
}

func (m *mockCivoClient) ListSSHKeys() ([]civogo.SSHKey, error) {
	if m.FnListSSHKeys == nil {
		return nil, &errNotImplemented{funcName: "ListSSHKeys"}
	}

	return m.FnListSSHKeys()
}

func (m *mockCivoClient) DeleteSSHKey(id string) (*civogo.SimpleResponse, error) {
	if m.FnDeleteSSHKey == nil {
		return nil, &errNotImplemented{funcName: "DeleteSSHKey"}
	}

	return m.FnDeleteSSHKey(id)
}

func (m *mockCivoClient) ListVolumes() ([]civogo.Volume, error) {
	if m.FnListVolumes == nil {
		return nil, &errNotImplemented{funcName: "ListVolumes"}
	}

	return m.FnListVolumes()
}

func (m *mockCivoClient) DeleteVolume(id string) (*civogo.SimpleResponse, error) {
	if m.FnDeleteVolume == nil {
		return nil, &errNotImplemented{funcName: "DeleteVolume"}
	}

	return m.FnDeleteVolume(id)
}

func (m *mockCivoClient) ListKubernetesClusters() (*civogo.PaginatedKubernetesClusters, error) {
	if m.FnListKubernetesClusters == nil {
		return nil, &errNotImplemented{funcName: "ListKubernetesClusters"}
	}

	return m.FnListKubernetesClusters()
}

func (m *mockCivoClient) DeleteKubernetesCluster(id string) (*civogo.SimpleResponse, error) {
	if m.FnDeleteKubernetesCluster == nil {
		return nil, &errNotImplemented{funcName: "DeleteKubernetesCluster"}
	}

	return m.FnDeleteKubernetesCluster(id)
}

func (m *mockCivoClient) ListVolumesForCluster(clusterID string) ([]civogo.Volume, error) {
	if m.FnListVolumesForCluster == nil {
		return nil, &errNotImplemented{funcName: "ListVolumesForCluster"}
	}

	return m.FnListVolumesForCluster(clusterID)
}

func (m *mockCivoClient) ListNetworks() ([]civogo.Network, error) {
	if m.FnListNetworks == nil {
		return nil, &errNotImplemented{funcName: "ListNetworks"}
	}

	return m.FnListNetworks()
}

func (m *mockCivoClient) FindNetwork(search string) (*civogo.Network, error) {
	if m.FnFindNetwork == nil {
		return nil, &errNotImplemented{funcName: "FindNetwork"}
	}

	return m.FnFindNetwork(search)
}

func (m *mockCivoClient) DeleteNetwork(id string) (*civogo.SimpleResponse, error) {
	if m.FnDeleteNetwork == nil {
		return nil, &errNotImplemented{funcName: "DeleteNetwork"}
	}

	return m.FnDeleteNetwork(id)
}

func (m *mockCivoClient) ListObjectStoreCredentials() (*civogo.PaginatedObjectStoreCredentials, error) {
	if m.FnListObjectStoreCredentials == nil {
		return nil, &errNotImplemented{funcName: "ListObjectStoreCredentials"}
	}

	return m.FnListObjectStoreCredentials()
}

func (m *mockCivoClient) GetObjectStoreCredential(id string) (*civogo.ObjectStoreCredential, error) {
	if m.FnGetObjectStoreCredential == nil {
		return nil, &errNotImplemented{funcName: "GetObjectStoreCredential"}
	}

	return m.FnGetObjectStoreCredential(id)
}

func (m *mockCivoClient) DeleteObjectStoreCredential(id string) (*civogo.SimpleResponse, error) {
	if m.FnDeleteObjectStoreCredential == nil {
		return nil, &errNotImplemented{funcName: "DeleteObjectStoreCredential"}
	}

	return m.FnDeleteObjectStoreCredential(id)
}

func (m *mockCivoClient) ListObjectStores() (*civogo.PaginatedObjectstores, error) {
	if m.FnListObjectStores == nil {
		return nil, &errNotImplemented{funcName: "ListObjectStores"}
	}

	return m.FnListObjectStores()
}

func (m *mockCivoClient) DeleteObjectStore(id string) (*civogo.SimpleResponse, error) {
	if m.FnDeleteObjectStore == nil {
		return nil, &errNotImplemented{funcName: "DeleteObjectStore"}
	}

	return m.FnDeleteObjectStore(id)
}

func (m *mockCivoClient) ListFirewalls() ([]civogo.Firewall, error) {
	if m.FnListFirewalls == nil {
		return nil, &errNotImplemented{funcName: "ListFirewalls"}
	}

	return m.FnListFirewalls()
}

func (m *mockCivoClient) DeleteFirewall(id string) (*civogo.SimpleResponse, error) {
	if m.FnDeleteFirewall == nil {
		return nil, &errNotImplemented{funcName: "DeleteFirewall"}
	}

	return m.FnDeleteFirewall(id)
}

type mockLogger struct {
	output io.Writer
}

func (m *mockLogger) Infof(format string, args ...interface{}) {
	fmt.Fprintf(m.output, "[INFO] "+format+"\n", args...)
}

func (m *mockLogger) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(m.output, "[ERROR] "+format+"\n", args...)
}

func (m *mockLogger) Warnf(format string, args ...interface{}) {
	fmt.Fprintf(m.output, "[WARN] "+format+"\n", args...)
}

func TestNew(t *testing.T) {
	cases := []struct {
		Name       string
		Opts       []Option
		Token      string
		Region     string
		APIURL     string
		Context    context.Context
		Logger     *logger.Logger
		Nuke       bool
		NameFilter string
		WantErr    bool
	}{
		{
			Name:       "all good",
			Token:      "token",
			Region:     "region",
			APIURL:     "https://api.example.com",
			Context:    context.Background(),
			Logger:     logger.None,
			Nuke:       true,
			NameFilter: "filter",
			WantErr:    false,
		},
		{
			Name: "errored out option",
			Opts: []Option{
				func(c *Civo) error {
					return fmt.Errorf("error")
				},
			},
			WantErr: true,
		},
		{
			Name:    "missing token",
			Token:   "",
			WantErr: true,
		},
		{
			Name:    "missing region",
			Token:   "token",
			Region:  "",
			WantErr: true,
		},
		{
			Name:    "impossible client using invalid URL",
			Token:   "token",
			Region:  "region",
			APIURL:  "#@$%^&*",
			WantErr: true,
		},
		{
			Name:    "missing context",
			Token:   "token",
			Region:  "region",
			Context: nil,
			Logger:  nil,
		},
		{
			Name:    "default logger",
			Token:   "token",
			Region:  "region",
			Logger:  nil,
			WantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(tt *testing.T) {
			opts := append([]Option{}, tc.Opts...)

			if tc.Token != "" {
				opts = append(opts, WithToken(tc.Token))
			}

			if tc.Region != "" {
				opts = append(opts, WithRegion(tc.Region))
			}

			if tc.APIURL != "" {
				opts = append(opts, WithAPIURL(tc.APIURL))
			}

			if tc.Context != nil {
				opts = append(opts, WithContext(tc.Context))
			}

			if tc.Logger != nil {
				opts = append(opts, WithLogger(tc.Logger))
			}

			if tc.NameFilter != "" {
				opts = append(opts, WithNameFilter(tc.NameFilter))
			}

			opts = append(opts, WithNuke(tc.Nuke))

			client, err := New(opts...)

			if tc.WantErr {
				if err == nil {
					tt.Fatal("expected err to not be nil")
				}

				return
			}

			if err != nil {
				tt.Fatalf("expected err to be nil, got %v", err)
			}

			if client == nil {
				tt.Fatal("expected client to not be nil")
			}

			if client.client == nil {
				tt.Fatal("expected client.client to not be nil")
			}

			if client.token != tc.Token {
				tt.Fatalf("expected client.token to be %q, got %q", tc.Token, client.token)
			}

			if client.region != tc.Region {
				tt.Fatalf("expected client.region to be %q, got %q", tc.Region, client.region)
			}

			if tc.APIURL == "" && client.apiURL != civoAPIURL {
				tt.Fatalf("expected client.apiURL to be %q, got %q", civoAPIURL, client.apiURL)
			}

			if tc.APIURL != "" && client.apiURL != tc.APIURL {
				tt.Fatalf("expected client.apiURL to be %q, got %q", tc.APIURL, client.apiURL)
			}

			if tc.Context == nil && client.context != context.Background() {
				tt.Fatalf("expected client.context to be %v, got %v", context.Background(), client.context)
			}

			if tc.Context != nil && client.context != tc.Context {
				tt.Fatalf("expected client.context to be %v, got %v", tc.Context, client.context)
			}

			t.Logf("client.logger: %#v", client.logger)
			t.Logf("tc.Logger: %#v", tc.Logger)

			if tc.Logger == nil && client.logger != logger.None {
				tt.Fatalf("expected client.logger to be the default logger, got: %#v", client.logger)
			}

			if tc.Logger != nil && client.logger != tc.Logger {
				tt.Fatalf("expected client.logger to be %#v, got %#v", tc.Logger, client.logger)
			}

			if client.nuke != tc.Nuke {
				tt.Fatalf("expected client.nuke to be %v, got %v", tc.Nuke, client.nuke)
			}

			if client.nameFilter != tc.NameFilter {
				tt.Fatalf("expected client.nameFilter to be %q, got %q", tc.NameFilter, client.nameFilter)
			}
		})
	}
}
