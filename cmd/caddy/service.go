package caddy

import (
	"fmt"
	"sbx/internal"
	"sbx/shared"

	"github.com/spf13/cobra"
)

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Kelola layanan Caddy (start, stop, restart)",
	Long:  "Perintah untuk enable, disable, start, stop, restart layanan caddy dan membuat systemd internal.",
}

func init() {
	ServiceCmd.AddCommand(&cobra.Command{
        Use:   "enable",
        Short: "Enable layanan caddy",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.EnableService("caddy")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "disable",
        Short: "Disable layanan caddy",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.DisableService("caddy")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "start",
        Short: "Start layanan caddy",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.StartService("caddy")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "stop",
        Short: "Stop layanan caddy",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.StopService("caddy")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "restart",
        Short: "Restart layanan caddy",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.RestartService("caddy")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

	ServiceCmd.AddCommand(&cobra.Command{
        Use:   "status",
        Short: "Status layanan caddy",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.StatusService("caddy")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "create-service",
        Short: "Buat systemd service untuk caddy",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.CreateServiceFile("caddy")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

	ServiceCmd.AddCommand(&cobra.Command{
		Use:   "logs",
		Short: "Ambil log layanan Caddy",
		Run: func(cmd *cobra.Command, args []string) {
			if err := internal.StreamLogs("caddy.service"); err != nil {
				fmt.Println("Gagal menampilkan log:", err)
			}
		},
	})
}