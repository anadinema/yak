# yak role

Switch the active role tier.

## Usage

```sh
yak role <tier> [flags]
```

## Description

Switches the active role tier across all accounts, then regenerates `~/.aws/config`. Accepts full role tier names or any configured alias.

If the target account is protected by safeguards and the requested tier exceeds `max_power_grade_allowed`, yak blocks the operation with a clear error.

## Flags

| Flag                  | Description                                                                                                                           |
|-----------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `--account <name>`    | Apply the role switch to a specific account only                                                                                      |
| `--bypass-safeguards` | Override safeguard restrictions for this operation. Requires confirmation. All bypasses are logged to `~/.local/share/yak/audit.log`. |

## Examples

```sh
yak role developer
yak role adm              # alias for admin
yak role read_only --account prod
yak role root --account prod --bypass-safeguards
```

## Bypass confirmation

When `--bypass-safeguards` is used on a protected account, yak displays a warning and requires you to type the account name to confirm:

```
⚠ WARNING: Safeguard bypass requested
  Account:          prod
  Requested role:   root (power grade 400)
  Safeguard limit:  developer (power grade 200)

  Type the account name to confirm: _
```

This action is logged with timestamp to `~/.local/share/yak/audit.log`.
