package cmd

import (
	"boilerplate/internal/app"
	"boilerplate/internal/pkg/logger"
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

	ctx := logger.NewTagContext(context.Background(), "__main__")
	err := app.Run(ctx, app.SetConfigFile(cfgFile))
	if err != nil {
		logger.WithContext(ctx).Errorln(err.Error())
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file")
}
