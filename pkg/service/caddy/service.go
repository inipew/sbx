package service

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

// StartService starts the Caddy service and returns an error if it fails.
func StartService() error {
	cmd := exec.Command("systemctl", "start", "caddy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Caddy service: %w", err)
	}
	return nil
}

// StopService stops the Caddy service and returns an error if it fails.
func StopService() error {
	cmd := exec.Command("systemctl", "stop", "caddy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Caddy service: %w", err)
	}
	return nil
}

// RestartService restarts the Caddy service and returns an error if it fails.
func RestartService() error {
	cmd := exec.Command("systemctl", "restart", "caddy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart Caddy service: %w", err)
	}
	return nil
}

// GetLogs retrieves the logs for the Caddy service and returns them as a string.
func GetLogs(w io.Writer) error {
	cmd := exec.Command("journalctl", "-u", "caddy", "--no-pager")
	cmd.Stdout = w
	cmd.Stderr = w // Capture stderr as well

	return cmd.Run()
}

func extractVersion(output string) (string, error) {
	// Menyusun pola regex untuk mengekstrak versi
	re := regexp.MustCompile(`v([\d\.]+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) < 2 {
		return "", fmt.Errorf("versi tidak ditemukan dalam output")
	}

	return matches[1], nil
}
// GetCurrentVersion retrieves the current version of the Caddy binary and returns it.
func GetCurrentVersion() (string, error) {
	cmd := exec.Command("caddy", "version")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("gagal menjalankan perintah 'caddy version': %w", err)
	}
	version, err := extractVersion(string(out))
	if err != nil {
		fmt.Println("Error:", err)
	}
	return version, nil
}

// checkServiceStatus checks if a service is active.
func CheckServiceStatus() (bool, error) {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "caddy.service")
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