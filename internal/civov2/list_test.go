package civov2

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func Test_GetPaginatedGeneric(t *testing.T) {
	t.Run("should return all instances in the Civo account", func(t *testing.T) {
		type paginated struct {
			Page    int        `json:"page"`
			PerPage int        `json:"per_page"`
			Pages   int        `json:"pages"`
			Items   []Instance `json:"items"`
		}

		token := "foobarbaz"

		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+token {
				t.Fatalf("unexpected Authorization header: %s", r.Header.Get("Authorization"))
			}

			w.Header().Set("Content-Type", "application/json")

			switch r.URL.Query().Get("page") {
			case "1":
				json.NewEncoder(w).Encode(paginated{
					Page:    1,
					PerPage: 2,
					Pages:   3,
					Items: []Instance{
						{ID: "1", Name: "test1"},
						{ID: "2", Name: "test2"},
					},
				})

			case "2":
				json.NewEncoder(w).Encode(paginated{
					Page:    2,
					PerPage: 2,
					Pages:   3,
					Items: []Instance{
						{ID: "3", Name: "test3"},
						{ID: "4", Name: "test4"},
					},
				})

			case "3":
				json.NewEncoder(w).Encode(paginated{
					Page:    3,
					PerPage: 2,
					Pages:   3,
					Items: []Instance{
						{ID: "5", Name: "test5"},
					},
				})

			case "4":
				// in Civo, a page overflow returns the first page
				json.NewEncoder(w).Encode(paginated{
					Page:    1,
					PerPage: 2,
					Pages:   3,
					Items: []Instance{
						{ID: "1", Name: "test1"},
						{ID: "2", Name: "test2"},
					},
				})
			}
		}

		srv := createServer(t, http.MethodGet, "/v2/instances", http.HandlerFunc(handler))
		defer srv.Close()

		client := &Client{
			requester: newCivoJSONClient(nil, srv.URL, token),
		}

		ctx := context.Background()

		// deliberately calling the generic function instead of `client.GetInstances()`
		// since all the functions use the exact same code to fetch paginated results
		// so there shouldn't be a need to test each one of them individually
		instances, err := getPaginated[Instance](ctx, client, "/v2/instances")
		if err != nil {
			t.Fatalf("not expecting an error, got %v", err)
		}

		if got, want := len(instances), 5; got != want {
			t.Fatalf("expecting %d instances, got %d", want, got)
		}

		expected := []Instance{
			{ID: "1", Name: "test1"},
			{ID: "2", Name: "test2"},
			{ID: "3", Name: "test3"},
			{ID: "4", Name: "test4"},
			{ID: "5", Name: "test5"},
		}

		for i, instance := range instances {
			if got, want := instance, expected[i]; got != want {
				t.Fatalf("expecting instance %d to be %v, got %v", i, want, got)
			}
		}
	})
}

func Test_Paginated(t *testing.T) {
	t.Run("successfully fetch 3 pages of instances on nyc1 with overflow for page 4", func(t *testing.T) {
		const region = "nyc1"

		handler := func(w http.ResponseWriter, r *http.Request) {
			foundRegion := r.URL.Query().Get("region")
			if foundRegion != region {
				t.Fatalf("http request received unexpected region: %q", foundRegion)
			}

			switch r.URL.Query().Get("page") {
			case "1":
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"page":1,"per_page":2,"pages":3,"items":[{"id":"1","name":"test1"},{"id":"2","name":"test2"}]}`))

			case "2":
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"page":2,"per_page":2,"pages":3,"items":[{"id":"3","name":"test3"},{"id":"4","name":"test4"}]}`))

			case "3":
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"page":3,"per_page":2,"pages":3,"items":[{"id":"5","name":"test5"}]}`))

			default:
				t.Fatalf("unexpected page number: %s", r.URL.Query().Get("page"))
			}
		}

		srv := createServer(t, http.MethodGet, "/v2/instances", handler)
		defer srv.Close()

		client := &Client{
			region:    region,
			requester: newCivoJSONClient(nil, srv.URL, ""),
		}

		type example struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		instances, err := getPaginated[example](context.Background(), client, "/v2/instances")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := len(instances), 5; got != want {
			t.Fatalf("expecting %d instances, got %d", want, got)
		}
	})

	t.Run("single page with no items", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"page":1,"per_page":2,"pages":1,"items":[]}`))
		}

		srv := createServer(t, http.MethodGet, "/v2/instances", handler)
		defer srv.Close()

		client := &Client{
			requester: newCivoJSONClient(nil, srv.URL, ""),
		}

		type example struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		instances, err := getPaginated[example](context.Background(), client, "/v2/instances")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := len(instances), 0; got != want {
			t.Fatalf("expecting %d instances, got %d", want, got)
		}
	})

	t.Run("error fetching items", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}

		srv := createServer(t, http.MethodGet, "/v2/instances", handler)
		defer srv.Close()

		client := &Client{
			requester: newCivoJSONClient(nil, srv.URL, ""),
		}

		type example struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		_, err := getPaginated[example](context.Background(), client, "/v2/instances")
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})
}
