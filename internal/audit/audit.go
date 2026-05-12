package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logDirOverride string

func SetLogDir(path string) {
	logDirOverride = path
}

// Log appends a bypass event to the audit log file.
func Log(account, requestedRole string) error {
	path, err := auditFile()

	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}

	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f,
		"%s  SAFEGUARD BYPASSED  account=%s role=%s user=%s\n",
		time.Now().UTC().Format(time.RFC3339),
		account, requestedRole, username,
	)
	if err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func auditFile() (string, error) {
	if logDirOverride != "" {
		return filepath.Join(expandConfiguredDir(logDirOverride), "audit.log"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "yak", "audit.log"), nil
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
