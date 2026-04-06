# Install

## Requirements

- AWS CLI v2 installed and on your `PATH`
- For 1Password secret resolution: `op` CLI installed and signed in
- macOS 12+, Linux (amd64/arm64), or Windows via Git Bash

---

## macOS

```sh
brew tap anadinema/tap
brew install yak
```

To upgrade:

```sh
brew upgrade yak
```

---

## Linux

### Install script

```sh
curl -fsSL https://raw.githubusercontent.com/anadinema/yak/main/install.sh | bash
```

The script detects your architecture and installs the binary to `~/.local/bin`. Make sure `~/.local/bin` is on your `PATH`.

### Manual

Download the binary for your architecture from the [releases page](https://github.com/anadinema/yak/releases/latest), make it executable, and move it to somewhere on your `PATH`.

```sh
# Example for linux/amd64
curl -Lo yak https://github.com/anadinema/yak/releases/latest/download/yak_linux_amd64
chmod +x yak
mv yak ~/.local/bin/yak
```

---

## Windows (Git Bash)

yak works best inside Git Bash using the `linux_amd64` binary. This avoids Windows terminal compatibility issues with the interactive TUI.

1. Download `yak_linux_amd64` from the [latest release](https://github.com/anadinema/yak/releases/latest)
2. Place it in `~/bin` (Git Bash adds this to `PATH` automatically)
3. Rename it to `yak` with no extension

```sh
# Inside Git Bash
mkdir -p ~/bin
mv ~/Downloads/yak_linux_amd64 ~/bin/yak
chmod +x ~/bin/yak
```

See the [Windows guide](../guides/windows-git-bash.md) for full setup including shell integration.

---

## Verify

```sh
yak --version
```

---

## Next steps

- [Quickstart](quickstart.md) — get up and running in five minutes
- [Configuration reference](../config/reference.md) — full config file documentation
- [Shell integration](../guides/shell-integration.md) — set up `AWS_PROFILE` export and optional passthrough
