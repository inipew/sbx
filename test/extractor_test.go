package test

import (
	"os"
	"path/filepath"
	"sbx/internal"
	"testing"
)

func TestExtractTarGz(t *testing.T) {
	// src := "/home/5678cxz/test-folder/sing-box-1.9.3-linux-amd64.tar.gz"
	dest := "/home/5678cxz/test-folder"

	// err := extractor.ExtractTarGz(src, dest)
	// if err != nil {
	// 	t.Fatalf("Extraction failed: %v", err)
	// }
	src2 := "/home/5678cxz/test-folder/caddy_2.8.4_linux_amd64.tar.gz"

	err2 := internal.ExtractTarGz(src2, dest)
	if err2 != nil {
		t.Fatalf("Extraction failed: %v", err2)
	}

	// Periksa keberadaan file LICENSE
	expectedFile1 := filepath.Join(dest, "caddy")
	if _, err := os.Stat(expectedFile1); os.IsNotExist(err) {
		t.Errorf("Expected file does not exist: %s", expectedFile1)
	}

	// Periksa keberadaan file sing-box
	// expectedFile2 := filepath.Join(dest, "sing-box")
	// if _, err := os.Stat(expectedFile2); os.IsNotExist(err) {
	// 	t.Errorf("Expected file does not exist: %s", expectedFile2)
	// }
}
