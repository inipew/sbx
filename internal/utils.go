package internal

import (
	"fmt"
	"os"
	"os/exec"
	"sbx/shared"
)

func CheckErr(err error) {
    if err != nil {
        shared.Error(fmt.Sprintln("Error:", err))
        os.Exit(1)
    }
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