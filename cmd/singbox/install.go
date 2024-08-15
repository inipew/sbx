package singbox

import (
	"fmt"
	"sbx/pkg/download"
	"sbx/pkg/extractor"
	"sbx/shared"

	"github.com/spf13/cobra"
)


var InstallCmd = &cobra.Command{
    Use:   "install",
    Short: "Unduh dan instal sing-box",
    Run: func(cmd *cobra.Command, args []string) {
        err := download.InstallCaddy()
        if err != nil {
            shared.Error(fmt.Sprintln("Gagal menginstal Sing-box:", err))
            return
        }
        shared.Info(fmt.Sprintln("Sing-box berhasil diinstal."))
    },
}

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract a .tar.gz file to a specified directory",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
            shared.Info(fmt.Sprintln("Usage: extract <source-file> <destination-directory>"))
			return
		}
		src := args[0]
		dest := args[1]

		if err := extractor.ExtractTarGz(src, dest); err != nil {
			shared.Error(fmt.Sprintf("Error extracting file: %v\n", err))
			return
		}

		shared.Info(fmt.Sprint("Extraction complete."))
	},
}