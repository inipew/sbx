package caddy

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sbx/pkg/constant"
	service "sbx/pkg/service/caddy"
	"sbx/shared"

	"github.com/spf13/cobra"
)

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Kelola layanan Caddy (start, stop, restart)",
	Long:  "Perintah untuk enable, disable, start, stop, restart layanan caddy dan membuat systemd service.",
}

func init() {
	ServiceCmd.AddCommand(&cobra.Command{
        Use:   "enable",
        Short: "Enable layanan sing-box",
        Run:   enableService,
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "disable",
        Short: "Disable layanan sing-box",
        Run:   disableService,
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "start",
        Short: "Start layanan sing-box",
        Run:   startService,
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "stop",
        Short: "Stop layanan sing-box",
        Run:   stopService,
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "restart",
        Short: "Restart layanan sing-box",
        Run:   restartService,
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "create-service",
        Short: "Buat systemd service untuk sing-box",
        Run:   createServiceFile,
    })

	ServiceCmd.AddCommand(&cobra.Command{
		Use:   "logs",
		Short: "Ambil log layanan Caddy",
		Run: func(cmd *cobra.Command, args []string) {
			if err := streamLogs(); err != nil {
				fmt.Println("Gagal menampilkan log:", err)
			}
		},
	})
}
func enableService(cmd *cobra.Command, args []string) {
    runSystemdCommand("enable")
}

func disableService(cmd *cobra.Command, args []string) {
    runSystemdCommand("disable")
}

func startService(cmd *cobra.Command, args []string) {
    runSystemdCommand("start")
}

func stopService(cmd *cobra.Command, args []string) {
    runSystemdCommand("stop")
}

func restartService(cmd *cobra.Command, args []string) {
    runSystemdCommand("restart")
}

// Fungsi untuk menampilkan log dengan paging menggunakan 'less'
func streamLogs() error {
	pr, pw := io.Pipe()

	// Memulai goroutine untuk menjalankan perintah dan menulis output ke pipe
	go func() {
		defer pw.Close()
		if err := service.GetLogs(pw); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to retrieve Caddy logs: %w", err))
		}
	}()

	// Menampilkan log dengan paging
	cmd := exec.Command("less")
	cmd.Stdin = pr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func createServiceFile(cmd *cobra.Command, args []string) {
	filePath := constant.SingBoxServicePath
	serviceContent := constant.SingBoxServiceContent

	file, err := os.Create(filePath)
	if err != nil {
		shared.Error(fmt.Sprintf("Error creating service file: %v", err))
		return
	}
	defer file.Close()

	_, err = file.WriteString(serviceContent)
	if err != nil {
		shared.Error(fmt.Sprintf("Error writing service file: %v", err))
		return
	}

	shared.Info("Systemd service file created successfully.")
	shared.Info("Remember to run `systemctl daemon-reload` and `systemctl enable sing-box.service`.")
}

func runSystemdCommand(action string) {
    cmd := exec.Command("systemctl", action, "caddy.service")
    err := cmd.Run()
    if err != nil {
        shared.Error(fmt.Sprintf("Error executing systemctl %s: %v", action, err))
        return
    }
    shared.Info(fmt.Sprintf("Successfully executed systemctl %s for caddy service.", action))
}