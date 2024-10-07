package json

import "testing"

func Test_mergeHostPath(t *testing.T) {
	tests := []struct {
		name string
		host string
		path string
		want string
	}{
		{
			name: "host only",
			host: "https://example.com",
			want: "https://example.com/",
		},
		{
			name: "host with trailing slash",
			host: "https://example.com/",
			want: "https://example.com/",
		},
		{
			name: "path only",
			path: "/example",
			want: "/example",
		},
		{
			name: "host and path both with slashes",
			host: "https://example.com/",
			path: "/example",
			want: "https://example.com/example",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeHostPath(tt.host, tt.path); got != tt.want {
				t.Fatalf("expecting %q, got %q", tt.want, got)
			}
		})
	}
}
