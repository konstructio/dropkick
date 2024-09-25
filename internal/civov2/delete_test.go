package civov2

import (
	"context"
	"net/http"
	"testing"
)

func Test_Delete(t *testing.T) {
	t.Run("successfully delete an instance", func(t *testing.T) {
		id := "abc123"

		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Fatalf("unexpected http method: %s", r.Method)
			}

			if r.URL.Path != "/v2/instances/"+id {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id":"abc123","result":"deleted"}`))
		}

		srv := createServer(t, http.MethodDelete, "/v2/instances/"+id, handler)
		defer srv.Close()

		client := &Client{
			requester: newCivoJSONClient(nil, srv.URL, ""),
		}

		if err := delete(client, context.Background(), "/v2/instances", id); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error deleting an instance: not found", func(t *testing.T) {
		id := "abc123"

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}

		srv := createServer(t, http.MethodDelete, "/v2/instances/"+id, handler)
		defer srv.Close()

		client := &Client{
			requester: newCivoJSONClient(nil, srv.URL, ""),
		}

		if err := delete(client, context.Background(), "/v2/instances", id); err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})
}
