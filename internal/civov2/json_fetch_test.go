package civov2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServer(t *testing.T, method, path string, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("%s %s", method, path), handler)

	return httptest.NewServer(mux)
}

func Test_civoError(t *testing.T) {
	t.Run("known error code", func(t *testing.T) {
		err := &CivoError{Code: "database_account_not_found"}

		if got, want := err.Error(), "authentication failed: invalid token"; got != want {
			t.Fatalf("expecting error to be %q, got %q", want, got)
		}
	})

	t.Run("unknown error code", func(t *testing.T) {
		err := &CivoError{Code: "unknown_error"}

		if got, want := err.Error(), "unknown_error: unknown Civo error"; got != want {
			t.Fatalf("expecting error to be %q, got %q", want, got)
		}
	})

	t.Run("is error", func(t *testing.T) {
		err := &CivoError{Code: "database_account_not_found"}

		if !err.Is(&CivoError{Code: "database_account_not_found"}) {
			t.Fatalf("expecting error to be the same when compared via object")
		}

		if !errors.Is(err, &CivoError{Code: "database_account_not_found"}) {
			t.Fatalf("expecting error to be the same when compared via errors.Is")
		}
	})
}

func Test_jsonClient_get(t *testing.T) {
	t.Run("successful get json request and decoding", func(t *testing.T) {
		token := "foobarbaz"

		handler := func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", token); got != want {
				t.Fatalf("expecting Authorization header to be %q, got %q", want, got)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"test"}`))
		}

		srv := createServer(t, http.MethodGet, "/users/me", handler)
		defer srv.Close()

		client := newCivoJSONClient(nil, srv.URL, token)

		var output struct {
			Name string `json:"name"`
		}

		if err := client.doCivo(context.Background(), "/users/me", http.MethodGet, nil, &output, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.Name != "test" {
			t.Fatalf("expecting value for name to be %q, got %q", "test", output.Name)
		}
	})

	t.Run("successfuul get json request with query parameters", func(t *testing.T) {
		token := "foobarbaz"

		handler := func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", token); got != want {
				t.Fatalf("expecting Authorization header to be %q, got %q", want, got)
			}

			if got, want := r.URL.Query().Get("name"), "test"; got != want {
				t.Fatalf("expecting query parameter name to be %q, got %q", want, got)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"test"}`))
		}

		srv := createServer(t, http.MethodGet, "/users/me", handler)
		defer srv.Close()

		client := newCivoJSONClient(nil, srv.URL, token)

		var output struct {
			Name string `json:"name"`
		}

		params := map[string]string{
			"name": "test",
		}

		if err := client.doCivo(context.Background(), "/users/me", http.MethodGet, nil, &output, params); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.Name != "test" {
			t.Fatalf("expecting value for name to be %q, got %q", "test", output.Name)
		}
	})

	t.Run("successful post json request and decoding", func(t *testing.T) {
		token := "foobarbaz"

		handler := func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", token); got != want {
				t.Fatalf("expecting Authorization header to be %q, got %q", want, got)
			}

			var input struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				t.Fatalf("unable to decode body: %v", err)
			}

			if input.Name != "test" {
				t.Fatalf("expecting value for name to be %q, got %q", "test", input.Name)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"name":"%s"}`, input.Name)
		}

		srv := createServer(t, http.MethodGet, "/users/me", handler)
		defer srv.Close()

		client := newCivoJSONClient(nil, srv.URL, token)

		input := map[string]string{
			"name": "test",
		}

		var output struct {
			Name string `json:"name"`
		}

		if err := client.doCivo(context.Background(), "/users/me", http.MethodGet, input, &output, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.Name != "test" {
			t.Fatalf("expecting value for name to be %q, got %q", "test", output.Name)
		}
	})

	t.Run("unauthorized request", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}

		srv := createServer(t, http.MethodGet, "/users/me", handler)
		defer srv.Close()

		client := newCivoJSONClient(nil, srv.URL, "")

		err := client.doCivo(context.Background(), "/users/me", http.MethodGet, nil, nil, nil)
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})

	t.Run("unexpected status code", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}

		srv := createServer(t, http.MethodGet, "/users/me", handler)
		defer srv.Close()

		client := newCivoJSONClient(nil, srv.URL, "")

		err := client.doCivo(context.Background(), "/users/me", http.MethodGet, nil, nil, nil)
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})

	t.Run("unencodeable json post body", func(t *testing.T) {
		client := newCivoJSONClient(nil, "https://example.com", "")

		body := make(chan struct{})

		err := client.doCivo(context.Background(), "/users/me", http.MethodGet, body, nil, nil)
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})

	t.Run("unable to parse requested url", func(t *testing.T) {
		client := newCivoJSONClient(nil, ":", "")

		err := client.doCivo(context.Background(), "/users/me", http.MethodGet, nil, nil, nil)
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})

	t.Run("unable to create request", func(t *testing.T) {
		client := newCivoJSONClient(nil, "https://example.com", "")

		err := client.doCivo(context.Background(), "/users/me", "GE\nT", http.MethodGet, nil, nil)
		t.Logf("error: %v", err)
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})

	t.Run("unable to do request", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatalf("unexpected request")
		}))
		srv.Close()

		client := newCivoJSONClient(nil, srv.URL, "")

		err := client.doCivo(context.Background(), "/users/me", http.MethodGet, nil, nil, nil)
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})

	t.Run("undecodeable response body", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{`))
		}

		srv := createServer(t, http.MethodGet, "/users/me", handler)
		defer srv.Close()

		client := newCivoJSONClient(nil, srv.URL, "")

		var output struct {
			Name string `json:"name"`
		}

		err := client.doCivo(context.Background(), "/users/me", http.MethodGet, nil, &output, nil)
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})

	t.Run("emulate civo unauthorized with a 404", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"code":"database_account_not_found","message":"account not found"}`))
		}

		srv := createServer(t, http.MethodGet, "/users/me", handler)
		defer srv.Close()

		client := newCivoJSONClient(nil, srv.URL, "")

		err := client.doCivo(context.Background(), "/users/me", http.MethodGet, nil, nil, nil)
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}

		if err.Error() != "authentication failed: invalid token" {
			t.Fatalf("expecting error to be %q, got %q", "authentication failed: invalid token", err.Error())
		}
	})

	t.Run("not found status with unparseable body if civo failed to encode", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{`))
		}

		srv := createServer(t, http.MethodGet, "/users/me", handler)
		defer srv.Close()

		client := newCivoJSONClient(nil, srv.URL, "")

		err := client.doCivo(context.Background(), "/users/me", http.MethodGet, nil, nil, nil)
		if err == nil {
			t.Fatalf("expecting an error, got nil")
		}
	})
}

func Test_getClient(t *testing.T) {
	t.Run("return custom client", func(t *testing.T) {
		client := &http.Client{}
		jsonClient := newCivoJSONClient(client, "", "")

		if got, want := jsonClient.getClient(), client; got != want {
			t.Fatalf("expecting client to be %v, got %v", want, got)
		}
	})

	t.Run("return default client", func(t *testing.T) {
		jsonClient := newCivoJSONClient(nil, "", "")

		if got, want := jsonClient.getClient(), http.DefaultClient; got != want {
			t.Fatalf("expecting client to be %v, got %v", want, got)
		}
	})
}
