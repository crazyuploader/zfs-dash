package cmd

import (
	"fmt"

	"github.com/crazyuploader/zfs-dash/internal/config"
	"github.com/crazyuploader/zfs-dash/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the ZFS dashboard web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := configInitError(); err != nil {
			return fmt.Errorf("read config: %w", err)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}
		if len(cfg.Endpoints) == 0 {
			return fmt.Errorf("no endpoints configured; use --endpoints or config.yaml")
		}
		return server.Start(cfg)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
