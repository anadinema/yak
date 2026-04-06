# Shell Integration

## AWS_PROFILE export

A Go binary cannot set environment variables in the parent shell — this is an OS constraint, not a yak limitation. To have `AWS_PROFILE` updated in your terminal after `yak account` or `yak login`, add this shell wrapper to your `.zshrc` or `.bashrc`:

```sh
yak() {
  command yak "$@"
  if [[ "$1" == "account" || "$1" == "login" ]]; then
    export AWS_PROFILE="$(command yak status --print-profile)"
  fi
}
```

Reload:

```sh
source ~/.zshrc   # zsh
source ~/.bashrc  # bash
```

This pattern is used by `rbenv`, `nvm`, and other tools that need to influence the shell environment.

---

## Passthrough

When `passthrough.enabled = true` in your config, any command yak does not recognise is forwarded verbatim to the AWS CLI. Exit codes, stdout, and stderr pass through unchanged.

```toml
[passthrough]
enabled         = true
cli_binary_path = ""   # auto-detected from PATH
```

With passthrough enabled you can alias `aws` to `yak`:

```sh
alias aws=yak
```

Now `aws s3 ls`, `aws ec2 describe-instances`, and everything else continues to work exactly as before — while `yak account`, `yak login`, and `yak role` give you the extra functionality.

!!! note
    Passthrough is disabled by default. It is an explicit opt-in because forwarding unknown commands may have unintended consequences if yak is misconfigured.

---

## Optional: short aliases

```sh
alias ya='yak account'
alias yr='yak role'
alias yl='yak login'
alias ys='yak status'
```
