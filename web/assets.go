package web

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type AssetManifest map[string]string

var manifest AssetManifest

func InitAssets() error {
	manifestPath := filepath.Join("static", "build", "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &manifest)
}

func GetAssetPath(name string) string {
	if path, ok := manifest[name]; ok {
		return path
	}
	return ""
}

func GetScriptPaths() []string {
	var paths []string
	for name, path := range manifest {
		if strings.HasSuffix(name, ".js") && !strings.Contains(path, "chunks/") {
			paths = append(paths, path)
		}
	}
	return paths
}

func GetStylePaths() []string {
	var paths []string
	for name, path := range manifest {
		if strings.HasSuffix(name, ".css") {
			paths = append(paths, path)
		}
	}
	return paths
}
