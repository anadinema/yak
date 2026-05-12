package resolver

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const onePasswordPrefix = "op://"
const envVariablePrefix = "$"

type Resolver struct {
	cacheTTL    time.Duration
	cacheFile   string
	memCache    map[string]cacheEntry
	mu          sync.Mutex
	bypassCache bool // set to true for yak login
}

var cacheDirOverride string

func SetCacheDir(path string) {
	cacheDirOverride = path
}

type cacheEntry struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

func New(cacheTTLMinutes int, bypassCache bool) (*Resolver, error) {
	cacheFile, err := defaultCacheFile()
	if err != nil {
		return nil, err
	}

	r := &Resolver{
		cacheTTL:    time.Duration(cacheTTLMinutes) * time.Minute,
		cacheFile:   cacheFile,
		memCache:    make(map[string]cacheEntry),
		bypassCache: bypassCache,
	}
	if !bypassCache && cacheTTLMinutes > 0 {
		_ = r.loadCache()
	}
	return r, nil
}

func (r *Resolver) Resolve(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	switch {
	case strings.HasPrefix(raw, onePasswordPrefix):
		return r.resolveFrom1Password(raw)
	case strings.HasPrefix(raw, envVariablePrefix):
		return r.resolveFromEnvVar(raw)
	default:
		return raw, nil
	}
}

func (r *Resolver) ResolveAll(raw map[string]string) (map[string]string, error) {
	resolved := make(map[string]string, len(raw))
	for k, v := range raw {
		val, err := r.Resolve(v)
		if err != nil {
			return nil, fmt.Errorf("resolving %q: %w", k, err)
		}
		resolved[k] = val
	}
	return resolved, nil
}

func (r *Resolver) resolveFromEnvVar(raw string) (string, error) {
	// Strip $ or ${}
	name := strings.TrimPrefix(raw, "$")
	if strings.HasPrefix(name, "{") {
		name = name[1 : len(name)-1]
	}

	val, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("environment variable %q is not set (referenced in config as %s)", name, raw)
	}
	return val, nil
}

func (r *Resolver) resolveFrom1Password(ref string) (string, error) {
	// Check op CLI is available
	opPath, err := exec.LookPath("op")
	if err != nil {
		return "", fmt.Errorf(
			"config references %q but the 1Password CLI (op) is not installed or not in PATH.\n"+
				"Install it from: https://developer.1password.com/docs/cli/get-started/",
			ref,
		)
	}

	if !r.bypassCache && r.cacheTTL > 0 {
		r.mu.Lock()
		if entry, ok := r.memCache[ref]; ok && time.Now().Before(entry.ExpiresAt) {
			r.mu.Unlock()
			return entry.Value, nil
		}
		r.mu.Unlock()
	}

	out, err := exec.Command(opPath, "read", "--no-newline", ref).Output()
	if err != nil {
		return "", fmt.Errorf("failed to read %q from 1Password: %w", ref, err)
	}
	value := string(out)

	// Store in cache
	if !r.bypassCache && r.cacheTTL > 0 {
		r.mu.Lock()
		r.memCache[ref] = cacheEntry{
			Value:     value,
			ExpiresAt: time.Now().Add(r.cacheTTL),
		}
		r.mu.Unlock()
		_ = r.saveCache()
	}

	return value, nil
}

func (r *Resolver) FlushCache() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.memCache = make(map[string]cacheEntry)
	return os.Remove(r.cacheFile)
}

func (r *Resolver) loadCache() error {
	data, err := os.ReadFile(r.cacheFile)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return json.Unmarshal(data, &r.memCache)
}

func (r *Resolver) saveCache() error {
	r.mu.Lock()
	data, err := json.Marshal(r.memCache)
	r.mu.Unlock()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(r.cacheFile), 0700); err != nil {
		return err
	}

	return os.WriteFile(r.cacheFile, data, 0600)
}

func defaultCacheFile() (string, error) {
	if cacheDirOverride != "" {
		return filepath.Join(expandConfiguredDir(cacheDirOverride), "secret_cache"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".local", "share", "yak", "secret_cache"), nil
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
