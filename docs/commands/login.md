# yak login

SSO login and credential export.

## Usage

```sh
yak login [flags]
```

## Description

Runs `aws sso login` for the active account, then exports the resulting credentials to `~/.aws/credentials`. Always bypasses the 1Password secret cache and fetches fresh values.

## Flags

| Flag               | Description                                           |
|--------------------|-------------------------------------------------------|
| `--account <name>` | Login as a specific account instead of the active one |
| `--role <tier>`    | Use a specific role tier for this login               |

## Example

```sh
yak login
yak login --account prod
```

## Notes

After running `yak login`, your shell wrapper exports `AWS_PROFILE` automatically if configured. See [shell integration](../guides/shell-integration.md).
