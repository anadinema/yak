package state

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// SSOConfig contains parsed role mappings from ~/.aws/config profile sections.
type SSOConfig struct {
	profileRoleNames map[string]string
}

// RoleForProfile returns the current sso_role_name for a profile from AWS config.
func (s *SSOConfig) RoleForProfile(profileName string) string {
	if s == nil {
		return ""
	}
	return s.profileRoleNames[profileName]
}

// LoadSSOConfig reads AWS config and maps [profile X] sections to sso_role_name.
// Missing files return an empty mapping.
func LoadSSOConfig(path string) (*SSOConfig, error) {
	path = expandHomePath(path)
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &SSOConfig{profileRoleNames: map[string]string{}}, nil
		}
		return nil, err
	}

	file, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	state := &SSOConfig{profileRoleNames: make(map[string]string)}
	for _, section := range file.Sections() {
		name := strings.TrimSpace(section.Name())
		if !strings.HasPrefix(name, "profile ") {
			continue
		}
		profileName := strings.TrimSpace(strings.TrimPrefix(name, "profile "))
		roleName := strings.TrimSpace(section.Key("sso_role_name").String())
		if profileName != "" && roleName != "" {
			state.profileRoleNames[profileName] = roleName
		}
	}
	return state, nil
}

func expandHomePath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
