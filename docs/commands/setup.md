# yak setup

Generate `~/.aws/config` from your yak configuration.

## Usage

```sh
yak setup [flags]
```

## Description

Reads `~/.config/yak/config.toml`, resolves all secret values, and writes a valid `~/.aws/config` with one `[profile name]` block per configured account. The active role tier from state is used to determine which role goes into each profile block.

Run `yak setup` after:

- First-time setup
- Adding or removing accounts
- Changing role names in config
- Changing `sso_start_url`

## Flags

| Flag        | Description                                            |
|-------------|--------------------------------------------------------|
| `--dry-run` | Print what would be written without touching any files |

## Example

```sh
yak setup
yak setup --dry-run
```

## Generated output

For each account, yak writes a profile block like:

```ini
[profile dev]
sso_start_url  = https://yourorg.awsapps.com/start
sso_region     = eu-west-1
sso_account_id = 123456789012
sso_role_name  = PrivilegedRole
region         = eu-west-1
```
