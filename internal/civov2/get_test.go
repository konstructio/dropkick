package civov2

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func Test_GetByID(t *testing.T) {
	t.Run("get single instance by id", func(t *testing.T) {
		id := "abc123"
		region := "nyc1"

		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v2/instances/"+id {
				t.Fatalf("expecting request to %q, got %q", "/v2/instances/"+id, r.URL.Path)
			}

			if r.URL.Query().Get("region") != region {
				t.Fatalf("expecting region %q, got %q", region, r.URL.Query().Get("region"))
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Instance{
				ID:   id,
				Name: "test-instance",
			})
		}

		srv := createServer(t, http.MethodGet, "/v2/instances/"+id, handler)
		defer srv.Close()

		client := &Client{
			region:    region,
			requester: newCivoJSONClient(nil, srv.URL, region),
		}

		instance, err := getByID[Instance](context.Background(), client, "/v2/instances", id)
		if err != nil {
			t.Fatalf("not expecting an error, got %v", err)
		}

		if instance.ID != id {
			t.Fatalf("expecting instance id %q, got %q", id, instance.ID)
		}
	})

	t.Run("failed to get instance by id", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}

		srv := createServer(t, http.MethodGet, "/v2/instances/foobar", handler)
		defer srv.Close()

		client := &Client{
			region:    "n/a",
			requester: newCivoJSONClient(nil, srv.URL, "n/a"),
		}

		_, err := getByID[Instance](context.Background(), client, "/v2/instances", "foobar")
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})
}
