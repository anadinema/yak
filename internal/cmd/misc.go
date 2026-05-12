package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anadinema/yak/internal/aws"
	"github.com/anadinema/yak/internal/config"
	"github.com/anadinema/yak/internal/resolver"
	"github.com/anadinema/yak/internal/state"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// ── export ────────────────────────────────────────────────────────────────────

func exportCmd(app *appContext) *cobra.Command {
	var accountFlag string
	var allFlag bool

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export credentials to ~/.aws/credentials",
		Long:  `Exports SSO credentials for one account (default: active) or all accounts.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := app.cfg

			if allFlag {
				profiles := make([]string, 0, len(cfg.Accounts))
				for _, a := range cfg.Accounts {
					profiles = append(profiles, a.Name)
				}
				if err := aws.ExportAllCredentials(cfg.AWS.CredentialsPath, profiles); err != nil {
					return err
				}
				fmt.Printf("✓ Exported credentials for %d accounts\n", len(profiles))
				return nil
			}

			target := accountFlag
			if target == "" {
				st, err := state.Load()
				if err != nil {
					return err
				}
				target = st.ActiveAccount
			}
			if target == "" {
				target = cfg.DefaultAccount
			}
			target = config.ResolveAccountAlias(cfg, target)

			if err := aws.ExportCredentials(cfg.AWS.CredentialsPath, target); err != nil {
				return err
			}
			fmt.Printf("✓ Exported credentials for account: %s\n", target)
			return nil
		},
	}

	cmd.Flags().StringVarP(&accountFlag, "account", "a", "", "export credentials for a specific account")
	cmd.Flags().BoolVarP(&allFlag, "all", "f", false, "export credentials for all configured accounts")
	return cmd
}

// ── status ────────────────────────────────────────────────────────────────────

type statusOutput struct {
	Account          string `json:"account"`
	Role             string `json:"role"`
	Region           string `json:"region,omitempty"`
	LastLogin        string `json:"last_login,omitempty"`
	SafeguardActive  bool   `json:"safeguard_active"`
	SafeguardMaxRole string `json:"safeguard_max_role,omitempty"`
}

func statusCmd(app *appContext) *cobra.Command {
	var jsonFlag bool
	var printProfile bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current active account, role, and login state",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := app.cfg

			st, err := state.Load()
			if err != nil {
				return err
			}

			activeAccount := st.ActiveAccount
			if activeAccount == "" {
				activeAccount = cfg.DefaultAccount
			}
			activeRole := st.ActiveRole
			if activeRole == "" {
				activeRole = cfg.DefaultRole
			}

			// --print-profile is used by the shell wrapper
			if printProfile {
				fmt.Println(activeAccount)
				return nil
			}

			account := config.FindAccount(cfg, activeAccount)
			region := cfg.Region
			if account != nil {
				region = config.EffectiveRegion(cfg, account)
			}

			out := statusOutput{
				Account:          activeAccount,
				Role:             activeRole,
				Region:           region,
				SafeguardActive:  cfg.Safeguards.Enabled,
				SafeguardMaxRole: cfg.Safeguards.MaxPowerGradeAllowed,
			}
			if !st.LastLogin.IsZero() {
				out.LastLogin = st.LastLogin.Format(time.RFC3339)
			}

			if jsonFlag {
				return json.NewEncoder(os.Stdout).Encode(out)
			}

			fmt.Printf("\nAccount:     %s\n", out.Account)
			fmt.Printf("Role:        %s\n", out.Role)
			fmt.Printf("Region:      %s\n", out.Region)
			if out.LastLogin != "" {
				t, _ := time.Parse(time.RFC3339, out.LastLogin)
				fmt.Printf("Last login:  %s (%s ago)\n",
					t.Format("2006-01-02 15:04:05 UTC"),
					time.Since(t).Round(time.Minute),
				)
			} else {
				fmt.Printf("Last login:  never\n")
			}
			if cfg.Safeguards.Enabled {
				fmt.Printf("Safeguard:   active (max: %s)\n", cfg.Safeguards.MaxPowerGradeAllowed)
			} else {
				fmt.Printf("Safeguard:   disabled\n")
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&jsonFlag, "json", "o", false, "output as JSON")
	cmd.Flags().BoolVarP(&printProfile, "print-profile", "p", false, "print only the active AWS profile name (used by the shell wrapper)")
	return cmd
}

// ── accounts ──────────────────────────────────────────────────────────────────

func accountsCmd(app *appContext) *cobra.Command {
	var jsonFlag bool

	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "List all configured accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := app.cfg

			if jsonFlag {
				return json.NewEncoder(os.Stdout).Encode(cfg.Accounts)
			}

			re, _ := resolver.New(cfg.Secrets.CacheTTL, false)

			tbl := table.New("NAME", "ID", "DEFAULT ROLE", "ALLOWED ROLES", "SAFEGUARD")
			tbl.WithWriter(os.Stdout)
			for _, a := range cfg.Accounts {
				resolvedId, err := re.Resolve(a.AccountID)
				var id string
				if err != nil {
					id = maskAccountID(a.AccountID)
				} else {
					id = maskAccountID(resolvedId)
				}
				defaultRole := a.DefaultRole
				if defaultRole == "" {
					defaultRole = cfg.DefaultRole
				}
				allowed := "all"
				if len(a.AllowedRoles) > 0 {
					allowed = strings.Join(a.AllowedRoles, ", ")
				}
				safeguard := "—"
				if cfg.Safeguards.Enabled {
					if a.Safeguard != nil && !*a.Safeguard {
						safeguard = "disabled"
					} else {
						safeguard = "active"
					}
				}
				tbl.AddRow(a.Name, id, defaultRole, allowed, safeguard)
			}
			tbl.Print()
			return nil
		},
	}

	cmd.Flags().BoolVarP(&jsonFlag, "json", "o", false, "output as JSON")
	return cmd
}

// ── roles ─────────────────────────────────────────────────────────────────────

func rolesCmd(app *appContext) *cobra.Command {
	var accountFlag string
	var jsonFlag bool

	cmd := &cobra.Command{
		Use:   "roles",
		Short: "List available role tiers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := app.cfg

			roleAliases := make(map[string][]string)
			for alias, role := range cfg.Aliases.Roles {
				roleAliases[role] = append(roleAliases[role], alias)
			}

			var allowedFilter map[string]bool
			if accountFlag != "" {
				accountName := config.ResolveAccountAlias(cfg, accountFlag)
				account := config.FindAccount(cfg, accountName)
				if account == nil {
					return fmt.Errorf("account %q not found", accountFlag)
				}
				if len(account.AllowedRoles) > 0 {
					allowedFilter = make(map[string]bool)
					for _, r := range account.AllowedRoles {
						allowedFilter[r] = true
					}
				}
			}

			if jsonFlag {
				return json.NewEncoder(os.Stdout).Encode(cfg.Roles)
			}

			tbl := table.New("TIER", "POWER GRADE", "ALIASES", "AVAILABLE")
			tbl.WithWriter(os.Stdout)
			// Print in power grade order
			order := []string{"read_only", "developer", "admin", "root"}
			for _, tier := range order {
				if _, ok := cfg.Roles[tier]; !ok {
					continue
				}
				available := "yes"
				if allowedFilter != nil && !allowedFilter[tier] {
					available = "no (not in allowed_roles)"
				}
				aliases := strings.Join(roleAliases[tier], ", ")
				if aliases == "" {
					aliases = "—"
				}
				tbl.AddRow(tier, config.PowerGrade(tier), aliases, available)
			}
			tbl.Print()
			return nil
		},
	}

	cmd.Flags().StringVarP(&accountFlag, "account", "a", "", "show only roles available for a specific account")
	cmd.Flags().BoolVarP(&jsonFlag, "json", "o", false, "output as JSON")
	return cmd
}

// ── helpers ───────────────────────────────────────────────────────────────────

// maskAccountID partially masks a 12-digit account ID for safe display.
// e.g. "123456789012" → "123456****12"
func maskAccountID(id string) string {
	if len(id) < 6 {
		return id
	}
	return id[:6] + "****" + id[len(id)-2:]
}

// confirmBypass prompts user for confirmation to bypass safeguards.
func confirmBypass(accountName, role string, cfg *config.Config) error {
	ceiling := cfg.Safeguards.MaxPowerGradeAllowed
	_, _ = fmt.Fprintf(os.Stderr, "\n")
	_, _ = fmt.Fprintf(os.Stderr, "WARNING: Safeguard bypass requested\n")
	_, _ = fmt.Fprintf(os.Stderr, "Account:         %s\n", accountName)
	_, _ = fmt.Fprintf(os.Stderr, "Requested role:  %s (power grade %d)\n",
		role, config.PowerGrade(role))
	_, _ = fmt.Fprintf(os.Stderr, "Safeguard limit: %s (power grade %d)\n",
		ceiling, config.PowerGrade(ceiling))
	_, _ = fmt.Fprintf(os.Stderr, "Continue? [y/N] ")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		return fmt.Errorf("aborted")
	}
	return nil
}

// confirmBypassMany prompts user once for bypass affecting multiple accounts.
func confirmBypassMany(accountNames []string, role string, cfg *config.Config) error {
	ceiling := cfg.Safeguards.MaxPowerGradeAllowed
	_, _ = fmt.Fprintf(os.Stderr, "\n")
	_, _ = fmt.Fprintf(os.Stderr, "WARNING: Safeguard bypass requested\n")
	_, _ = fmt.Fprintf(os.Stderr, "Accounts:        %s\n", strings.Join(accountNames, ", "))
	_, _ = fmt.Fprintf(os.Stderr, "Requested role:  %s (power grade %d)\n",
		role, config.PowerGrade(role))
	_, _ = fmt.Fprintf(os.Stderr, "Safeguard limit: %s (power grade %d)\n",
		ceiling, config.PowerGrade(ceiling))
	_, _ = fmt.Fprintf(os.Stderr, "Continue? [y/N] ")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		return fmt.Errorf("aborted")
	}
	return nil
}

// resolveAccounts builds resolved AWS account entries.
// If targetRole is provided:
//   - targetAccountName set: only that account uses targetRole
//   - targetAccountName empty: all accounts use targetRole
//
// Otherwise, accounts use their effective default role.
func resolveAccounts(cfg *config.Config, targetAccountName, targetRole string) ([]aws.ResolvedAccount, error) {
	res, err := resolver.New(cfg.Secrets.CacheTTL, false)
	if err != nil {
		return nil, err
	}

	resolved := make([]aws.ResolvedAccount, 0, len(cfg.Accounts))
	for i := range cfg.Accounts {
		acc := &cfg.Accounts[i]
		roleTier := config.EffectiveDefaultRole(cfg, acc)
		if targetRole != "" && (targetAccountName == "" || acc.Name == targetAccountName) {
			roleTier = targetRole
		}

		ra, err := aws.ResolveAccount(cfg, acc, roleTier, res)
		if err != nil {
			return nil, err
		}
		resolved = append(resolved, *ra)
	}

	return resolved, nil
}
