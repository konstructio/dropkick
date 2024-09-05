package civo

import (
	"fmt"
	"io"

	"github.com/civo/civogo"
)

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
}

func (m *mockCivoClient) ListInstances(page int, perPage int) (*civogo.PaginatedInstanceList, error) {
	return m.FnListInstances(page, perPage)
}

func (m *mockCivoClient) ListSSHKeys() ([]civogo.SSHKey, error) {
	return m.FnListSSHKeys()
}

func (m *mockCivoClient) DeleteSSHKey(id string) (*civogo.SimpleResponse, error) {
	return m.FnDeleteSSHKey(id)
}

func (m *mockCivoClient) ListVolumes() ([]civogo.Volume, error) {
	return m.FnListVolumes()
}

func (m *mockCivoClient) DeleteVolume(id string) (*civogo.SimpleResponse, error) {
	return m.FnDeleteVolume(id)
}

func (m *mockCivoClient) ListKubernetesClusters() (*civogo.PaginatedKubernetesClusters, error) {
	return m.FnListKubernetesClusters()
}

func (m *mockCivoClient) DeleteKubernetesCluster(id string) (*civogo.SimpleResponse, error) {
	return m.FnDeleteKubernetesCluster(id)
}

func (m *mockCivoClient) ListVolumesForCluster(clusterID string) ([]civogo.Volume, error) {
	return m.FnListVolumesForCluster(clusterID)
}

func (m *mockCivoClient) ListNetworks() ([]civogo.Network, error) {
	return m.FnListNetworks()
}

func (m *mockCivoClient) FindNetwork(search string) (*civogo.Network, error) {
	return m.FnFindNetwork(search)
}

func (m *mockCivoClient) DeleteNetwork(id string) (*civogo.SimpleResponse, error) {
	return m.FnDeleteNetwork(id)
}

func (m *mockCivoClient) ListObjectStoreCredentials() (*civogo.PaginatedObjectStoreCredentials, error) {
	return m.FnListObjectStoreCredentials()
}

func (m *mockCivoClient) DeleteObjectStoreCredential(id string) (*civogo.SimpleResponse, error) {
	return m.FnDeleteObjectStoreCredential(id)
}

func (m *mockCivoClient) ListObjectStores() (*civogo.PaginatedObjectstores, error) {
	return m.FnListObjectStores()
}

func (m *mockCivoClient) DeleteObjectStore(id string) (*civogo.SimpleResponse, error) {
	return m.FnDeleteObjectStore(id)
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
