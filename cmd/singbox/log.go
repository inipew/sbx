package singbox

import (
	"fmt"
	"sbx/pkg/constant"
	"sbx/pkg/log"
	"sbx/shared"

	"github.com/spf13/cobra"
)

var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "Melihat log sing-box",
	Run:   func(cmd *cobra.Command, args []string) {
			err := log.WatchLog(constant.SingboxLogFilePath)
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
}
