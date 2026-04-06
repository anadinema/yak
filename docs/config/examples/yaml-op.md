# YAML + 1Password

Complete `~/.config/yak/config.yaml` using 1Password for all sensitive values.

```yaml title="~/.config/yak/config.yaml"
# ─── Global defaults ──────────────────────────────────────────────────────────
default_account: dev
default_role: admin
region: eu-north-1
sso_start_url: "op://Private/AWS/sso_start_url"

aws:
  credentials_path: ~/.aws/credentials
  config_path: ~/.aws/config

secrets:
  cache_ttl: 60

passthrough:
  enabled: false
  cli_binary_path: ""

safeguards:
  enabled: true
  protected_accounts: [prod, spp, spt]
  protected_pattern: ""
  max_power_grade_allowed: developer

# ─── Roles ────────────────────────────────────────────────────────────────────
roles:
  read_only: "op://Private/AWS Roles/viewer_role"
  developer: "op://Private/AWS Roles/developer_role"
  admin: "op://Private/AWS Roles/privileged_role"
  root: "op://Private/AWS Roles/godmode_role"

# ─── Accounts ─────────────────────────────────────────────────────────────────
accounts:
  - name: dev
    account_id: "op://Private/AWS Accounts/dev_account_id"

  - name: sandbox
    account_id: "op://Private/AWS Accounts/sandbox_account_id"

  - name: test
    account_id: "op://Private/AWS Accounts/test_account_id"

  - name: prod
    account_id: "op://Private/AWS Accounts/prod_account_id"
    default_role: developer
    allowed_roles: [read_only, developer, admin]
    role_overrides:
      admin: "op://Private/AWS Roles/prod_admin_role"
      read_only: "op://Private/AWS Roles/prod_read_role"

  - name: mgmt
    account_id: "op://Private/AWS Accounts/mgmt_account_id"
    default_role: admin
    allowed_roles: [read_only, developer, admin]
    role_overrides:
      admin: "op://Private/AWS Roles/prod_admin_role"
      read_only: "op://Private/AWS Roles/prod_read_role"

# ─── Aliases ──────────────────────────────────────────────────────────────────
aliases:
  accounts:
    d: dev
    sb: sandbox
    t: test
    p: prod
    m: mgmt
  roles:
    ro: read_only
    dev: developer
    adm: admin
    rt: root
```
