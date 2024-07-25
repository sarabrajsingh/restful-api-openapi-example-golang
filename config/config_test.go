package config_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sarabrajsingh/restful-openapi/config"
)

// TestNewConfig tests the NewConfig function
func TestNewConfig(t *testing.T) {
	config, err := config.NewConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Get the path of the directory where the current source file is located
	_, currentFilePath, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filepath.Dir(currentFilePath))

	expectedOpenAPI3YamlFileLocation := filepath.Join(baseDir, "api", "openapi.yaml")
	expectedSwaggerUIFolder := filepath.Join(baseDir, "swaggerui", "dist")

	if config.OpenAPI3YamlFileLocation != expectedOpenAPI3YamlFileLocation {
		t.Errorf("expected OpenAPI3YamlFileLocation to be %v, got %v", expectedOpenAPI3YamlFileLocation, config.OpenAPI3YamlFileLocation)
	}

	if config.SwaggerUIFolder != expectedSwaggerUIFolder {
		t.Errorf("expected SwaggerUIFolder to be %v, got %v", expectedSwaggerUIFolder, config.SwaggerUIFolder)
	}
}
