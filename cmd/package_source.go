package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/webfunction-protocol/webfunction-go"
)

// loadedPackage is a package together with where it came from.
type loadedPackage struct {
	pkg  *webfunction.Package
	from string // where it came from, as a URL: the endpoint, or a file:// URL
	verb string // "Fetched" or "Loaded", for the status line
}

// loadPackage returns the package a command works on: read from file when
// it is set, otherwise fetched from endpointURL.
func loadPackage(endpointURL, file string) (*loadedPackage, error) {
	if file != "" {
		pkg, from, err := loadPackageFile(file)
		if err != nil {
			return nil, err
		}
		return &loadedPackage{pkg: pkg, from: from, verb: "Loaded"}, nil
	}
	client, err := webfunction.FromPackageEndpoint(endpointURL, webfunction.Options{})
	if err != nil {
		return nil, fmt.Errorf("fetching package: %w", err)
	}
	return &loadedPackage{pkg: client.Package(), from: endpointURL, verb: "Fetched"}, nil
}

// loadPackageFile reads a package from a JSON file, returning it with the
// file's absolute file:// URL. The file holds the same document a package
// endpoint returns.
func loadPackageFile(path string) (*webfunction.Package, string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("reading package file: %w", err)
	}
	var pkg webfunction.Package
	if err := json.Unmarshal(raw, &pkg); err != nil {
		return nil, "", fmt.Errorf("parsing package file %s: %w", path, err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return &pkg, fileURL(abs), nil
}

// fileURL returns the file:// URL for an absolute path. Windows drive
// paths (C:\dir\pkg.json) become file:///C:/dir/pkg.json.
func fileURL(path string) string {
	p := filepath.ToSlash(path)
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return (&url.URL{Scheme: "file", Path: p}).String()
}

// packageSourceFlags checks --url and --file for commands that accept
// exactly one of them (convert, validate). It returns the flag name to list
// as missing when neither was given, or an error when both were.
func packageSourceFlags(endpointURL, file string) (missing string, err error) {
	switch {
	case endpointURL != "" && file != "":
		return "", errors.New("--url and --file cannot be used together")
	case endpointURL == "" && file == "":
		return "--url or --file", nil
	}
	return "", nil
}