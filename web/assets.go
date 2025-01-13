package web

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// AssetType represents the type of asset (JS or CSS)
type AssetType string

const (
	JavaScript AssetType = ".js"
	CSS        AssetType = ".css"
)

// AssetManifest represents the mapping of original asset names to their hashed versions
type AssetManifest map[string]string

var (
	manifest AssetManifest
	mu       sync.RWMutex
)

// Init initializes the asset manifest
func AssetInit() error {
	return loadManifest()
}

// loadManifest reads and parses the asset manifest file
func loadManifest() error {
	manifestPath := filepath.Join("static", "build", "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	return json.Unmarshal(data, &manifest)
}

// reloadManifestIfNeeded reloads the manifest in development mode
func reloadManifestIfNeeded() {
	if os.Getenv("APP_ENV") == "development" {
		loadManifest()
	}
}

// GetAssetPaths returns paths for assets of the specified type, excluding chunks if specified
func GetAssetPaths(assetType AssetType, excludeChunks bool) []string {
	reloadManifestIfNeeded()

	var paths []string
	mu.RLock()
	defer mu.RUnlock()

	for name, path := range manifest {
		if strings.HasSuffix(name, string(assetType)) {
			if excludeChunks && strings.Contains(path, "chunks/") {
				continue
			}
			paths = append(paths, path)
		}
	}
	return paths
}

// GetScriptPaths returns JavaScript files excluding chunk files
func GetScriptPaths() []string {
	return GetAssetPaths(JavaScript, true)
}

// GetStylePaths returns CSS files
func GetStylePaths() []string {
	return GetAssetPaths(CSS, false)
}

// GetPath returns the path for a given asset name
func GetPath(name string) string {
	reloadManifestIfNeeded()

	mu.RLock()
	defer mu.RUnlock()
	return manifest[name]
}
