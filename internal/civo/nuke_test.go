package civo

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/konstructio/dropkick/internal/civo/sdk"
	"github.com/konstructio/dropkick/internal/civo/sdk/testutils"
	"github.com/konstructio/dropkick/internal/logger"
)

func Test_NukeEverything(t *testing.T) {
	t.Run("successfully delete all resources", func(t *testing.T) {
		// generate random resources
		var (
			lbList                    = generator[sdk.LoadBalancer](rand.IntN(10) + 1)
			kubernetesList            = generator[sdk.KubernetesCluster](rand.IntN(10) + 1)
			instanceList              = generator[sdk.Instance](rand.IntN(10) + 1)
			volumeList                = generator[sdk.Volume](rand.IntN(10) + 1)
			sshKeyList                = generator[sdk.SSHKey](rand.IntN(10) + 1)
			objectStoreList           = generator[sdk.ObjectStore](rand.IntN(10) + 1)
			objectStoreCredentialList = generator[sdk.ObjectStoreCredential](rand.IntN(10) + 1)
			firewallList              = generator[sdk.Firewall](rand.IntN(10) + 1)
			networkList               = generator[sdk.Network](rand.IntN(10) + 1)
		)

		// count the number of calls to "Each"
		callCount := 0

		// there are 9 resources supported to be deleted
		// in the civo API defined by us
		numberOfResources := 9

		mock := &mockClient{
			fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
				callCount++
				switch resource.(type) {
				case sdk.LoadBalancer:
					return runEach(lbList, fn)
				case sdk.KubernetesCluster:
					return runEach(kubernetesList, fn)
				case sdk.Instance:
					return runEach(instanceList, fn)
				case sdk.Volume:
					return runEach(volumeList, fn)
				case sdk.SSHKey:
					return runEach(sshKeyList, fn)
				case sdk.ObjectStore:
					return runEach(objectStoreList, fn)
				case sdk.ObjectStoreCredential:
					return runEach(objectStoreCredentialList, fn)
				case sdk.Firewall:
					return runEach(firewallList, fn)
				case sdk.Network:
					return runEach(networkList, fn)
				default:
					t.Fatalf("unexpected resource type: %T", resource)
				}

				return nil
			},

			fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
				return nil
			},
		}

		c := &Civo{
			client: mock,
			logger: logger.None,
			nuke:   true,
		}

		err := c.NukeEverything(context.Background())
		testutils.AssertNoErrorf(t, err, "expected no error when calling NukeEverything, got %v", err)
		testutils.AssertEqualf(t, callCount, numberOfResources, "expected all resources to be deleted, got %d", callCount)
	})

	t.Run("error when deleting resources", func(t *testing.T) {
		// While I would love to shrink this test, each function uses a generic
		// parameter which cannot be passed programatically unless we also use
		// reflection. We'll keep it as-is for now.

		cases := []struct {
			name     string
			fnEach   func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error
			fnDelete func(ctx context.Context, resource sdk.APIResource) error
		}{
			{
				name: "error when deleting load balancers",
				fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
					if _, ok := resource.(sdk.LoadBalancer); ok {
						return runEach(generator[sdk.LoadBalancer](rand.IntN(10)+1), fn)
					}
					return nil
				},
				fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
					if _, ok := resource.(sdk.LoadBalancer); ok {
						return fmt.Errorf("expected error when deleting %T", resource)
					}

					return nil
				},
			},
			{
				name: "error when deleting kubernetes clusters",
				fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
					if _, ok := resource.(sdk.KubernetesCluster); ok {
						return runEach(generator[sdk.KubernetesCluster](rand.IntN(10)+1), fn)
					}
					return nil
				},
				fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
					if _, ok := resource.(sdk.KubernetesCluster); ok {
						return fmt.Errorf("expected error when deleting %T", resource)
					}

					return nil
				},
			},
			{
				name: "error when deleting instances",
				fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
					if _, ok := resource.(sdk.Instance); ok {
						return runEach(generator[sdk.Instance](rand.IntN(10)+1), fn)
					}
					return nil
				},
				fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
					if _, ok := resource.(sdk.Instance); ok {
						return fmt.Errorf("expected error when deleting %T", resource)
					}

					return nil
				},
			},
			{
				name: "error when deleting volumes",
				fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
					if _, ok := resource.(sdk.Volume); ok {
						return runEach(generator[sdk.Volume](rand.IntN(10)+1), fn)
					}
					return nil
				},
				fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
					if _, ok := resource.(sdk.Volume); ok {
						return fmt.Errorf("expected error when deleting %T", resource)
					}

					return nil
				},
			},
			{
				name: "error when deleting ssh keys",
				fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
					if _, ok := resource.(sdk.SSHKey); ok {
						return runEach(generator[sdk.SSHKey](rand.IntN(10)+1), fn)
					}
					return nil
				},
				fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
					if _, ok := resource.(sdk.SSHKey); ok {
						return fmt.Errorf("expected error when deleting %T", resource)
					}

					return nil
				},
			},
			{
				name: "error when deleting object stores",
				fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
					if _, ok := resource.(sdk.ObjectStore); ok {
						return runEach(generator[sdk.ObjectStore](rand.IntN(10)+1), fn)
					}
					return nil
				},
				fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
					if _, ok := resource.(sdk.ObjectStore); ok {
						return fmt.Errorf("expected error when deleting %T", resource)
					}

					return nil
				},
			},
			{
				name: "error when deleting object store credentials",
				fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
					if _, ok := resource.(sdk.ObjectStoreCredential); ok {
						return runEach(generator[sdk.ObjectStoreCredential](rand.IntN(10)+1), fn)
					}
					return nil
				},
				fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
					if _, ok := resource.(sdk.ObjectStoreCredential); ok {
						return fmt.Errorf("expected error when deleting %T", resource)
					}

					return nil
				},
			},
			{
				name: "error when deleting firewalls",
				fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
					if _, ok := resource.(sdk.Firewall); ok {
						return runEach(generator[sdk.Firewall](rand.IntN(10)+1), fn)
					}
					return nil
				},
				fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
					if _, ok := resource.(sdk.Firewall); ok {
						return fmt.Errorf("expected error when deleting %T", resource)
					}

					return nil
				},
			},
			{
				name: "error when deleting networks",
				fnEach: func(ctx context.Context, resource sdk.APIResource, fn func(sdk.APIResource) error) error {
					if _, ok := resource.(sdk.Network); ok {
						return runEach(generator[sdk.Network](rand.IntN(10)+1), fn)
					}
					return nil
				},
				fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
					if _, ok := resource.(sdk.Network); ok {
						return fmt.Errorf("expected error when deleting %T", resource)
					}

					return nil
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				mock := &mockClient{
					fnEach:   tc.fnEach,
					fnDelete: tc.fnDelete,
				}

				c := &Civo{
					client: mock,
					logger: logger.None,
					nuke:   true,
				}

				err := c.NukeEverything(context.Background())
				testutils.AssertErrorf(t, err, "expected error when calling NukeEverything, got %v", err)
			})
		}
	})
}
