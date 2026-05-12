# Commands

yak provides a focused set of commands for AWS account and credential management.

| Command                       | Description                                   |
|-------------------------------|-----------------------------------------------|
| [`yak generate`](generate.md) | Generate a minimal config file template       |
| [`yak setup`](setup.md)       | Generate `~/.aws/config` from your yak config |
| [`yak login`](login.md)       | SSO login and credential export               |
| [`yak account`](account.md)   | Switch active AWS account                     |
| [`yak role`](role.md)         | Switch active role tier                       |
| [`yak export`](export.md)     | Export credentials without switching          |
| [`yak status`](status.md)     | Show current account, role, and login state   |
| [`yak accounts`](accounts.md) | List all configured accounts                  |
| [`yak roles`](roles.md)       | List available role tiers                     |

## Global flags

These flags work with every command.

| Flag              | Description                                               |
|-------------------|-----------------------------------------------------------|
| `--help`          | Show help for the command                                 |
| `--version`       | Print yak version                                         |

## Global environment variables

| Variable          | Description                                                                    |
|-------------------|--------------------------------------------------------------------------------|
| `YAK_CONFIG_FILE` | Override config file path; when unset, defaults to `~/.config/yak/config.toml` |
