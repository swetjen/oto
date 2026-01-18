package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

// TestExampleGeneration tests that the example project can be generated
// and compiled successfully.
func TestExampleGeneration(t *testing.T) {
	is := is.New(t)

	// Get the absolute path to the project root
	projectRoot, err := os.Getwd()
	is.NoErr(err)

	exampleDir := filepath.Join(projectRoot, "example")

	// Step 1: Build the oto CLI tool
	buildCmd := exec.Command("go", "build", "-o", "oto", ".")
	buildCmd.Dir = projectRoot
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build oto CLI: %v\nOutput: %s", err, buildOutput)
	}
	defer os.Remove(filepath.Join(projectRoot, "oto"))

	// Step 2: Add a replace directive to use local otohttp
	// This is needed before generation because the parser resolves dependencies
	goModPath := filepath.Join(exampleDir, "go.mod")
	originalGoMod, err := os.ReadFile(goModPath)
	is.NoErr(err)

	// Restore original go.mod after test
	defer func() {
		if err := os.WriteFile(goModPath, originalGoMod, 0644); err != nil {
			t.Logf("warning: failed to restore go.mod: %v", err)
		}
	}()

	modEditCmd := exec.Command("go", "mod", "edit",
		"-replace=github.com/pacedotdev/oto/otohttp=../otohttp")
	modEditCmd.Dir = exampleDir
	modEditOutput, err := modEditCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to add replace directive: %v\nOutput: %s", err, modEditOutput)
	}

	// Step 3: Run go mod tidy to ensure dependencies are resolved
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = exampleDir
	tidyOutput, err := tidyCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run go mod tidy: %v\nOutput: %s", err, tidyOutput)
	}

	// Step 4: Generate the server code (replicates generate.sh)
	genServerCmd := exec.Command(
		filepath.Join(projectRoot, "oto"),
		"-template", "server.go.plush",
		"-out", "server.gen.go",
		"-pkg", "main",
		"-types", "types.yaml",
		"./def",
	)
	genServerCmd.Dir = exampleDir
	genServerOutput, err := genServerCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to generate server code: %v\nOutput: %s", err, genServerOutput)
	}

	// Step 5: Generate the JavaScript client (replicates generate.sh)
	genClientCmd := exec.Command(
		filepath.Join(projectRoot, "oto"),
		"-template", "client.js.plush",
		"-out", "client.gen.js",
		"-pkg", "main",
		"-types", "types.yaml",
		"./def",
	)
	genClientCmd.Dir = exampleDir
	genClientOutput, err := genClientCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to generate JS client: %v\nOutput: %s", err, genClientOutput)
	}

	// Step 6: Format the generated Go code
	fmtCmd := exec.Command("gofmt", "-w", "server.gen.go")
	fmtCmd.Dir = exampleDir
	fmtOutput, err := fmtCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to format generated code: %v\nOutput: %s", err, fmtOutput)
	}

	// Step 7: Verify the example project compiles
	compileCmd := exec.Command("go", "build", "-o", "/dev/null", ".")
	compileCmd.Dir = exampleDir
	compileOutput, err := compileCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("example project failed to compile: %v\nOutput: %s", err, compileOutput)
	}

	// Step 8: Verify the generated files exist
	generatedFiles := []string{
		filepath.Join(exampleDir, "server.gen.go"),
		filepath.Join(exampleDir, "client.gen.js"),
	}
	for _, f := range generatedFiles {
		_, err := os.Stat(f)
		if os.IsNotExist(err) {
			t.Errorf("expected generated file to exist: %s", f)
		}
	}
}
