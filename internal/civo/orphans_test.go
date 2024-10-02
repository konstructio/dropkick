package civo

import (
	"context"
	"net/http"
	"testing"

	"github.com/konstructio/dropkick/internal/civo/sdk"
	"github.com/konstructio/dropkick/internal/civo/sdk/testutils"
	"github.com/konstructio/dropkick/internal/logger"
)

func Test_NukeOrphanedResources(t *testing.T) {
	t.Run("find at least one orphan per resource and delete it without name filter", func(t *testing.T) {
		const token = "abc123"
		const region = "nyc1"

		instances := []sdk.Instance{
			{ID: "1", Name: "test-instance-1", SSHKeyID: "1", NetworkID: "1", FirewallID: "1"},
			{ID: "2", Name: "test-instance-2", SSHKeyID: "2", NetworkID: "3", FirewallID: "1"},
			{ID: "3", Name: "test-instance-3", SSHKeyID: "2", NetworkID: "1", FirewallID: "1"},
			{ID: "4", Name: "test-instance-4"},
			{ID: "5", Name: "test-instance-5"},
		}

		sshkeys := []sdk.SSHKey{
			{ID: "1", Name: "test-sshkey-1"}, // in use by instance
			{ID: "2", Name: "test-sshkey-2"}, // in use by instance
			{ID: "3", Name: "test-sshkey-3"}, // orphan
			{ID: "4", Name: "test-sshkey-4"}, // orphan
		}

		networks := []sdk.Network{
			{ID: "1", Name: "test-network-1"}, // in use by instance
			{ID: "2", Name: "test-network-2"}, // in use by volume
			{ID: "3", Name: "test-network-3"}, // in use by instance
			{ID: "4", Name: "test-network-4"}, // in use by volume
			{ID: "5", Name: "test-network-5"}, // orphan
			{ID: "6", Name: "test-network-6"}, // in use by volume
		}

		firewalls := []sdk.Firewall{
			{ID: "1", Name: "test-firewall-1"}, // in use by instance
			{ID: "2", Name: "test-firewall-2"}, // orphan
			{ID: "3", Name: "test-firewall-3"}, // orphan
			{ID: "4", Name: "test-firewall-4"}, // orphan
		}

		volumes := []sdk.Volume{
			{ID: "1", Name: "test-volume-1", Status: "attached", InstanceID: "1", NetworkID: "2"}, // in use by instance
			{ID: "2", Name: "test-volume-2", Status: "attached", InstanceID: "1", NetworkID: "4"}, // in use by instance
			{ID: "3", Name: "test-volume-3"}, // orphan
			{ID: "4", Name: "test-volume-4"}, // orphan
			{ID: "5", Name: "test-volume-5", Status: "attached", InstanceID: "3", NetworkID: "6"}, // in use by instance
			{ID: "6", Name: "test-volume-6"}, // orphan
		}

		lbs := []sdk.LoadBalancer{
			{ID: "1", Name: "test-lb-1", ClusterID: "1"},  // in use by k8s cluster
			{ID: "2", Name: "test-lb-2", ClusterID: "1"},  // in use by k8s cluster
			{ID: "3", Name: "test-lb-3"},                  // orphan
			{ID: "4", Name: "test-lb-4", FirewallID: "1"}, // in use by firewall
			{ID: "5", Name: "test-lb-5", FirewallID: "1"}, // in use by firewall
			{ID: "6", Name: "test-lb-6"},
			{ID: "7", Name: "test-lb-7"},
		}

		objstorecreds := []sdk.ObjectStoreCredential{
			{ID: "1", Name: "test-objstore-cred-1"}, // in use by objstore 1
			{ID: "2", Name: "test-objstore-cred-2"}, // in use by objstore 2
			{ID: "3", Name: "test-objstore-cred-3"}, // orphan
			{ID: "4", Name: "test-objstore-cred-4"}, // orphan
			{ID: "5", Name: "test-objstore-cred-5"}, // in use by objstore 5
			{ID: "6", Name: "test-objstore-cred-6"}, // orphan
		}

		objstores := []sdk.ObjectStore{
			{ID: "1", Name: "test-objstore-1", Credentials: objstorecreds[1]}, // in use by objstore
			{ID: "2", Name: "test-objstore-2", Credentials: objstorecreds[2]}, // in use by objstore
			{ID: "3", Name: "test-objstore-3"},
			{ID: "4", Name: "test-objstore-4"},
			{ID: "5", Name: "test-objstore-5", Credentials: objstorecreds[5]}, // in use by objstore
			{ID: "6", Name: "test-objstore-6"},
		}

		client := testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				t.Logf("%s %q -> params: %#v", method, location, params)

				// Check region in params
				if params["region"] != region {
					t.Fatalf("unexpected region: %q", params["region"])
				}

				// On resource deletions, we return no errors
				if method == http.MethodDelete {
					return nil
				}

				// Ensure we only send GET requests
				if method != http.MethodGet {
					t.Fatalf("unexpected method: %q", method)
				}

				var page, perPage int

				paginated, err := sdk.IsPaginatedResource(location)
				testutils.AssertNoErrorf(t, err, "resource %q doesn't exist", location)

				if paginated {
					// ensure paginated resources are handled correctly
					page = testutils.ParseParamNumber(t, location, params, "page", 1, true)
					perPage = testutils.ParseParamNumber(t, location, params, "per_page", 100, true)
				} else {
					// sanity check: ensure non-paginated resources don't
					// attempt to send page or per_page parameters
					if _, ok := params["page"]; ok {
						t.Fatal("unexpected page parameter")
					}
					if _, ok := params["per_page"]; ok {
						t.Fatal("unexpected per_page parameter")
					}
				}

				switch out := output.(type) {
				case *sdk.PaginatedResponse[sdk.Instance]:
					res, page, perPage, totalPages := testutils.GetResultsForPage(t, instances, page, perPage)
					out.Items = res
					out.Page = page
					out.PerPage = perPage
					out.Pages = totalPages

				case *[]sdk.Firewall:
					testutils.InjectIntoSlice(t, out, firewalls)

				case *[]sdk.Volume:
					testutils.InjectIntoSlice(t, out, volumes)

				case *[]sdk.Network:
					testutils.InjectIntoSlice(t, out, networks)

				case *[]sdk.SSHKey:
					testutils.InjectIntoSlice(t, out, sshkeys)

				case *[]sdk.LoadBalancer:
					testutils.InjectIntoSlice(t, out, lbs)

				case *sdk.PaginatedResponse[sdk.ObjectStore]:
					res, page, perPage, totalPages := testutils.GetResultsForPage(t, objstores, page, perPage)
					out.Items = res
					out.Page = page
					out.PerPage = perPage
					out.Pages = totalPages

				case *sdk.PaginatedResponse[sdk.ObjectStoreCredential]:
					res, page, perPage, totalPages := testutils.GetResultsForPage(t, objstorecreds, page, perPage)
					out.Items = res
					out.Page = page
					out.PerPage = perPage
					out.Pages = totalPages

				case nil:
					// no-op

				default:
					t.Fatalf("unexpected output type: %T", output)
				}

				return nil
			},

			FnGetRegion: func() string {
				return region
			},
		}

		civo := &Civo{
			client: &client,
			nuke:   true,
			region: region,
			token:  token,
			logger: logger.None,
		}

		err := civo.NukeOrphanedResources(context.Background())
		testutils.AssertNoError(t, err)
	})
}
