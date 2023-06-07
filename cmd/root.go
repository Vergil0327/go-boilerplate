package cmd

import (
	"boilerplate/internal/app"
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "app",
		Short: "backend boilerplate",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if cfgFile == "" {
				return errors.New("missing configuration file")
			}
			return nil
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	app.Run(context.Background(), app.SetConfigFile(cfgFile))
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file")
}
