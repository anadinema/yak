package cmd

import (
	"fmt"

	"github.com/anadinema/yak/internal/aws"
	"github.com/spf13/cobra"
)

func setupCmd(app *appContext) *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Generate ~/.aws/config from your yak config",
		Long: `Reads ~/.config/yak/config.toml, resolves all secret values, and writes
a valid ~/.aws/config with one [profile name] block per configured account.

Per-account default_role is used for each profile (falling back to global
default_role), with role_overrides applied when present. Run 'yak setup' after
adding accounts, changing role names, or changing sso_start_url.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := app.cfg

			resolved, err := resolveAccounts(cfg, "", "")
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Println("# Dry run — would write the following to", cfg.AWS.ConfigPath)
				fmt.Println()
				fmt.Printf("[sso-session %s]\n", cfg.SSOSessionName)
				fmt.Printf("sso_start_url            = %s\n", resolved[0].SSOStart)
				fmt.Printf("sso_region               = %s\n", resolved[0].SSORegion)
				fmt.Printf("sso_registration_scopes  = sso:account:access\n")
				fmt.Println()
				for _, ra := range resolved {
					fmt.Printf("[profile %s]\n", ra.Name)
					fmt.Printf("sso_session    = %s\n", cfg.SSOSessionName)
					fmt.Printf("sso_account_id = %s\n", ra.AccountID)
					fmt.Printf("sso_role_name  = %s\n", ra.RoleName)
					fmt.Printf("region         = %s\n", ra.Region)
					fmt.Println()
				}
				return nil
			}

			if err := aws.WriteConfig(cfg.AWS.ConfigPath, cfg.SSOSessionName, resolved); err != nil {
				return err
			}

			fmt.Printf("✓ AWS config written to %s (%d profiles, using per-account default roles)\n",
				cfg.AWS.ConfigPath, len(resolved))
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print what would be written without touching any files")
	return cmd
}
