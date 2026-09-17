package transit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalProjectConfig(t *testing.T) {
	// Create a temporary directory and change to it
	tmpDir, err := os.MkdirTemp("", "transit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a local transit
	err = CreateTransit("build-local", []string{"echo local-build", "echo done"}, true)
	if err != nil {
		t.Fatalf("CreateTransit failed: %v", err)
	}

	// Verify .transit.yaml was created
	if _, err := os.Stat(filepath.Join(tmpDir, LocalConfigFileName)); os.IsNotExist(err) {
		t.Fatalf("expected %s to be created", LocalConfigFileName)
	}

	// Retrieve local transit
	tr, err := GetTransit("build-local")
	if err != nil {
		t.Fatalf("GetTransit failed: %v", err)
	}
	if !tr.IsLocal {
		t.Error("expected IsLocal to be true")
	}
	if len(tr.Commands) != 2 || tr.Commands[0] != "echo local-build" {
		t.Errorf("unexpected commands: %v", tr.Commands)
	}

	// Search
	results, err := SearchCommands("local-build")
	if err != nil {
		t.Fatalf("SearchCommands failed: %v", err)
	}
	if len(results["build-local"]) != 1 {
		t.Errorf("expected 1 search result, got %v", results)
	}

	// Delete
	err = DeleteTransit("build-local")
	if err != nil {
		t.Fatalf("DeleteTransit failed: %v", err)
	}

	// Should not exist anymore
	_, err = GetTransit("build-local")
	if err == nil {
		t.Fatal("expected transit to be deleted")
	}
}
