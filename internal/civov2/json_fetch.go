package civov2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// civoJSONClient is a client that can make requests to the Civo API.
type civoJSONClient struct {
	endpoint    string
	bearerToken string
	client      *http.Client
}

// newCivoJSONClient creates a new civoJSONClient.
func newCivoJSONClient(client *http.Client, endpoint, bearerToken string) *civoJSONClient {
	return &civoJSONClient{
		endpoint:    endpoint,
		bearerToken: bearerToken,
		client:      client,
	}
}

// getClient returns the http client to use for requests as configured
// when creating the civoJSONClient, or the default Go http client if none
// was provided.
func (j *civoJSONClient) getClient() *http.Client {
	if j.client != nil {
		return j.client
	}
	return http.DefaultClient
}

// knownCodes is a map of known Civo error codes to their corresponding
// error messages. This is used to standardize the error messages returned
// by the Civo API.
var knownCodes = map[string]string{
	"database_account_not_found": "authentication failed: invalid token",
}

// CivoError is an error returned by the Civo API.
type CivoError struct {
	Code string `json:"code"`
}

// Error returns the error message for the CivoError.
func (e *CivoError) Error() string {
	if msg, ok := knownCodes[e.Code]; ok {
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
	Code int
}

// Error returns the error message for the HTTPError.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("unexpected status code: \"%d %s\"", e.Code, http.StatusText(e.Code))
}

// Is checks if the target error is an HTTPError with the same code.
func (e *HTTPError) Is(target error) bool {
	err, ok := target.(*HTTPError)
	return ok && err.Code == e.Code
}

// doCivo makes a raw HTTP request to the Civo API.
func (j *civoJSONClient) doCivo(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
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

	res, err := j.getClient().Do(req)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted, http.StatusNotFound:
		// do nothing, we successfully got the data
		// for not found, we need to check the code in the response

	case http.StatusUnauthorized:
		// Standardize the error message for authentication failure
		return &CivoError{Code: "database_account_not_found"}
	default:
		return &HTTPError{Code: res.StatusCode}
	}

	// Civo treats 404 as a special case for authentication failure
	// by returning a specific error code "database_account_not_found"
	// so we need to check for that.
	if res.StatusCode == http.StatusNotFound {
		var resp CivoError
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			return &HTTPError{Code: res.StatusCode}
		}
		return &resp
	}

	if output != nil {
		if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
			return fmt.Errorf("unable to decode response: %w", err)
		}
	}

	return nil
}
