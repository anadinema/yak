package cmd

import (
	"fmt"
	"os"

	"github.com/anadinema/yak/internal/audit"
	"github.com/anadinema/yak/internal/config"
	"github.com/anadinema/yak/internal/resolver"
	"github.com/anadinema/yak/internal/state"
	"github.com/anadinema/yak/internal/version"
	"github.com/spf13/cobra"
)

var configPath string

type appContext struct {
	cfg            *config.Config
	awsConfigState *state.AWSConfigState
}

func (appCtx *appContext) loadConfig() error {
	if appCtx.cfg != nil {
		return nil
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	appCtx.cfg = cfg
	resolver.SetCacheDir(cfg.Paths.CacheDir)
	state.SetStateDir(cfg.Paths.StateDir)
	audit.SetLogDir(cfg.Paths.LogDir)
	return nil
}

func (appCtx *appContext) loadAWSConfigState() error {
	if appCtx.awsConfigState != nil {
		return nil
	}
	cfg, err := state.LoadAWSConfigState()
	if err != nil {
		return err
	}
	appCtx.awsConfigState = cfg
	return nil
}

func Root() *cobra.Command {
	app := &appContext{}

	root := &cobra.Command{
		Use:           "yak",
		Short:         "Yet Another AWS Kit — AWS SSO and credential management done right",
		Version:       version.String(),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVar(&configPath, "config", "", "path to config file (default: ~/.config/yak/config.toml)")
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if shouldSkipConfigLoad(cmd, args) {
			return nil
		}
		return app.loadConfig()
	}

	root.AddCommand(
		setupCmd(app),
		loginCmd(app),
		accountCmd(app),
		roleCmd(app),
		exportCmd(app),
		statusCmd(app),
		accountsCmd(app),
		rolesCmd(app),
	)

	root.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}

	return root
}

func shouldSkipConfigLoad(cmd *cobra.Command, args []string) bool {
	if cmd.Name() == "help" {
		return true
	}
	if cmd.Parent() == nil && len(args) == 0 {
		return true
	}
	versionFlag, err := cmd.Flags().GetBool("version")
	return err == nil && versionFlag
}

func Execute() {
	root := Root()
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
