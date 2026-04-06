# yak

**Yet Another AWS Kit** — because AWS SSO should work like this out of the box.

yak wraps the AWS CLI to make account switching, SSO login, credential management, and role switching something you actually want to do. It adds safeguards for production accounts, integrates with 1Password, and gets out of your way for everything else.

---

## What it does

```sh
# Switch to your prod account with a single command
yak account prod

# Switch role tiers across all accounts
yak role admin

# Login and export credentials in one step
yak login

# See exactly where you are
yak status
```

## Why

Managing multiple AWS accounts via SSO involves too many steps: logging in, exporting credentials, switching profiles, regenerating config files when roles change. `yak` turns all of this into single short commands, with a config file that mirrors how your AWS org is actually structured.

If `yak` doesn't know a command, it [passes it through](guides/shell-integration.md) to the AWS CLI transparently — so you can alias `aws` to `yak` and nothing breaks.

---

## Quick install

=== "macOS (Homebrew)"

    ```sh
    brew tap anadinema/tap
    brew install yak
    ```

=== "Linux"

    ```sh
    curl -fsSL https://raw.githubusercontent.com/anadinema/yak/main/install.sh | bash
    ```

=== "Windows (Git Bash)"

    Download the `linux_amd64` binary from the [latest release](https://github.com/anadinema/yak/releases/latest) and place it in `~/bin`.

    See the full [Windows guide](guides/windows-git-bash.md).

---

[Get started →](install.md){ .md-button .md-button--primary }
[Configuration reference →](config/reference.md){ .md-button }
