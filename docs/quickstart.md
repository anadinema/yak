# Quickstart

Get yak working in about five minutes.

---

## 1. Create your config file

```sh
mkdir -p ~/.config/yak
```

Create `~/.config/yak/config.toml`. Start minimal — you can add accounts and roles as you go.

```toml
default_account = "dev"
default_role    = "developer"
region          = "eu-west-1"
sso_start_url   = "${AWS_SSO_URL}"

[aws]
credentials_path = "~/.aws/credentials"
config_path      = "~/.aws/config"

[secrets]
cache_ttl = 60

[roles]
read_only = "${AWS_VIEWER_ROLE}"
developer = "${AWS_DEVELOPER_ROLE}"
admin     = "${AWS_ELEVATED_ROLE}"
root      = "${AWS_GODMODE_ROLE}"

[[accounts]]
name       = "dev"
account_id = "${AWS_DEV_ACCOUNT_ID}"

[[accounts]]
name       = "prod"
account_id    = "${AWS_PROD_ACCOUNT_ID}"
default_role  = "developer"
allowed_roles = ["read_only", "developer", "admin"]

[aliases.accounts]
d = "dev"
p = "prod"

[aliases.roles]
ro  = "read_only"
dev = "developer"
adm = "admin"
```

---

## 2. Set your environment variables

Add these to your `.zshrc`, `.bashrc`, or `auth-env.shrc`:

```sh
export AWS_SSO_URL="https://yourorg.awsapps.com/start"
export AWS_DEV_ACCOUNT_ID="123456789012"
export AWS_PROD_ACCOUNT_ID="987654321098"
export AWS_DEVELOPER_ROLE="DeveloperRole"
export AWS_VIEWER_ROLE="ViewerRole"
export AWS_ELEVATED_ROLE="PrivilegedRole"
export AWS_GODMODE_ROLE="GodModeRole"
```

Or use [1Password](../guides/onepassword.md) — replace any value with `op://vault/item/field`.

---

## 3. Generate your AWS config

```sh
yak setup
```

This reads your `config.toml`, resolves all values, and writes `~/.aws/config`.

---

## 4. Add shell integration

Add this to your `.zshrc` or `.bashrc` so `AWS_PROFILE` is exported after account switches:

```sh
yak() {
  command yak "$@"
  if [[ "$1" == "account" || "$1" == "login" ]]; then
    export AWS_PROFILE="$(command yak status --print-profile)"
  fi
}
```

Reload your shell:

```sh
source ~/.zshrc  # or source ~/.bashrc
```

---

## 5. Login

```sh
yak login
```

This opens your browser for SSO authentication, then exports your credentials.

---

## 6. Switch accounts

```sh
yak account prod   # switch to prod
yak account d      # switch to dev (alias)
yak status         # see where you are
```

That's it. See the [command reference](../commands/index.md) for everything else.
