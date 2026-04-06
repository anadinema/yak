# Migrating from Shell Functions

If you're coming from a set of shell functions like `awsi`, `awsp`, `awsr`, `awsc`, and `aws_setup`, this guide maps each to its yak equivalent.

## Command mapping

| Shell function | yak equivalent             | Notes                                                   |
|----------------|----------------------------|---------------------------------------------------------|
| `aws_setup`    | `yak setup`                | Config-driven, no env var template or `envsubst` needed |
| `awsi`         | `yak login`                | Same SSO browser flow + credential export               |
| `awsp dev`     | `yak account dev`          | Also accepts aliases: `yak account d`                   |
| `awsp prod`    | `yak account prod`         | Safeguard is active on prod by default                  |
| `awsr admin`   | `yak role admin`           | Switches role tier, regenerates config automatically    |
| `awsc dev`     | `yak export --account dev` | Exports credentials for one account                     |
| `awsc all`     | `yak export --all`         | Exports credentials for all accounts                    |

## What changes

**Config file instead of a template.**
You no longer need `config.template` + `envsubst`. yak reads account IDs and role names directly from `~/.config/yak/config.toml`, resolving `${ENV_VAR}` references at runtime. Your existing env vars (`$AWS_DEVELOPER_ROLE`, `$AWS_DEV_ACCOUNT_ID` etc.) map directly into the config.

**No more manual `aws_setup` after every role switch.**
`yak role` regenerates `~/.aws/config` automatically as part of the switch.

**`AWS_PROFILE` via shell wrapper.**
The shell wrapper replaces the `export AWS_PROFILE=...` lines that your shell functions handled directly. See [shell integration](shell-integration.md).

**Safeguards replace manual discipline.**
Previously avoiding `admin` on prod was a personal rule. yak enforces it automatically via `[safeguards]` — you have to explicitly bypass it with confirmation.

## Migration steps

1. Install yak — see [install](../install.md)
2. Create `~/.config/yak/config.toml` using the [TOML + env vars example](../config/examples/toml-env.md) as a starting point
3. Your existing env vars map directly — `$AWS_DEVELOPER_ROLE` becomes `developer = "${AWS_DEVELOPER_ROLE}"` in `[roles]`
4. Add the shell wrapper to your `.zshrc` — see [shell integration](shell-integration.md)
5. Run `yak setup` to verify config parses and generates `~/.aws/config` correctly
6. Run `yak login` to authenticate
7. Test `yak account dev`, `yak role admin`, `yak status`
8. Remove the old shell functions from your `.zshrc` once everything works

## Keeping old functions temporarily

You can keep both working during the transition by renaming the old functions:

```sh
# Temporarily renamed while testing yak
_legacy_awsp() { ... }
_legacy_awsr() { ... }
```

This way you can fall back if needed without losing anything.
