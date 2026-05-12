package cmd

import (
	"fmt"
	"os"

	"github.com/anadinema/yak/internal/config"
	"github.com/anadinema/yak/internal/state"
	"github.com/anadinema/yak/internal/version"
	"github.com/spf13/cobra"
)

var configPath string

type appContext struct {
	cfg       *config.Config
	ssoConfig *state.SSOConfig
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
	return nil
}

func (appCtx *appContext) loadSSOConfig() error {
	if appCtx.ssoConfig != nil {
		return nil
	}
	cfg, err := state.LoadSSOConfig(appCtx.cfg.AWS.ConfigPath)
	if err != nil {
		return err
	}
	appCtx.ssoConfig = cfg
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
