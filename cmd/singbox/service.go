package singbox

import (
	"fmt"
	"sbx/internal"
	"sbx/shared"

	"github.com/spf13/cobra"
)

var ServiceCmd = &cobra.Command{
    Use:   "service [command]",
    Short: "Kelola layanan sing-box",
    Long:  "Perintah untuk enable, disable, start, stop, restart layanan sing-box dan membuat systemd internal.",
}

func init() {
    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "enable",
        Short: "Enable layanan sing-box",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.EnableService("sing-box")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "disable",
        Short: "Disable layanan sing-box",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.DisableService("sing-box")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "start",
        Short: "Start layanan sing-box",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.StartService("sing-box")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "stop",
        Short: "Stop layanan sing-box",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.StopService("sing-box")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "status",
        Short: "Status layanan sing-box",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.StatusService("sing-box")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "restart",
        Short: "Restart layanan sing-box",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.RestartService("sing-box")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

    ServiceCmd.AddCommand(&cobra.Command{
        Use:   "create-service",
        Short: "Buat systemd service untuk sing-box",
        Run:   func(cmd *cobra.Command, args []string) {
			err := internal.CreateServiceFile("sing-box")
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
    })

	ServiceCmd.AddCommand(&cobra.Command{
		Use:   "logs",
		Short: "Ambil log layanan Caddy",
		Run: func(cmd *cobra.Command, args []string) {
			if err := internal.StreamLogs("sing-box.service"); err != nil {
				fmt.Println("Gagal menampilkan log:", err)
			}
		},
	})
}