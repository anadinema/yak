# Safeguards

yak includes a safeguard system that prevents accidentally using elevated roles on production-like accounts.

## How it works

Each role tier has a built-in **power grade** — a numeric value yak uses to compare role privilege levels:

| Role        | Power grade |
|-------------|-------------|
| `read_only` | 100         |
| `developer` | 200         |
| `admin`     | 300         |
| `root`      | 400         |

When you switch role on a protected account, yak checks whether the requested role's power grade exceeds `max_power_grade_allowed`. If it does, the operation is blocked with a clear error.

## Configuration

```toml
[safeguards]
enabled                 = true
protected_accounts      = ["prod", "mgmt"]   # exact account names
protected_pattern       = ""                        # optional regex, e.g. "^prod-"
max_power_grade_allowed = "developer"               # power grade 200 ceiling
```

With this config:

- `yak role read_only --account prod` — allowed (100 ≤ 200)
- `yak role developer --account prod` — allowed (200 ≤ 200)
- `yak role admin --account prod` — **blocked** (300 > 200)
- `yak role root --account prod` — **blocked** (400 > 200)

`protected_accounts` matches exact names. `protected_pattern` is a regex applied additionally — any account whose name matches either is protected. Both can be used together.

## Bypassing safeguards

For genuine emergency access, use `--bypass-safeguards` on the `role` subcommand:

```sh
yak role root --account prod --bypass-safeguards
```

yak displays a prominent warning and requires you to type the account name to confirm:

```
⚠ WARNING: Safeguard bypass requested
  Account:         prod
  Requested role:  root (power grade 400)
  Safeguard limit: developer (power grade 200)

  Type the account name to confirm: _
```

Typing the wrong account name cancels the operation. All bypasses are logged with timestamp, account, role, and username to `~/.local/share/yak/audit.log`.

## Per-account safeguard override

To exclude a specific account from safeguards even if it matches a protected pattern:

```toml
[[accounts]]
name       = "prod-sandbox"
account_id = "${AWS_PROD_SANDBOX_ACCOUNT_ID}"
safeguard  = false   # not a real prod account, override the pattern match
```

## Disabling safeguards entirely

```toml
[safeguards]
enabled = false
```

This disables all safeguard checks globally. Not recommended for setups with real production accounts.

!!! danger
    Disabling safeguards removes all protection. If you need temporary elevated access, prefer `--bypass-safeguards` on a per-operation basis rather than disabling globally.
