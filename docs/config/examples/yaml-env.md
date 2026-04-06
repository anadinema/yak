# YAML + Environment Variables

Complete `~/.config/yak/config.yaml` using environment variables for all sensitive values.

```yaml title="~/.config/yak/config.yaml"
# ─── Global defaults ──────────────────────────────────────────────────────────
default_account: dev
default_role: admin
region: eu-north-1
sso_start_url: "${AWS_SSO_URL}"

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
  protected_accounts: [prod, mgmt]
  protected_pattern: ""
  max_power_grade_allowed: developer

# ─── Roles ────────────────────────────────────────────────────────────────────
roles:
  read_only: "${AWS_VIEWER_ROLE}"
  developer: "${AWS_DEVELOPER_ROLE}"
  admin: "${AWS_ELEVATED_ROLE}"
  root: "${AWS_GODMODE_ROLE}"

# ─── Accounts ─────────────────────────────────────────────────────────────────
accounts:
  - name: dev
    account_id: "${AWS_DEV_ACCOUNT_ID}"

  - name: sandbox
    account_id: "${AWS_SANDBOX_ACCOUNT_ID}"

  - name: test
    account_id: "${AWS_TEST_ACCOUNT_ID}"

  - name: prod
    account_id: "${AWS_PROD_ACCOUNT_ID}"
    default_role: developer
    allowed_roles: [read_only, developer, admin]
    role_overrides:
      admin: "${AWS_ADMIN_ROLE}"
      read_only: "${AWS_READ_ROLE}"

  - name: mgmt
    account_id: "${AWS_MGMT_ACCOUNT_ID}"
    default_role: admin
    allowed_roles: [read_only, developer, admin]
    role_overrides:
      admin: "${AWS_ADMIN_ROLE}"
      read_only: "${AWS_READ_ROLE}"

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

## Required environment variables

```sh
export AWS_SSO_URL="https://yourorg.awsapps.com/start"
export AWS_VIEWER_ROLE="ViewerRole"
export AWS_DEVELOPER_ROLE="DeveloperRole"
export AWS_ELEVATED_ROLE="PrivilegedRole"
export AWS_GODMODE_ROLE="GodModeRole"
export AWS_ADMIN_ROLE="AdministratorRole"
export AWS_READ_ROLE="ReaderRole"
export AWS_DEV_ACCOUNT_ID="123456789012"
export AWS_SANDBOX_ACCOUNT_ID="234567890123"
export AWS_TEST_ACCOUNT_ID="345678901234"
export AWS_PROD_ACCOUNT_ID="456789012345"
export AWS_SPT_ACCOUNT_ID="567890123456"
export AWS_SPP_ACCOUNT_ID="678901234567"
export AWS_MGMT_ACCOUNT_ID="789012345678"
```
