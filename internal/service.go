package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"sbx/shared"
	"strings"
	"syscall"
)

// executeCommand menjalankan perintah dan mengembalikan output sebagai string serta error.
func executeCommand(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out // Capture stderr as well
	err := cmd.Run()
	return out.String(), err
}

// GetLogs retrieves the logs for a given service and writes them to the provided writer.
func GetLogs(serviceName string, w io.Writer) error {
	cmd := exec.Command("journalctl", "-xeu", serviceName, "--no-pager")
	cmd.Stdout = w
	cmd.Stderr = w // Capture stderr as well

	return cmd.Run()
}

// CheckServiceStatus checks if the specified service is active.
func CheckServiceStatus(serviceName string) (bool, error) {
	cmd := exec.Command("systemctl", "is-active", "--quiet", serviceName)
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

// GetCurrentVersion retrieves the current version of the specified binary and returns it.
func GetCurrentVersion(binaryPath string) (string, error) {
	cmd := exec.Command(binaryPath, "version")
	output, err := executeCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute '%s version': %w", binaryPath, err)
	}
	return strings.TrimSpace(output), nil
}

func GetCaddyVersion() (string, error) {
	output, err := GetCurrentVersion(CaddyBinPath)
	if err != nil {
		return "", fmt.Errorf("failed to get current version of Caddy: %w", err)
	}

	// Extract and format the version from the command output
	re := regexp.MustCompile(`v([\d\.]+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return "", fmt.Errorf("version not found in output")
	}

	return matches[1], nil
}

// GetSingBoxVersion retrieves, extracts, and formats the version of the sing-box binary.
func GetSingBoxVersion() (string, error) {
	// Execute the command to get the version
	output, err := GetCurrentVersion(SingboxBinPath)
	if err != nil {
		return "", fmt.Errorf("failed to get current version of Sing-box: %w", err)
	}

	// Extract and format the version from the command outputs
	re := regexp.MustCompile(`version (\S+)`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", fmt.Errorf("version not found in output")
	}

	version := match[1]
	versionParts := strings.Split(version, "-")
	if len(versionParts) > 2 {
		return strings.Join(versionParts[:len(versionParts)-1], "-"), nil
	}

	return version, nil
}

func RunSystemdCommand(serviceName, action string) error {
    output, err := shared.RunCommand("systemctl", action, serviceName)
	if err != nil {
		return err
	}
	shared.Info(fmt.Sprintf("Successfully executed systemctl %s for %s.", action, serviceName))
	if output != "" {
		shared.Info(fmt.Sprintf("Command output: %s", output))
	}
	return nil
}

// CreateServiceFile creates a systemd service file for the specified service name.
func CreateServiceFile(serviceName string) error {
	var filePath, serviceContent string

	// Determine file path and content based on service name
	switch serviceName {
	case "caddy":
		filePath = CaddyServicePath
		serviceContent = CaddyServiceContent
	case "sing-box":
		filePath = SingBoxServicePath
		serviceContent = SingBoxServiceContent
	default:
		return fmt.Errorf("unsupported service name: %s", serviceName)
	}

	// Create the service file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create service file '%s': %w", filePath, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("Warning: failed to close service file '%s': %v", filePath, cerr)
		}
	}()

	// Write content to the service file
	_, err = file.WriteString(serviceContent)
	if err != nil {
		return fmt.Errorf("failed to write to service file '%s': %w", filePath, err)
	}

	return nil
}

// EnableService enables the systemd service and returns an error if it fails.
func EnableService(serviceName string) error {
	return RunSystemdCommand(serviceName, "enable")
}

// DisableService disables the systemd service and returns an error if it fails.
func DisableService(serviceName string) error {
	return RunSystemdCommand(serviceName, "disable")
}

// StartService starts the systemd service and returns an error if it fails.
func StartService(serviceName string) error {
	return RunSystemdCommand(serviceName, "start")
}

// StopService stops the systemd service and returns an error if it fails.
func StopService(serviceName string) error {
	return RunSystemdCommand(serviceName, "stop")
}

// StatusService checks the status of the systemd service and returns an error if it fails.
func StatusService(serviceName string) error {
	return RunSystemdCommand(serviceName, "status")
}

// RestartService restarts the systemd service and returns an error if it fails.
func RestartService(serviceName string) error {
	return RunSystemdCommand(serviceName, "restart")
}

func StreamLogs(serviceName string) error {
	pr, pw := io.Pipe()

	// Membuat channel untuk menangani sinyal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Memulai goroutine untuk menjalankan perintah dan menulis output ke pipe
	go func() {
		defer pw.Close()
		if err := GetLogs(serviceName, pw); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to retrieve %s logs: %w", serviceName, err))
		}
	}()

	// Menampilkan log dengan paging
	cmd := exec.Command("less", "-F", "-X", "-R")
	cmd.Stdin = pr
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start 'less': %w", err)
	}

	// Menunggu sinyal atau proses selesai
	select {
	case <-signalChan:
		// Menghentikan proses jika sinyal diterima
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	case err := <-cmdDone(cmd):
		if err != nil {
			return fmt.Errorf("failed to run 'less': %w", err)
		}
	}

	return nil
}

// cmdDone mengembalikan channel yang menerima error dari cmd.Wait
func cmdDone(cmd *exec.Cmd) <-chan error {
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	return done
}