package sdk

import (
	"strings"
	"testing"
)

func Test_IsPaginatedResource(t *testing.T) {
	tests := []struct {
		endpoint  string
		paginated bool
		wantErr   bool
	}{
		{
			endpoint:  Instance{}.GetAPIEndpoint(),
			paginated: true,
			wantErr:   false,
		},
		{
			endpoint:  Firewall{}.GetAPIEndpoint(),
			paginated: false,
			wantErr:   false,
		},
		{
			endpoint:  Volume{}.GetAPIEndpoint(),
			paginated: false,
			wantErr:   false,
		},
		{
			endpoint:  KubernetesCluster{}.GetAPIEndpoint(),
			paginated: true,
			wantErr:   false,
		},
		{
			endpoint:  Network{}.GetAPIEndpoint(),
			paginated: false,
			wantErr:   false,
		},
		{
			endpoint:  ObjectStore{}.GetAPIEndpoint(),
			paginated: true,
			wantErr:   false,
		},
		{
			endpoint:  ObjectStoreCredential{}.GetAPIEndpoint(),
			paginated: true,
			wantErr:   false,
		},
		{
			endpoint:  SSHKey{}.GetAPIEndpoint(),
			paginated: false,
			wantErr:   false,
		},
		{
			endpoint:  LoadBalancer{}.GetAPIEndpoint(),
			paginated: false,
			wantErr:   false,
		},
		{
			endpoint:  "invalid",
			paginated: false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(strings.TrimPrefix(tt.endpoint, "/v2/"), func(t *testing.T) {
			got, err := IsPaginatedResource(tt.endpoint)
			if (err != nil) != tt.wantErr {
				t.Fatalf("not expecting an error, got %v", err)
			}

			if got != tt.paginated {
				t.Fatalf("expected %s to be paginated=%v, got paginated=%v", tt.endpoint, tt.paginated, got)
			}
		})
	}
}
