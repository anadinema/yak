package config

// Config is the top-level structure of ~/.config/yak/config.toml
type Config struct {
	DefaultAccount   string `mapstructure:"default_account"`
	DefaultRole      string `mapstructure:"default_role"`
	Region           string `mapstructure:"region"`
	SSOStartURL      string `mapstructure:"sso_start_url"`
	SSOSessionName   string `mapstructure:"sso_session_name"`
	UseRoleFromState bool   `mapstructure:"use_state_role"`

	AWS        AWSConfig         `mapstructure:"aws"`
	Secrets    SecretsConfig     `mapstructure:"secrets"`
	Safeguards SafeguardsConfig  `mapstructure:"safeguards"`
	Roles      map[string]string `mapstructure:"roles"`
	Accounts   []Account         `mapstructure:"accounts"`
	Aliases    AliasConfig       `mapstructure:"aliases"`
}

type AWSConfig struct {
	CredentialsPath string `mapstructure:"credentials_path"`
	ConfigPath      string `mapstructure:"config_path"`
}

type SecretsConfig struct {
	CacheTTL int `mapstructure:"cache_ttl"` // minutes
}

type SafeguardsConfig struct {
	Enabled              bool     `mapstructure:"enabled"`
	ProtectedAccounts    []string `mapstructure:"protected_accounts"`
	ProtectedPattern     string   `mapstructure:"protected_pattern"`
	MaxPowerGradeAllowed string   `mapstructure:"max_power_grade_allowed"`
}

type Account struct {
	Name          string            `mapstructure:"name"`
	AccountID     string            `mapstructure:"account_id"`
	DefaultRole   string            `mapstructure:"default_role"`
	Region        string            `mapstructure:"region"`
	AllowedRoles  []string          `mapstructure:"allowed_roles"`
	RoleOverrides map[string]string `mapstructure:"role_overrides"`
	Safeguard     *bool             `mapstructure:"safeguard"` // nil = inherit from safeguards config
}

type AliasConfig struct {
	Accounts map[string]string `mapstructure:"accounts"`
	Roles    map[string]string `mapstructure:"roles"`
}

var BuiltinPowerGrades = map[string]int{
	"read_only": 100,
	"developer": 200,
	"admin":     300,
	"root":      400,
}

func PowerGrade(role string) int {
	g, ok := BuiltinPowerGrades[role]
	if !ok {
		return -1
	}
	return g
}
