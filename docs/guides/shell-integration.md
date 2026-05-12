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

## Optional: short aliases

```sh
alias ya='yak account'
alias yr='yak role'
alias yl='yak login'
alias ys='yak status'
```
