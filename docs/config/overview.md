# Configuration Overview

yak reads a single configuration file on startup.

## File location

yak looks for its config file in this order:

1. `~/.config/yak/config.toml`
2. `~/.config/yak/config.yaml`

If both exist, TOML takes precedence. The directory is created automatically on first run.

---

## Format

Both TOML and YAML are fully supported. The structure is identical — choose whichever you prefer.

- **TOML** is the default and what all examples in this documentation use. It has no indentation-sensitive syntax and handles inline tables cleanly.
- **YAML** may feel more familiar if you work with Kubernetes or GitHub Actions configs day-to-day.

See the [examples](examples/toml-env.md) for all four combinations.

---

## Value formats

Any string value in the config supports three formats. yak detects the format automatically based on the prefix.

### Plain text

Used as-is. Appropriate for non-sensitive values like region names and file paths.

```toml
region = "eu-north-1"
```

### Environment variable — `${VAR_NAME}`

Resolved from the process environment at runtime.

```toml
account_id = "${AWS_DEV_ACCOUNT_ID}"
```

!!! warning
    yak exits with an error if a referenced environment variable is not set. There is no silent fallback to an empty string.

### 1Password — `op://vault/item/field`

Fetched via the `op` CLI at runtime. Requires `op` to be installed and the user to be authenticated.

```toml
account_id = "op://Private/AWS Accounts/dev_account_id"
```

Results are cached in `~/.local/share/yak/secret_cache` (mode `0600`) for `secrets.cache_ttl` minutes. `yak login` always bypasses this cache.

!!! info "Mixing formats"
    You can mix formats freely within a single config file. For example, use plain text for a region, env vars for role names, and 1Password for account IDs.

---

## Runtime files

yak maintains two files at runtime that are separate from your config and never hand-edited.

| File                              | Purpose                                       |
|-----------------------------------|-----------------------------------------------|
| `~/.local/share/yak/state.json`   | Active account and role tier, last login time |
| `~/.local/share/yak/secret_cache` | Cached `op://` values with timestamps         |
| `~/.local/share/yak/audit.log`    | Record of safeguard bypasses                  |

---

## Sections

| Section              | Purpose                                                      |
|----------------------|--------------------------------------------------------------|
| Global defaults      | `default_account`, `default_role`, `region`, `sso_start_url` |
| `[aws]`              | Paths for generated AWS files                                |
| `[secrets]`          | `op://` cache TTL                                            |
| `[passthrough]`      | Forwarding unknown commands to the AWS CLI                   |
| `[safeguards]`       | Production account protection                                |
| `[roles]`            | Role tier definitions and values                             |
| `[[accounts]]`       | Per-account configuration                                    |
| `[aliases.accounts]` | Short aliases for account names                              |
| `[aliases.roles]`    | Short aliases for role tier names                            |

See the [full reference](reference.md) for every field.
