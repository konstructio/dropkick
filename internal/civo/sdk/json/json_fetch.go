package json

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client is a client that can make requests to the Civo API.
type Client struct {
	endpoint    string
	bearerToken string
	client      *http.Client
}

// New creates a new civoJSONClient.
func New(client *http.Client, endpoint, bearerToken string) *Client {
	return &Client{
		endpoint:    endpoint,
		bearerToken: bearerToken,
		client:      client,
	}
}

// GetClient returns the http client to use for requests as configured
// when creating the civoJSONClient, or the default Go http client if none
// was provided.
func (j *Client) GetClient() *http.Client {
	if j.client != nil {
		return j.client
	}
	return http.DefaultClient
}

// getEndpoint returns the endpoint of the client.
func (j *Client) GetEndpoint() string {
	return j.endpoint
}

// knownCodes is a map of known Civo error codes to their corresponding
// error messages. This is used to standardize the error messages returned
// by the Civo API.
var knownCodes = map[string]string{
	"database_account_not_found":         "authentication failed: invalid token",
	"database_network_inuse_by_instance": "the network is in use",
}

// CivoError is an error returned by the Civo API.
type CivoError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}

// Error returns the error message for the CivoError.
func (e *CivoError) Error() string {
	if msg, ok := knownCodes[e.Code]; ok {
		if e.Details != "" {
			return fmt.Sprintf("%s: %s", msg, e.Details)
		}

		return msg
	}

	return e.Code + ": unknown Civo error"
}

// Is checks if the target error is a CivoError with the same code.
func (e *CivoError) Is(target error) bool {
	err, ok := target.(*CivoError)
	return ok && err.Code == e.Code
}

// HTTPError is an error returned when an unexpected HTTP status code is
// returned by the Civo API.
type HTTPError struct {
	Code     int
	Contents string
}

// Error returns the error message for the HTTPError.
func (e *HTTPError) Error() string {
	msg := fmt.Sprintf("unexpected status code: \"%d %s\"", e.Code, http.StatusText(e.Code))

	if e.Contents != "" {
		return fmt.Sprintf("%s: %s", msg, e.Contents)
	}

	return msg
}

// Is checks if the target error is an HTTPError with the same code.
func (e *HTTPError) Is(target error) bool {
	err, ok := target.(*HTTPError)
	return ok && err.Code == e.Code
}

// Do makes a raw HTTP request to the Civo API.
func (j *Client) Do(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
	u, err := url.Parse(mergeHostPath(j.endpoint, location))
	if err != nil {
		return fmt.Errorf("unable to parse requested url: %w", err)
	}

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", j.bearerToken))

	res, err := j.GetClient().Do(req)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	// Successful cases where we expect the body to have what's supposed
	// to be returned.
	case http.StatusOK, http.StatusAccepted:
		// Do nothing

	// Civo treats 404 as a special case for authentication failure
	// by returning a specific error code "database_account_not_found".
	// For conflicts (like attempting to delete a network in use by a node)
	// Civo returns a 409 status code, but the response body is still JSON.
	case http.StatusNotFound, http.StatusConflict:
		var resp CivoError
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(res.Body) // ignoring error since we can handle it at the HTTPError level
			return &HTTPError{Code: res.StatusCode, Contents: buf.String()}
		}
		return &resp

	case http.StatusUnauthorized:
		// Standardize the error message for authentication failure
		return &CivoError{Code: "database_account_not_found"}
	default:
		var buf bytes.Buffer
		buf.ReadFrom(res.Body)
		return &HTTPError{Code: res.StatusCode, Contents: buf.String()}
	}

	if output != nil {
		if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
			var buf bytes.Buffer
			buf.ReadFrom(res.Body)
			return fmt.Errorf("unable to decode response: %s: %w", buf.String(), err)
		}
	}

	return nil
}
