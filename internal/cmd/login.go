package cmd

import (
	"fmt"

	"github.com/anadinema/yak/internal/aws"
	"github.com/anadinema/yak/internal/config"
	"github.com/anadinema/yak/internal/state"
	"github.com/spf13/cobra"
)

func loginCmd(app *appContext) *cobra.Command {
	var accountFlag string
	var roleFlag string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "SSO login and credential export",
		Long:  `Runs 'aws sso login' for the active account, then exports credentials.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := app.cfg
			if err := app.loadSSOConfig(); err != nil {
				return err
			}

			st, err := state.Load()
			if err != nil {
				return err
			}

			targetAccountName := st.ActiveAccount
			if accountFlag != "" {
				targetAccountName = config.ResolveAccountAlias(cfg, accountFlag)
			}
			if targetAccountName == "" {
				targetAccountName = cfg.DefaultAccount
			}

			account := config.FindAccount(cfg, targetAccountName)
			if account == nil {
				return fmt.Errorf("account %q not found in config", targetAccountName)
			}

			targetRole := ""
			if cfg.UseRoleFromState {
				targetRole = st.ActiveRole
			}
			if roleFlag != "" {
				targetRole = config.ResolveRoleAlias(cfg, roleFlag)
			}
			if targetRole == "" {
				targetRole = config.EffectiveDefaultRole(cfg, account)
			}
			configRole := app.ssoConfig.RoleForProfile(targetAccountName)
			roleChangedFromConfig := targetRole != configRole

			if roleChangedFromConfig {
				resolved, err := resolveAccounts(cfg, targetAccountName, targetRole)
				if err != nil {
					return err
				}
				if err := aws.WriteConfig(cfg.AWS.ConfigPath, cfg.SSOSessionName, resolved); err != nil {
					return fmt.Errorf("writing AWS config: %w", err)
				}
			}

			fmt.Printf("-> Logging in to account: %s (role: %s)\n", targetAccountName, targetRole)

			if err := aws.SsoLoginForProfile(targetAccountName); err != nil {
				return fmt.Errorf("SSO login failed: %w", err)
			}

			if err := aws.ExportCredentials(cfg.AWS.CredentialsPath, targetAccountName); err != nil {
				return fmt.Errorf("exporting credentials: %w", err)
			}

			if err := state.SetLogin(targetAccountName, targetRole); err != nil {
				return fmt.Errorf("saving state: %w", err)
			}

			fmt.Printf("Login successful! (%s / %s)\n", targetAccountName, targetRole)
			fmt.Printf("\nExport: export AWS_PROFILE=%s\n", targetAccountName)
			return nil
		},
	}

	cmd.Flags().StringVar(&accountFlag, "account", "", "login as a specific account instead of the active one")
	cmd.Flags().StringVar(&roleFlag, "role", "", "use a specific role tier for this login")
	return cmd
}
