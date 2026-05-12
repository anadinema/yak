package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultConfigDir  = ".config/yak"
	defaultConfigName = "config"
)

// Load reads the config file from the default location or the path override.
// TOML is tried first, then YAML.
func Load(overridePath string) (*Config, error) {
	v := viper.New()

	if overridePath != "" {
		v.SetConfigFile(overridePath)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not determine home directory: %w", err)
		}
		v.AddConfigPath(filepath.Join(home, defaultConfigDir))
		v.SetConfigName(defaultConfigName)
		// viper will try .toml then .yaml automatically
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return nil, fmt.Errorf("no config file found at ~/.config/yak/config.toml (or .yaml) — run 'yak setup --help' to get started")
		}
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks required fields and logical consistency.
func validate(cfg *Config) error {
	var errs []string

	if cfg.DefaultAccount == "" {
		errs = append(errs, "default_account is required")
	}
	if cfg.DefaultRole == "" {
		errs = append(errs, "default_role is required")
	}
	if cfg.Region == "" {
		errs = append(errs, "region is required")
	}
	if cfg.SSOStartURL == "" {
		errs = append(errs, "sso_start_url is required")
	}
	if cfg.SSOSessionName == "" {
		errs = append(errs, "sso_session_name is required")
	}
	if len(cfg.Roles) == 0 {
		errs = append(errs, "at least one role must be defined in [roles]")
	}
	if len(cfg.Accounts) == 0 {
		errs = append(errs, "at least one account must be defined in [[accounts]]")
	}

	// Validate default_role references a known role
	if _, ok := cfg.Roles[cfg.DefaultRole]; !ok && cfg.DefaultRole != "" {
		errs = append(errs, fmt.Sprintf("default_role %q is not defined in [roles]", cfg.DefaultRole))
	}

	// Validate default_account references a known account
	found := false
	for _, a := range cfg.Accounts {
		if a.Name == cfg.DefaultAccount {
			found = true
			break
		}
	}
	if !found && cfg.DefaultAccount != "" {
		errs = append(errs, fmt.Sprintf("default_account %q is not defined in [[accounts]]", cfg.DefaultAccount))
	}

	// Validate each account
	for _, a := range cfg.Accounts {
		if a.Name == "" {
			errs = append(errs, "every account must have a name")
			continue
		}
		if a.AccountID == "" {
			errs = append(errs, fmt.Sprintf("account %q: account_id is required", a.Name))
		}
		// If allowed_roles is set, default_role must be in it
		effectiveDefault := a.DefaultRole
		if effectiveDefault == "" {
			effectiveDefault = cfg.DefaultRole
		}
		if len(a.AllowedRoles) > 0 {
			allowed := false
			for _, r := range a.AllowedRoles {
				if r == effectiveDefault {
					allowed = true
					break
				}
			}
			if !allowed {
				errs = append(errs, fmt.Sprintf(
					"account %q: default_role %q must be listed in allowed_roles",
					a.Name, effectiveDefault,
				))
			}
		}
		// Validate role_overrides reference known roles
		for k := range a.RoleOverrides {
			if _, ok := cfg.Roles[k]; !ok {
				errs = append(errs, fmt.Sprintf(
					"account %q: role_override key %q is not defined in [roles]",
					a.Name, k,
				))
			}
		}
	}

	// Validate safeguards max_power_grade_allowed is a known role
	if cfg.Safeguards.Enabled && cfg.Safeguards.MaxPowerGradeAllowed != "" {
		if PowerGrade(cfg.Safeguards.MaxPowerGradeAllowed) == -1 {
			errs = append(errs, fmt.Sprintf(
				"safeguards.max_power_grade_allowed %q is not a recognised role name",
				cfg.Safeguards.MaxPowerGradeAllowed,
			))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func ResolveAccountAlias(cfg *Config, input string) string {
	if realAccount, ok := cfg.Aliases.Accounts[input]; ok {
		return realAccount
	}
	return input
}

func ResolveRoleAlias(cfg *Config, input string) string {
	if realRole, ok := cfg.Aliases.Roles[input]; ok {
		return realRole
	}
	return input
}

func FindAccount(cfg *Config, name string) *Account {
	name = ResolveAccountAlias(cfg, name)
	for i := range cfg.Accounts {
		if cfg.Accounts[i].Name == name {
			return &cfg.Accounts[i]
		}
	}
	return nil
}

func EffectiveRegion(cfg *Config, accountCfg *Account) string {
	if accountCfg.Region != "" {
		return accountCfg.Region
	}
	return cfg.Region
}

func EffectiveDefaultRole(cfg *Config, accountCfg *Account) string {
	if accountCfg.DefaultRole != "" {
		return accountCfg.DefaultRole
	}
	return cfg.DefaultRole
}
