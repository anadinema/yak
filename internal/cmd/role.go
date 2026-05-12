package cmd

import (
	"fmt"

	"github.com/anadinema/yak/internal/audit"
	"github.com/anadinema/yak/internal/aws"
	"github.com/anadinema/yak/internal/config"
	"github.com/anadinema/yak/internal/state"
	"github.com/spf13/cobra"
)

func roleCmd(app *appContext) *cobra.Command {
	var accountFlag string
	var bypassSafeguards bool

	cmd := &cobra.Command{
		Use:   "role <tier>",
		Short: "Switch the active role tier",
		Long: `Switches the active role tier and regenerates ~/.aws/config.
Accepts full role tier names or any configured alias.

Protected accounts will block roles above max_power_grade_allowed.
Use --bypass-safeguards or -b to override (requires confirmation).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := app.cfg
			if err := app.loadSSOConfig(); err != nil {
				return err
			}

			targetRole := config.ResolveRoleAlias(cfg, args[0])

			if _, ok := cfg.Roles[targetRole]; !ok {
				return fmt.Errorf("role %q is not defined in [roles]", targetRole)
			}

			st, err := state.Load()
			if err != nil {
				return err
			}

			loginAccountName := st.ActiveAccount
			if loginAccountName == "" {
				loginAccountName = cfg.DefaultAccount
			}

			var scopedAccountName string
			var scopedAccounts []*config.Account
			if accountFlag != "" {
				scopedAccountName = config.ResolveAccountAlias(cfg, accountFlag)
				account := config.FindAccount(cfg, scopedAccountName)
				if account == nil {
					return fmt.Errorf("account %q not found in config", scopedAccountName)
				}
				scopedAccounts = []*config.Account{account}
				loginAccountName = scopedAccountName
			} else {
				scopedAccounts = make([]*config.Account, 0, len(cfg.Accounts))
				for i := range cfg.Accounts {
					scopedAccounts = append(scopedAccounts, &cfg.Accounts[i])
				}
			}

			violatingSafeguards := make([]string, 0, len(scopedAccounts))
			for _, account := range scopedAccounts {
				if !aws.IsRoleAllowed(account, targetRole) {
					return fmt.Errorf(
						"role %q is not in the allowed_roles list for account %q.\nAllowed: %v",
						targetRole, account.Name, account.AllowedRoles,
					)
				}

				if err := aws.CheckSafeguard(cfg, account, targetRole); err != nil {
					if !bypassSafeguards {
						return err
					}
					violatingSafeguards = append(violatingSafeguards, account.Name)
				}
			}

			if len(violatingSafeguards) > 0 {
				if len(violatingSafeguards) == 1 {
					if err := confirmBypass(violatingSafeguards[0], targetRole, cfg); err != nil {
						return err
					}
				} else {
					if err := confirmBypassMany(violatingSafeguards, targetRole, cfg); err != nil {
						return err
					}
				}
				for _, accountName := range violatingSafeguards {
					_ = audit.Log(accountName, targetRole)
				}
			}

			resolved, err := resolveAccounts(cfg, scopedAccountName, targetRole)

			if err != nil {
				return err
			}

			if err := aws.WriteConfig(cfg.AWS.ConfigPath, cfg.SSOSessionName, resolved); err != nil {
				return fmt.Errorf("writing AWS config: %w", err)
			}

			if err := aws.SsoLoginForProfile(loginAccountName); err != nil {
				return fmt.Errorf("SSO login failed: %w", err)
			}

			st.ActiveRole = targetRole
			if err := state.Save(st); err != nil {
				return fmt.Errorf("saving state: %w", err)
			}

			fmt.Printf("Switched to role: %s (power grade: %d)\n",
				targetRole, config.PowerGrade(targetRole))

			return nil
		},
	}

	cmd.Flags().StringVarP(&accountFlag, "account", "a", "", "apply role to only this account (default: apply to all accounts)")
	cmd.Flags().BoolVarP(&bypassSafeguards, "bypass-safeguards", "b", false, "override safeguard restrictions (requires confirmation)")
	return cmd
}
