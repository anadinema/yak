package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anadinema/yak/internal/config"
	"github.com/spf13/cobra"
)

const (
	configFormatTOML = "toml"
	configFormatYAML = "yaml"
)

func generateCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a minimal yak config file template",
		Long: `Creates a minimal config file template.

When YAK_CONFIG_FILE is set, that path is used.
Otherwise, writes to ~/.config/yak/config.toml or ~/.config/yak/config.yaml
based on --format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format = strings.ToLower(strings.TrimSpace(format))
			if format != configFormatTOML && format != configFormatYAML {
				return fmt.Errorf("unsupported format %q: use toml or yaml", format)
			}

			path, err := config.ResolveGeneratePath(format)
			if err != nil {
				return err
			}
			path = expandHome(path)

			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("config file already exists: %s", path)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("could not check config path %q: %w", path, err)
			}

			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				return fmt.Errorf("could not create config directory for %q: %w", path, err)
			}

			var content string
			switch format {
			case configFormatYAML:
				content = minimalYAMLConfig()
			default:
				content = minimalTOMLConfig()
			}

			if err := os.WriteFile(path, []byte(content), 0600); err != nil {
				return fmt.Errorf("could not write config template to %q: %w", path, err)
			}

			fmt.Printf("✓ Config template written to %s\n", path)
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", configFormatTOML, "output format: toml or yaml")
	return cmd
}

func minimalTOMLConfig() string {
	return `default_account = "dev"
default_role    = "developer"

[aws]
credentials_path = "~/.aws/credentials"
config_path      = "~/.aws/config"
region           = ""
sso_start_url    = ""
sso_session_name = "yak"

[roles]
developer = ""

[[accounts]]
name       = "dev"
account_id = ""
`
}

func minimalYAMLConfig() string {
	return `default_account: dev
default_role: developer

aws:
  credentials_path: ~/.aws/credentials
  config_path: ~/.aws/config
  region: ""
  sso_start_url: ""
  sso_session_name: yak

roles:
  developer: ""

accounts:
  - name: dev
    account_id: ""
`
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
