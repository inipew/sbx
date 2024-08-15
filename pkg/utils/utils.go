package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func CheckErr(err error) {
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}

// DetectArch returns the architecture of the host system.
func DetectArch() string {
	return runtime.GOARCH
}

// DetectOS returns the operating system of the host system.
func DetectOS() string {
	return runtime.GOOS
}

func CreateDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}
	return nil
}

// RemoveDir removes a directory and its contents.
func RemoveDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to remove directory %s: %w", dir, err)
	}
	return nil
}

// RemoveFile deletes a file at the given path.
func RemoveFile(filepath string) error {
	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to remove file %s: %w", filepath, err)
	}
	return nil
}

// MoveDir moves a directory from the source path to the destination path.
func MoveDir(srcPath, destPath string) error {
	cmd := exec.Command("mv", srcPath, destPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to move directory from %s to %s: %w", srcPath, destPath, err)
	}
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// RemoveDir removes a directory if it exists.
func RemoveDirIfExists(dir string) error {
	if PathExists(dir) {
		return RemoveDir(dir)
	}
	return nil
}

// ExtractTarGz extracts a .tar.gz file to a specified directory.
func ExtractTarGz(tarGzPath, destDir string) error {
	// Open the .tar.gz file
	file, err := os.Open(tarGzPath)
	if err != nil {
		return fmt.Errorf("failed to open tar.gz file %s: %w", tarGzPath, err)
	}
	defer file.Close()

	// Create a new gzip reader
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	// Create a new tar reader
	tr := tar.NewReader(gzr)

	// Iterate through the files in the archive
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		target := filepath.Join(destDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			// Create directories
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
		case tar.TypeReg:
			// Create files
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}
			_, err = io.Copy(outFile, tr)
			outFile.Close()
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}
		default:
			return fmt.Errorf("unsupported tar header type %c", header.Typeflag)
		}
	}
	return nil
}