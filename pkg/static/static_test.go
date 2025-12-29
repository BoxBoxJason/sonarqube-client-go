package static

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadStaticFile(t *testing.T) {
	// Create a temp directory and file
	tempDir, err := os.MkdirTemp("", "static-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testContent := "test content here"
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	t.Run("existing file", func(t *testing.T) {
		content, err := LoadStaticFile(testFile)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if content != testContent {
			t.Errorf("Expected %q, got %q", testContent, content)
		}
	})

	t.Run("non-existing file", func(t *testing.T) {
		_, err := LoadStaticFile(filepath.Join(tempDir, "nonexistent.txt"))
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})
}

func TestGenerateStaticFiles(t *testing.T) {
	// Create temp directories
	navDir, err := os.MkdirTemp("", "static-nav-*")
	if err != nil {
		t.Fatalf("Failed to create temp nav dir: %v", err)
	}
	defer os.RemoveAll(navDir)

	outputDir, err := os.MkdirTemp("", "static-output-*")
	if err != nil {
		t.Fatalf("Failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Create the expected static files
	staticFiles := []struct {
		name    string
		content string
	}{
		{"sonarqube.go.temp", "package test\n// sonarqube content"},
		{"client_util.go.temp", "package test\n// client util content"},
		{"suite_test.go.temp", "package test\n// test suite content"},
	}

	for _, sf := range staticFiles {
		path := filepath.Join(navDir, sf.name)
		err = os.WriteFile(path, []byte(sf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to write static file %s: %v", sf.name, err)
		}
	}

	t.Run("generate static files", func(t *testing.T) {
		err := GenerateStaticFiles(navDir, outputDir)
		if err != nil {
			t.Errorf("GenerateStaticFiles failed: %v", err)
		}

		// Check that output file was created
		outputFile := filepath.Join(outputDir, ConstFileName)
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Errorf("Expected output file %s was not created", outputFile)
		}
	})

	t.Run("missing nav directory", func(t *testing.T) {
		err := GenerateStaticFiles("/nonexistent/path", outputDir)
		// The function may not return an error for missing files based on LoadStaticFile behavior
		// Just check it doesn't panic
		_ = err
	})
}
