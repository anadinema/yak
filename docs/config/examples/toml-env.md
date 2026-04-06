# TOML + Environment Variables

Complete `~/.config/yak/config.toml` using environment variables for all sensitive values.

This is the recommended starting point. All secrets come from your shell environment — the same place your existing dotfiles already manage them.

```toml title="~/.config/yak/config.toml"
# ─── Global defaults ──────────────────────────────────────────────────────────
default_account = "dev"
default_role    = "admin"
region          = "eu-north-1"
sso_start_url   = "${AWS_SSO_URL}"

[aws]
credentials_path = "~/.aws/credentials"
config_path      = "~/.aws/config"

[secrets]
cache_ttl = 60    # minutes, 0 = no cache

[passthrough]
enabled         = false
cli_binary_path = ""

[safeguards]
enabled                 = true
protected_accounts      = ["prod", "mgmt"]
protected_pattern       = ""
max_power_grade_allowed = "developer"

# ─── Roles ────────────────────────────────────────────────────────────────────
# Built-in power grades: read_only=100, developer=200, admin=300, root=400
[roles]
read_only = "${AWS_VIEWER_ROLE}"
developer = "${AWS_DEVELOPER_ROLE}"
admin     = "${AWS_ELEVATED_ROLE}"
root      = "${AWS_GODMODE_ROLE}"

# ─── Accounts ─────────────────────────────────────────────────────────────────
[[accounts]]
name       = "dev"
account_id = "${AWS_DEV_ACCOUNT_ID}"

[[accounts]]
name       = "sandbox"
account_id = "${AWS_SANDBOX_ACCOUNT_ID}"

[[accounts]]
name       = "test"
account_id = "${AWS_TEST_ACCOUNT_ID}"

[[accounts]]
name           = "prod"
account_id     = "${AWS_PROD_ACCOUNT_ID}"
default_role   = "developer"
allowed_roles  = ["read_only", "developer", "admin"]  # root intentionally excluded
role_overrides = { admin = "${AWS_ADMIN_ROLE}", read_only = "${AWS_READ_ROLE}" }

[[accounts]]
name           = "mgmt"
account_id     = "${AWS_MGMT_ACCOUNT_ID}"
default_role   = "admin"
allowed_roles  = ["read_only", "developer", "admin"]  # root intentionally excluded
role_overrides = { admin = "${AWS_ADMIN_ROLE}", read_only = "${AWS_READ_ROLE}" }

# ─── Aliases ──────────────────────────────────────────────────────────────────
[aliases.accounts]
d      = "dev"
sb     = "sandbox"
t      = "test"
p      = "prod"
m      = "mgmt"

[aliases.roles]
ro  = "read_only"
dev = "developer"
adm = "admin"
rt  = "root"
```

## Required environment variables

Add these to your `.zshrc`, `.bashrc`, or `auth-env.shrc`:

```sh
export AWS_SSO_URL="https://yourorg.awsapps.com/start"

# Role names
export AWS_VIEWER_ROLE="ViewerRole"
export AWS_DEVELOPER_ROLE="DeveloperRole"
export AWS_ELEVATED_ROLE="PrivilegedRole"
export AWS_GODMODE_ROLE="GodModeRole"

# Role overrides for prod-like accounts
export AWS_ADMIN_ROLE="AdministratorRole"
export AWS_READ_ROLE="ReaderRole"

# Account IDs
export AWS_DEV_ACCOUNT_ID="123456789012"
export AWS_SANDBOX_ACCOUNT_ID="234567890123"
export AWS_TEST_ACCOUNT_ID="345678901234"
export AWS_PROD_ACCOUNT_ID="456789012345"
export AWS_MGMT_ACCOUNT_ID="789012345678"
```
