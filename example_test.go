package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestExampleGenerateAndBuild(t *testing.T) {
	is := is.New(t)

	root, err := os.Getwd()
	is.NoErr(err)
	binDir := t.TempDir()
	otoPath := filepath.Join(binDir, "oto")
	buildCmd := exec.Command("go", "build", "-o", otoPath, ".")
	buildCmd.Dir = root
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %s . failed: %v\n%s", otoPath, err, buildOut)
	}

	exampleDir := filepath.Join(root, "example")
	is.NoErr(os.Chdir(exampleDir))
	t.Cleanup(func() {
		_ = os.Chdir(root)
	})

	serverTemplate := "server.go.txt"
	serverOut := "server.gen.go"
	defPath := "./def"
	clientTemplate := "client.js.txt"
	clientOut := "client.gen.js"
	typeMap := "type-mapping.yaml"

	serverCmd := exec.Command(otoPath,
		"-template="+serverTemplate,
		"-out="+serverOut,
		"-pkg=main",
		defPath,
	)
	serverCmd.Dir = exampleDir
	serverOutBytes, err := serverCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("oto server generation failed: %v\n%s", err, serverOutBytes)
	}
	gofmtCmd := exec.Command("gofmt", "-w", serverOut)
	gofmtCmd.Dir = exampleDir
	gofmtOut, err := gofmtCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gofmt failed: %v\n%s", err, gofmtOut)
	}

	clientCmd := exec.Command(otoPath,
		"-template="+clientTemplate,
		"-out="+clientOut,
		"-pkg=main",
		"-type-map="+typeMap,
		defPath,
	)
	clientCmd.Dir = exampleDir
	clientOutBytes, err := clientCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("oto client generation failed: %v\n%s", err, clientOutBytes)
	}

	_, err = os.Stat(serverOut)
	is.NoErr(err)
	_, err = os.Stat(clientOut)
	is.NoErr(err)

	cmd := exec.Command("go", "build", ".")
	cmd.Dir = exampleDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build . failed: %v\n%s", err, out)
	}
}
