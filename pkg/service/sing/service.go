package service

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

// StartService starts the Caddy service and returns an error if it fails.
func StartService() error {
	cmd := exec.Command("systemctl", "start", "sing-box.service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Caddy service: %w", err)
	}
	return nil
}

// StopService stops the Caddy service and returns an error if it fails.
func StopService() error {
	cmd := exec.Command("systemctl", "stop", "sing-box.service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Caddy service: %w", err)
	}
	return nil
}

// RestartService restarts the Caddy service and returns an error if it fails.
func RestartService() error {
	cmd := exec.Command("systemctl", "restart", "sing-box.service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart Caddy service: %w", err)
	}
	return nil
}

// GetLogs retrieves the logs for the Caddy service and returns them as a string.
func GetLogs(w io.Writer) error {
	cmd := exec.Command("journalctl", "-u", "sing-box.service", "--no-pager")
	cmd.Stdout = w
	cmd.Stderr = w // Capture stderr as well

	return cmd.Run()
}

// GetCurrentVersion retrieves the current version of the Caddy binary and returns it.
func GetCurrentVersion() (string, error) {
	cmd := exec.Command("/usr/local/bin/sing-box", "version")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("gagal menjalankan perintah 'sing-box version': %w", err)
	}
	version := strings.TrimSpace(string(out))
	return version, nil
}

// checkServiceStatus checks if a service is active.
func CheckServiceStatus() (bool, error) {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "sing-box.service")
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 3 {
			// Service is not active
			return false, nil
		}
		return false, err
	}
	// Service is active
	return true, nil
}

// ExtractVersion mengambil dan memformat versi dari output perintah
func ExtractVersion(output string) (string, error) {
	// Menyusun regular expression untuk menemukan versi
	// Format dari versi yang dicontohkan: 1.10.0-alpha.29-f585961e
	// Kita hanya akan mengambil bagian sebelum tanda '-' terakhir
	re := regexp.MustCompile(`version (\S+)`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", fmt.Errorf("version not found in output")
	}

	version := match[1]
	// Hapus tag revision dari versi
	versionParts := strings.Split(version, "-")
	if len(versionParts) > 2 {
		return strings.Join(versionParts[:len(versionParts)-1], "-"), nil
	}

	return version, nil
}

// GetSingBoxVersion menjalankan perintah dan mengembalikan versi terformat
func GetSingBoxVersion() (string, error) {
	// Menjalankan perintah untuk mendapatkan versi
	cmd := exec.Command("/usr/local/bin/sing-box", "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running command: %w", err)
	}

	// Mengambil output dan mengekstrak versi
	version, err := ExtractVersion(out.String())
	if err != nil {
		return "", fmt.Errorf("error extracting version: %w", err)
	}

	return version, nil
}