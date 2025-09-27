package cloudflare

import (
	"reflect"
	"testing"

	cloudflarego "github.com/cloudflare/cloudflare-go"
)

func Test_filterRecords(t *testing.T) {
	type args struct {
		records   []cloudflarego.DNSRecord
		subdomain string
	}
	tests := []struct {
		name string
		args args
		want []cloudflarego.DNSRecord
	}{
		{
			name: "successful case scenario",
			args: args{
				records: []cloudflarego.DNSRecord{
					{
						ID:   "1",
						Type: "CNAME",
					},
					{
						ID:   "2",
						Type: "TXT",
					},
				},
			},
			want: []cloudflarego.DNSRecord{
				{
					ID:   "2",
					Type: "TXT",
				},
			},
		},
		{
			name: "successful case scenario with subdomain",
			args: args{
				records: []cloudflarego.DNSRecord{
					{
						ID:   "1",
						Type: "CNAME",
						Name: "metaphor-development.ci-k1-126b3ab2-civo-gh-cf.kubesecond.com",
					},
					{
						ID:   "2",
						Type: "A",
						Name: "metaphor-development.ci-k1-2c6ab0c9-civo-gh-cf.kubesecond.com",
					},
				},
				subdomain: "ci-k1-2c6ab0c9-civo-gh-cf",
			},
			want: []cloudflarego.DNSRecord{
				{
					ID:   "2",
					Type: "A",
					Name: "metaphor-development.ci-k1-2c6ab0c9-civo-gh-cf.kubesecond.com",
				},
			},
		},
		{
			name: "successful case scenario with subdomain",
			args: args{
				records: []cloudflarego.DNSRecord{
					{
						ID:   "1",
						Type: "TXT",
						Name: "metaphor-development.ci-k1-2c6ab0c9-civo-gh-cf.kubesecond.com",
					},
					{
						ID:   "2",
						Type: "A",
						Name: "metaphor-development.ci-k1-foobar-civo-gh-cf.kubesecond.com",
					},
				},
				subdomain: "ci-k1-2c6ab0c9-civo-gh-cf",
			},
			want: []cloudflarego.DNSRecord{
				{
					ID:   "1",
					Type: "TXT",
					Name: "metaphor-development.ci-k1-2c6ab0c9-civo-gh-cf.kubesecond.com",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterRecords(tt.args.records, tt.args.subdomain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}
