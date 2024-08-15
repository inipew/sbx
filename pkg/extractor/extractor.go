package extractor

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ExtractTarGz mengekstrak file .tar.gz ke direktori tujuan
func ExtractTarGz(src string, dest string) error {
	// Membuka file .tar.gz
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Membuat reader gzip
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Membuat reader tar
	tarReader := tar.NewReader(gzipReader)

	// Membaca file tar dan mengekstrak isinya
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // Selesai mengekstrak
		}
		if err != nil {
			return fmt.Errorf("error reading tar file: %w", err)
		}

		// Menentukan path file yang akan diekstrak
		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Membuat direktori jika header adalah direktori
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory: %w", err)
			}
		case tar.TypeReg:
			// Membuat file jika header adalah file
			file, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("error creating file: %w", err)
			}
			_, err = io.Copy(file, tarReader)
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}
			file.Close()
		default:
			return fmt.Errorf("unknown type: %b in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
