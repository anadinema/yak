# Windows (Git Bash)

yak works on Windows inside Git Bash using the Linux `amd64` binary. This avoids Windows terminal compatibility issues with interactive TUI rendering.

## Install

1. Download `yak_linux_amd64` from the [latest release](https://github.com/anadinema/yak/releases/latest)
2. Place it in `~/bin` — Git Bash adds this to `PATH` automatically
3. Rename and make executable:

```sh
mkdir -p ~/bin
mv ~/Downloads/yak_linux_amd64 ~/bin/yak
chmod +x ~/bin/yak
```

4. Verify:

```sh
yak --version
```

## Shell integration

Add the wrapper to `~/.bashrc` (Git Bash uses bash):

```sh
yak() {
  command yak "$@"
  if [[ "$1" == "account" || "$1" == "login" ]]; then
    export AWS_PROFILE="$(command yak status --print-profile)"
  fi
}
```

## AWS CLI

The Windows AWS CLI binary is accessible from Git Bash — `aws --version` should work without any extra setup.

## SSO login

`yak login` calls `aws sso login` which opens your browser. This works correctly in Git Bash — the AWS CLI handles the browser launch via Windows native APIs regardless of the terminal.

## VS Code integration

If you use VS Code with Git Bash as your default terminal, no additional setup is needed. yak works in the integrated terminal the same as in a standalone Git Bash window.

!!! tip
    Keep Git Bash as your VS Code default terminal — yak works there and your npm scripts continue to work as expected.
