package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ProfileAWSConfigState struct {
	AccountID string `json:"account_id"`
	RoleName  string `json:"role_name"`
	Region    string `json:"region"`
}

type AWSConfigState struct {
	UpdatedAt time.Time                        `json:"updated_at"`
	Profiles  map[string]ProfileAWSConfigState `json:"profiles"`
}

func (s *AWSConfigState) RoleNameForProfile(profileName string) string {
	if s == nil {
		return ""
	}
	profile, ok := s.Profiles[profileName]
	if !ok {
		return ""
	}
	return profile.RoleName
}

func awsConfigStateFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".local", "share", "yak", "config.state.json"), nil
}

func LoadAWSConfigState() (*AWSConfigState, error) {
	path, err := awsConfigStateFile()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &AWSConfigState{Profiles: map[string]ProfileAWSConfigState{}}, nil
		}
		return nil, fmt.Errorf("could not read config state file: %w", err)
	}

	var s AWSConfigState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("could not parse config state file: %w", err)
	}
	if s.Profiles == nil {
		s.Profiles = map[string]ProfileAWSConfigState{}
	}
	return &s, nil
}

func SaveAWSConfigState(s *AWSConfigState) error {
	path, err := awsConfigStateFile()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("could not create config state directory: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config state: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
