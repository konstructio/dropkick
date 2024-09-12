package civo

import (
	"fmt"
	"io"

	"github.com/civo/civogo"
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
