package civo

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/konstructio/dropkick/internal/civo/sdk"
	"github.com/konstructio/dropkick/internal/logger"
)

// mockClient is a mock implementation of the Client interface.
type mockClient struct {
	fnGetInstances              func(ctx context.Context) ([]sdk.Instance, error)
	fnGetFirewalls              func(ctx context.Context) ([]sdk.Firewall, error)
	fnGetVolumes                func(ctx context.Context) ([]sdk.Volume, error)
	fnGetKubernetesClusters     func(ctx context.Context) ([]sdk.KubernetesCluster, error)
	fnGetNetworks               func(ctx context.Context) ([]sdk.Network, error)
	fnGetObjectStores           func(ctx context.Context) ([]sdk.ObjectStore, error)
	fnGetObjectStoreCredentials func(ctx context.Context) ([]sdk.ObjectStoreCredential, error)
	fnGetLoadBalancers          func(ctx context.Context) ([]sdk.LoadBalancer, error)
	fnGetSSHKeys                func(ctx context.Context) ([]sdk.SSHKey, error)
	fnDelete                    func(ctx context.Context, resource sdk.APIResource) error
	fnEach                      func(ctx context.Context, v sdk.APIResource, iterator func(sdk.APIResource) error) error
}

// Ensure mockClient implements the Client interface.
var _ Client = &mockClient{}

func (m *mockClient) GetInstances(ctx context.Context) ([]sdk.Instance, error) {
	return m.fnGetInstances(ctx)
}

func (m *mockClient) GetFirewalls(ctx context.Context) ([]sdk.Firewall, error) {
	return m.fnGetFirewalls(ctx)
}

func (m *mockClient) GetVolumes(ctx context.Context) ([]sdk.Volume, error) {
	return m.fnGetVolumes(ctx)
}

func (m *mockClient) GetKubernetesClusters(ctx context.Context) ([]sdk.KubernetesCluster, error) {
	return m.fnGetKubernetesClusters(ctx)
}

func (m *mockClient) GetNetworks(ctx context.Context) ([]sdk.Network, error) {
	return m.fnGetNetworks(ctx)
}

func (m *mockClient) GetObjectStores(ctx context.Context) ([]sdk.ObjectStore, error) {
	return m.fnGetObjectStores(ctx)
}

func (m *mockClient) GetObjectStoreCredentials(ctx context.Context) ([]sdk.ObjectStoreCredential, error) {
	return m.fnGetObjectStoreCredentials(ctx)
}

func (m *mockClient) GetLoadBalancers(ctx context.Context) ([]sdk.LoadBalancer, error) {
	return m.fnGetLoadBalancers(ctx)
}

func (m *mockClient) GetSSHKeys(ctx context.Context) ([]sdk.SSHKey, error) {
	return m.fnGetSSHKeys(ctx)
}

func (m *mockClient) Delete(ctx context.Context, resource sdk.APIResource) error {
	return m.fnDelete(ctx, resource)
}

func (m *mockClient) Each(ctx context.Context, v sdk.APIResource, iterator func(sdk.APIResource) error) error {
	return m.fnEach(ctx, v, iterator)
}

// runEach runs the given function for each item in the list.
func runEach[T sdk.Resource](list []T, fn func(sdk.APIResource) error) error {
	for _, item := range list {
		if err := fn(item); err != nil {
			return err
		}
	}
	return nil
}

// generator generates a list of resources of type sdk.Resource, which are
// Civo-specific resources.
func generator[T sdk.Resource](n int) []T {
	list := make([]T, n)

	for i := 0; i < n; i++ {
		resType := reflect.TypeOf((*T)(nil)).Elem()
		resValue := reflect.New(resType).Elem()

		// lowercase type name
		typeName := strings.ToLower(resType.Name())

		// set the ID field
		if idField := resValue.FieldByName("ID"); idField.IsValid() && idField.CanSet() {
			idField.SetString("id-" + strconv.Itoa(i+1))
		}

		// set the Name field
		if nameField := resValue.FieldByName("Name"); nameField.IsValid() && nameField.CanSet() {
			nameField.SetString(typeName + "-" + strconv.Itoa(i+1))
		}

		// for firewalls, set the label as the name
		if labelField := resValue.FieldByName("Label"); labelField.IsValid() && labelField.CanSet() {
			labelField.SetString(typeName + "-" + strconv.Itoa(i+1))
		}

		list[i] = resValue.Interface().(T)
	}

	return list
}

func TestNew(t *testing.T) {
	cases := []struct {
		Name       string
		Opts       []Option
		Token      string
		Region     string
		Logger     *logger.Logger
		Nuke       bool
		NameFilter string
		WantErr    bool
	}{
		{
			Name:       "all good",
			Token:      "token",
			Region:     "region",
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
			Name:   "missing context",
			Token:  "token",
			Region: "region",
			Logger: nil,
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
