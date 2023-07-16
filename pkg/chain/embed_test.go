package chain_test

import (
	"embed"
	"path/filepath"
)

//go:embed testdata/**
var testDataFS embed.FS

func read(path string) ([]byte, error) {
	return testDataFS.ReadFile(filepath.Clean(path))
}
