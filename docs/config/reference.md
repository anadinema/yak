# Configuration Reference

Complete reference for every field in `~/.config/yak/config.toml` (or `config.yaml`).

---

## Global defaults

Top-level fields that apply to all accounts unless overridden at the account level.

| Field             | Type   | Required | Description                                                                      |
|-------------------|--------|----------|----------------------------------------------------------------------------------|
| `default_account` | string | yes      | Account name to use when none is specified. Must match a name in `[[accounts]]`. |
| `default_role`    | string | yes      | Role tier to use when none is specified. Must match a key in `[roles]`.          |
| `region`          | string | yes      | AWS region. Can be overridden per account.                                       |
| `sso_start_url`   | string | yes      | Your AWS SSO portal URL. Supports `${VAR}` and `op://`.                          |

```toml
default_account = "dev"
default_role    = "admin"
region          = "eu-west-1"
sso_start_url   = "${AWS_SSO_URL}"
```

---

## [aws]

Paths for the AWS files that yak generates and writes to.

| Field              | Type   | Default              | Description                                   |
|--------------------|--------|----------------------|-----------------------------------------------|
| `credentials_path` | string | `~/.aws/credentials` | Where to write exported SSO credentials.      |
| `config_path`      | string | `~/.aws/config`      | Where to write the generated AWS config file. |

```toml
[aws]
credentials_path = "~/.aws/credentials"
config_path      = "~/.aws/config"
```

---

## [secrets]

Controls caching of `op://` values. Has no effect on `${VAR}` or plain text values.

| Field       | Type    | Default | Description                                                     |
|-------------|---------|---------|-----------------------------------------------------------------|
| `cache_ttl` | integer | `60`    | Minutes to cache `op://` values. Set to `0` to disable caching. |

```toml
[secrets]
cache_ttl = 60
```

!!! note
    `yak login` always bypasses the cache regardless of this setting.

---

## [passthrough]

Controls whether unknown commands are forwarded to the AWS CLI.

| Field             | Type   | Default | Description                                                                                                                                       |
|-------------------|--------|---------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `enabled`         | bool   | `false` | When `true`, commands yak does not recognise are passed to the AWS CLI verbatim.                                                                  |
| `cli_binary_path` | string | `""`    | Path to the AWS CLI binary. Auto-detected from `PATH` if empty. yak exits with an error if passthrough is enabled but the binary cannot be found. |

```toml
[passthrough]
enabled         = false
cli_binary_path = ""
```

!!! tip "Using yak as an aws drop-in"
    With `enabled = true` and `alias aws=yak` in your shell, all your existing `aws` commands continue to work while you get yak's account and role management on top.

---

## [safeguards]

Guards against accidentally using elevated roles on production-like accounts.

| Field                     | Type     | Default       | Description                                                                                                     |
|---------------------------|----------|---------------|-----------------------------------------------------------------------------------------------------------------|
| `enabled`                 | bool     | `true`        | Master switch. `false` disables all safeguards globally.                                                        |
| `protected_accounts`      | string[] | `[]`          | Exact account names that are protected.                                                                         |
| `protected_pattern`       | string   | `""`          | Regex applied to account names. Any match is treated as protected. Applied in addition to `protected_accounts`. |
| `max_power_grade_allowed` | string   | `"developer"` | Maximum role tier (by name) allowed on protected accounts. Roles with a higher power grade are blocked.         |

```toml
[safeguards]
enabled                 = true
protected_accounts      = ["prod", "mgmt"]
protected_pattern       = ""
max_power_grade_allowed = "developer"
```

!!! warning "Bypassing safeguards"
    Use `yak role <tier> --bypass-safeguards` to override the safeguard for a single operation. yak will display a warning and require explicit confirmation before proceeding. All bypasses are logged to `~/.local/share/yak/audit.log`.

---

## [roles]

Defines the role tiers available across all accounts. Role names are used in `yak role`, `allowed_roles`, and `role_overrides`.

Each role has a built-in **power grade** — a numeric value yak uses to compare role privilege levels for safeguard enforcement. You cannot assign a power grade to a role that is not in this built-in list.

| Role name   | Built-in power grade | Typical use                   |
|-------------|----------------------|-------------------------------|
| `read_only` | 100                  | Diagnostics, read-only access |
| `developer` | 200                  | Day-to-day development work   |
| `admin`     | 300                  | Administrative operations     |
| `root`      | 400                  | Break-glass emergency access  |

!!! danger "root role"
    The `root` role has power grade 400 and bypasses normal restrictions. Treat it as a last resort. Consider excluding it from `allowed_roles` on all accounts and only using it with `--bypass-safeguards` when truly necessary.

### Plain string (recommended)

```toml
[roles]
read_only = "${AWS_VIEWER_ROLE}"
developer = "${AWS_DEVELOPER_ROLE}"
admin     = "${AWS_ELEVATED_ROLE}"
root      = "${AWS_GODMODE_ROLE}"
```

### Table form (only needed to override the built-in power grade)

```toml
[roles]
read_only = { value = "${AWS_VIEWER_ROLE}", power_grade = 100 }
developer = { value = "${AWS_DEVELOPER_ROLE}", power_grade = 200 }
admin     = { value = "${AWS_ELEVATED_ROLE}", power_grade = 300 }
root      = { value = "${AWS_GODMODE_ROLE}", power_grade = 400 }
```

You do not need to use the table form unless you have a specific reason to deviate from the built-in power grades.

---

## [[accounts]]

Defines an AWS account. Repeat the `[[accounts]]` block for each account.

| Field            | Type     | Required | Description                                                                                                                                                                          |
|------------------|----------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`           | string   | yes      | Unique account name. Used in `yak account <name>`.                                                                                                                                   |
| `account_id`     | string   | yes      | 12-digit AWS account ID. Supports `${VAR}` and `op://`.                                                                                                                              |
| `default_role`   | string   | no       | Default role tier for this account. Overrides global `default_role`.                                                                                                                 |
| `region`         | string   | no       | AWS region for this account. Overrides global `region`.                                                                                                                              |
| `allowed_roles`  | string[] | no       | Permitted role tier names on this account. When set, must include `default_role`. Omitting this field allows all defined roles.                                                      |
| `role_overrides` | map      | no       | Per-account role name overrides. Keys are tier names from `[roles]`. Values replace the global role value for this account only. Power grades cannot be overridden at account level. |
| `safeguard`      | bool     | no       | Per-account safeguard override. Defaults to `true` if the account matches `[safeguards].protected_accounts` or `protected_pattern`.                                                  |

### Minimal account (inherits all global defaults)

```toml
[[accounts]]
name       = "dev"
account_id = "${AWS_DEV_ACCOUNT_ID}"
```

### Account with restrictions and role overrides

```toml
[[accounts]]
name           = "prod"
account_id     = "${AWS_PROD_ACCOUNT_ID}"
default_role   = "developer"
allowed_roles  = ["read_only", "developer", "admin"]  # root intentionally excluded
role_overrides = { admin = "${AWS_ADMIN_ROLE}", read_only = "${AWS_READ_ROLE}" }
```

### Role resolution order

When yak needs to resolve the role ARN for an account, it applies these rules in order:

1. Determine active role tier — from CLI flag, `yak role` state, account `default_role`, falling back to global `default_role`
2. Check `allowed_roles` — if defined, the active tier must be listed or yak exits with an error
3. Check safeguards — if the account is protected and the tier's power grade exceeds `max_power_grade_allowed`, yak blocks the operation (unless `--bypass-safeguards` is passed)
4. Check `role_overrides` — if the account has an override for the active tier, use that role value
5. Fall back to global `[roles]` — use the global value for this tier
6. Construct the ARN — `arn:aws:iam::{account_id}:role/{resolved_role_name}`

---

## [aliases.accounts]

Short aliases that map to account names. Defined once, applies globally.

```toml
[aliases.accounts]
d      = "dev"
sb     = "sandbox"
t      = "test"
p      = "prod"
m      = "mgmt"
```

Usage: `yak account p` is equivalent to `yak account prod`.

---

## [aliases.roles]

Short aliases that map to role tier names.

```toml
[aliases.roles]
ro  = "read_only"
dev = "developer"
adm = "admin"
rt  = "root"
```

Usage: `yak role adm` is equivalent to `yak role admin`.

---

## Full example

See the [examples section](examples/toml-env.md) for complete configs in all four format and secret provider combinations.
