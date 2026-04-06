# yak account

Switch the active AWS account.

## Usage

```sh
yak account <name> [flags]
```

## Description

Switches the active account, regenerates `~/.aws/config` with the new active account context, and exports credentials for the new account. Accepts full account names or any configured alias.

Updates `~/.local/share/yak/state.json` with the new active account.

## Flags

| Flag            | Description                                         |
|-----------------|-----------------------------------------------------|
| `--role <tier>` | Also switch to this role tier when changing account |

## Examples

```sh
yak account dev
yak account prod
yak account p         # alias for prod
yak account p --role read_only
```

## Notes

`yak account` cannot set `AWS_PROFILE` in your parent shell on its own — this requires the [shell wrapper](../guides/shell-integration.md).
