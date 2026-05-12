# yak generate

Generate a minimal yak config template.

## Usage

```sh
yak generate [flags]
```

## Description

Creates a minimal config file with required sections only and placeholder empty strings for account IDs and role names.

Write target:

- `YAK_CONFIG_FILE` if set
- otherwise:
  - `~/.config/yak/config.toml` for `--format toml` (default)
  - `~/.config/yak/config.yaml` for `--format yaml`

The command fails if the target file already exists.

## Flags

| Flag              | Description                          |
|-------------------|--------------------------------------|
| `--format <type>` | Output format: `toml` or `yaml`      |

## Examples

```sh
yak generate
yak generate --format yaml
YAK_CONFIG_FILE=~/tmp/yak.dev.yaml yak generate --format yaml
```
