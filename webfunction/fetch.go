package webfunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FetchPackage retrieves a Package by invoking the given URL as a Web
// Function endpoint: an HTTP POST with a JSON object body, per
// https://webfunction.org/endpoint and the "Invocation response" retrieval
// method in https://webfunction.org/package#retrieval.
//
// This does not (yet) send any authentication or version headers - see
// FetchPackageWithOptions for that.
func FetchPackage(url string) (*Package, error) {
	return FetchPackageWithOptions(url, Options{})
}

// Options configures an endpoint invocation.
type Options struct {
	// BearerToken, if set, is sent as "Authorization: Bearer <token>".
	// See https://webfunction.org/authentication.
	BearerToken string

	// APIVersion, if set, is sent as the "Api-Version" header.
	// See https://webfunction.org/versioning.
	APIVersion string
}

// FetchPackageWithOptions is like FetchPackage but allows passing bearer
// auth / API version headers for packages that require them.
func FetchPackageWithOptions(url string, opts Options) (*Package, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, fmt.Errorf("building request for %s: %w", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if opts.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+opts.BearerToken)
	}
	if opts.APIVersion != "" {
		req.Header.Set("Api-Version", opts.APIVersion)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching package from %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %w", url, err)
	}

	// Per https://webfunction.org/endpoint, only 200 and 400 carry
	// protocol-level meaning; anything else must be surfaced as-is rather
	// than treated as a successful package response.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching package from %s: unexpected status %d: %s", url, resp.StatusCode, truncate(body, 500))
	}

	var pkg Package
	if err := json.Unmarshal(body, &pkg); err != nil {
		return nil, fmt.Errorf("parsing package from %s: %w", url, err)
	}

	return &pkg, nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}