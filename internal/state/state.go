package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// State tracks the currently active account and role tier across invocations.
type State struct {
	ActiveAccount string    `json:"active_account"`
	ActiveRole    string    `json:"active_role"`
	LastLogin     time.Time `json:"last_login,omitempty"`
}

var stateDirOverride string

func SetStateDir(path string) {
	stateDirOverride = path
}

func stateFile() (string, error) {
	if stateDirOverride != "" {
		return filepath.Join(expandConfiguredDir(stateDirOverride), "state.json"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".local", "share", "yak", "state.json"), nil
}

// Load reads the state file. Returns an empty State if the file doesn't exist.
func Load() (*State, error) {
	path, err := stateFile()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &State{}, nil
		}
		return nil, fmt.Errorf("could not read state file: %w", err)
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("could not parse state file: %w", err)
	}
	return &s, nil
}

// Save writes the state to disk.
func Save(s *State) error {
	path, err := stateFile()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("could not create state directory: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal state: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// SetAccount updates the active account and saves.
func SetAccount(account string) error {
	s, err := Load()
	if err != nil {
		return err
	}
	s.ActiveAccount = account
	return Save(s)
}

// SetRole updates the active role and saves.
func SetRole(role string) error {
	s, err := Load()
	if err != nil {
		return err
	}
	s.ActiveRole = role
	return Save(s)
}

// SetLogin updates the active account, role, and last login timestamp.
func SetLogin(account, role string) error {
	s, err := Load()
	if err != nil {
		return err
	}
	s.ActiveAccount = account
	s.ActiveRole = role
	s.LastLogin = time.Now().UTC()
	return Save(s)
}

func expandConfiguredDir(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
