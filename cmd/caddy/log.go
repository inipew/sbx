package caddy

import (
	"context"
	"fmt"
	"sbx/internal"
	"sbx/shared"
	"time"

	"github.com/spf13/cobra"
)

var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "Melihat log Caddy",
	Run:   func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			err := internal.WatchLog(ctx, internal.CaddylogFilePath)
			if err != nil {
				shared.Error(fmt.Sprintf("Error watching log file: %v", err))
			}
		},
}
