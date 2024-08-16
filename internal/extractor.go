package internal

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractTarGz mengekstrak file .tar.gz ke direktori tujuan dengan logika khusus
func ExtractTarGz(src, dest string) error {
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

	// Menyimpan entri file dalam folder sementara
	fileEntries := make(map[string]*tar.Header)

	// Membaca file tar dan menyimpan entri file
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // Selesai membaca
		}
		if err != nil {
			return fmt.Errorf("error reading tar file: %w", err)
		}

		// Abaikan file LICENSE dan README.md
		if filepath.Base(header.Name) == "LICENSE" || filepath.Base(header.Name) == "README.md" {
			continue
		}

		// Simpan entri file
		fileEntries[header.Name] = header
	}

	// Pastikan ada file yang diekstrak
	if len(fileEntries) == 0 {
		return fmt.Errorf("no files to extract")
	}

	// Ekstrak file yang terpilih
	for name := range fileEntries {
		// Menentukan path target tanpa struktur direktori
		target := filepath.Join(dest, filepath.Base(name))
		
		if !isValidPath(dest, target) {
			return fmt.Errorf("invalid path: %s", target)
		}

		// Membuat direktori jika belum ada
		if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory for file: %w", err)
		}

		// Mengatur ulang gzip reader untuk membaca ulang tarball
		file.Seek(0, io.SeekStart)
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("error creating gzip reader: %w", err)
		}
		defer gzipReader.Close()
		tarReader = tar.NewReader(gzipReader)

		// Cari dan ekstrak file yang sesuai
		for {
			currentHeader, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error reading tar file: %w", err)
			}

			if currentHeader.Name == name {
				// Ekstrak file
				err = writeFile(target, tarReader, currentHeader.FileInfo().Mode())
				if err != nil {
					return fmt.Errorf("error writing file: %w", err)
				}
				break
			}
		}
	}

	return nil
}

// writeFile menulis data dari tarReader ke file target dengan buffering
func writeFile(target string, tarReader io.Reader, fileMode os.FileMode) error {
	outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileMode)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer outFile.Close()

	// Menggunakan buffered writer untuk performa yang lebih baik
	bufWriter := bufio.NewWriterSize(outFile, 64*1024) // Buffer 64 KB
	defer bufWriter.Flush()

	_, err = io.Copy(bufWriter, tarReader)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func isValidPath(dest, target string) bool {
	absDest, err := filepath.Abs(dest)
	if err != nil {
		return false
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return false
	}

	// Ensure both paths end with a trailing separator for accurate comparison
	if !strings.HasSuffix(absDest, string(filepath.Separator)) {
		absDest += string(filepath.Separator)
	}

	return strings.HasPrefix(absTarget, absDest)
}