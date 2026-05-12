package cmd

import (
	"fmt"

	"github.com/anadinema/yak/internal/audit"
	"github.com/anadinema/yak/internal/aws"
	"github.com/anadinema/yak/internal/config"
	"github.com/anadinema/yak/internal/state"
	"github.com/spf13/cobra"
)

func accountCmd(app *appContext) *cobra.Command {
	var roleFlag string
	var bypassSafeguards bool

	cmd := &cobra.Command{
		Use:   "account <name>",
		Short: "Switch the active AWS account",
		Long:  `Switches the active account, and exports credentials. Accepts full account names or any configured alias.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := app.cfg
			if err := app.loadSSOConfig(); err != nil {
				return err
			}

			targetName := config.ResolveAccountAlias(cfg, args[0])
			account := config.FindAccount(cfg, targetName)
			if account == nil {
				return fmt.Errorf("account %q not found in config (checked aliases too)", args[0])
			}

			st, err := state.Load()
			if err != nil {
				return err
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
			configRole := app.ssoConfig.RoleForProfile(targetName)
			roleChangedFromConfig := targetRole != configRole

			if !aws.IsRoleAllowed(account, targetRole) {
				return fmt.Errorf(
					"role %q is not in the allowed_roles list for account %q.\nAllowed: %v",
					targetRole, account.Name, account.AllowedRoles,
				)
			}

			if safeguardError := aws.CheckSafeguard(cfg, account, targetRole); safeguardError != nil {
				if !bypassSafeguards {
					return safeguardError
				}
				if err := confirmBypass(account.Name, targetRole, cfg); err != nil {
					return err
				}
				_ = audit.Log(account.Name, targetRole)
			}

			if roleChangedFromConfig {
				resolved, err := resolveAccounts(cfg, targetName, targetRole)
				if err != nil {
					return err
				}
				if err := aws.WriteConfig(cfg.AWS.ConfigPath, cfg.SSOSessionName, resolved); err != nil {
					return fmt.Errorf("writing AWS config: %w", err)
				}
				if err := aws.SsoLoginForProfile(targetName); err != nil {
					return fmt.Errorf("SSO login failed: %w", err)
				}
			}

			if err := aws.ExportCredentials(cfg.AWS.CredentialsPath, targetName); err != nil {
				return fmt.Errorf("exporting credentials: %w", err)
			}

			st.ActiveAccount = targetName
			st.ActiveRole = targetRole
			if err := state.Save(st); err != nil {
				return fmt.Errorf("saving state: %w", err)
			}

			fmt.Printf("Switched to account: %s (role: %s)\n", targetName, targetRole)
			fmt.Printf("Export: export AWS_PROFILE=%s\n", targetName)
			return nil
		},
	}

	cmd.Flags().StringVar(&roleFlag, "role", "", "also switch to this role tier when changing account")
	cmd.Flags().BoolVar(&bypassSafeguards, "bypass-safeguards", false, "override safeguard restrictions (requires confirmation)")
	return cmd
}
