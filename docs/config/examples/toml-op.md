# TOML + 1Password

Complete `~/.config/yak/config.toml` using 1Password for all sensitive values.

This setup means the config file itself can be safely committed to your dotfiles repository — no secrets are stored in plaintext or environment variables.

## Prerequisites

- [1Password CLI (`op`)](https://developer.1password.com/docs/cli/) installed
- Signed in to 1Password: `op signin`

```toml title="~/.config/yak/config.toml"
# ─── Global defaults ──────────────────────────────────────────────────────────
default_account = "dev"
default_role    = "admin"
region          = "eu-north-1"
sso_start_url   = "op://Private/AWS/sso_start_url"

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
[roles]
read_only = "op://Private/AWS Roles/viewer_role"
developer = "op://Private/AWS Roles/developer_role"
admin     = "op://Private/AWS Roles/privileged_role"
root      = "op://Private/AWS Roles/godmode_role"

# ─── Accounts ─────────────────────────────────────────────────────────────────
[[accounts]]
name       = "dev"
account_id = "op://Private/AWS Accounts/dev_account_id"

[[accounts]]
name       = "sandbox"
account_id = "op://Private/AWS Accounts/sandbox_account_id"

[[accounts]]
name       = "test"
account_id = "op://Private/AWS Accounts/test_account_id"

[[accounts]]
name           = "prod"
account_id     = "op://Private/AWS Accounts/prod_account_id"
default_role   = "developer"
allowed_roles  = ["read_only", "developer", "admin"]
role_overrides = { admin = "op://Private/AWS Roles/prod_admin_role", read_only = "op://Private/AWS Roles/prod_read_role" }

[[accounts]]
name           = "mgmt"
account_id     = "op://Private/AWS Accounts/mgmt_account_id"
default_role   = "admin"
allowed_roles  = ["read_only", "developer", "admin"]
role_overrides = { admin = "op://Private/AWS Roles/prod_admin_role", read_only = "op://Private/AWS Roles/prod_read_role" }

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

## Suggested 1Password vault structure

```
Work/
  AWS/
    sso_start_url
  AWS Roles/
    viewer_role
    developer_role
    privileged_role
    godmode_role
    prod_admin_role
    prod_read_role
  AWS Accounts/
    dev_account_id
    sandbox_account_id
    test_account_id
    prod_account_id
    spt_account_id
    spp_account_id
    mgmt_account_id
```

## Sharing with team members

Since the config file contains no secrets, you can commit it to a shared repository. Team members only need to:

1. Install `op` CLI and sign in to the shared 1Password vault
2. Clone the repo and symlink or copy the config to `~/.config/yak/config.toml`
3. Run `yak setup`

See the [1Password guide](../../guides/onepassword.md) for more detail.
