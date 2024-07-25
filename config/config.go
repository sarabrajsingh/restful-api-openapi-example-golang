package config

import (
	"path/filepath"
	"runtime"
)

type Config struct {
	OpenAPI3YamlFileLocation string
	SwaggerUIFolder          string
}

// NewConfig creates a new Config instance with fully qualified paths
func NewConfig() (*Config, error) {
	// Get the path of the directory where the current source file is located
	_, currentFilePath, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filepath.Dir(currentFilePath))

	// Build the absolute paths
	openAPI3YamlFileLocation := filepath.Join(baseDir, "api", "openapi.yaml")
	swaggerUIFolder := filepath.Join(baseDir, "swaggerui", "dist")

	return &Config{
		OpenAPI3YamlFileLocation: openAPI3YamlFileLocation,
		SwaggerUIFolder:          swaggerUIFolder,
	}, nil
}
