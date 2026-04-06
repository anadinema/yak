# Using with 1Password

yak resolves any config value prefixed with `op://` using the 1Password CLI at runtime. This means your config file contains no secrets and can be safely committed to your dotfiles repository.

## Prerequisites

1. Install the [1Password CLI](https://developer.1password.com/docs/cli/get-started/)
2. Sign in: `op signin`
3. Keep `op` on your `PATH`

## Value format

```
op://vault/item/field
```

```toml
account_id    = "op://Private/AWS Accounts/dev_account_id"
sso_start_url = "op://Private/AWS/sso_start_url"
```

## Caching

To avoid calling `op` on every yak command, resolved values are cached in `~/.local/share/yak/secret_cache` with permissions `0600`. The cache lifetime is controlled by `secrets.cache_ttl` in minutes.

```toml
[secrets]
cache_ttl = 60   # re-fetch from 1Password after 60 minutes
```

Set `cache_ttl = 0` to disable caching entirely. `yak login` always bypasses the cache regardless of this setting.

!!! warning "Cache security"
    The secret cache stores resolved values in plaintext, protected only by file permissions (`0600`). If this is a concern, set `cache_ttl = 0`. Values will be fetched fresh from 1Password on every yak invocation.

## Sharing config with team members

With 1Password your entire config file is shareable — no secrets are embedded. Team members need access to the same vault and the `op` CLI installed.

Suggested workflow:

1. Commit `config.toml` (with `op://` references) to a shared dotfiles or team repo
2. Each person symlinks it: `ln -s ~/team-config/yak/config.toml ~/.config/yak/config.toml`
3. Each person authenticates: `op signin`
4. Run `yak setup`

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

## Mixing secret providers

You can mix `op://`, `${VAR}`, and plain text freely within a single config file:

```toml
region        = "eu-west-1"                          # plain text
sso_start_url = "${AWS_SSO_URL}"                     # env var
account_id    = "op://Private/AWS Accounts/dev_id"      # 1Password
```

yak detects the format automatically for each value.
