# yak export

Export credentials to `~/.aws/credentials` without switching accounts or regenerating config.

## Usage

```sh
yak export [flags]
```

## Description

Calls `aws configure export-credentials` for one or all accounts and writes the results to your credentials file. Useful for refreshing credentials that have expired without a full login.

## Flags

| Flag               | Description                                                         |
|--------------------|---------------------------------------------------------------------|
| `--account <name>` | Export credentials for a specific account (default: active account) |
| `--all`            | Export credentials for all configured accounts                      |

## Examples

```sh
yak export
yak export --account prod
yak export --all
```
