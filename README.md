# passgen

Simple and secure password generator CLI tool.

Generates strong, URL-safe passwords with ambiguous characters excluded by default. Copies to clipboard automatically.

## Features

- **Secure by default** — Uses `crypto/rand` for cryptographically secure random generation
- **Simple CLI** — Minimal flags, sensible defaults, zero configuration files
- **URL-safe characters** — Symbol set limited to `-_.~` (RFC 3986 unreserved)
- **Ambiguous characters excluded** — No `0O`, `1lI` confusion
- **Clipboard integration** — Auto-copies generated password to clipboard

## Installation

### Homebrew

```bash
brew install youyo/tap/passgen
```

### go install

```bash
go install github.com/youyo/passgen@latest
```

### GitHub Releases

Download pre-built binaries from the [Releases](https://github.com/youyo/passgen/releases) page.

## Usage

```bash
# Generate a 20-character password (default)
passgen

# Specify length
passgen 32

# Require at least 3 symbols and 2 digits
passgen --symbols 3 --digits 2

# Use short flags
passgen -s 3 -d 2 -u 2 -l 2

# Disable clipboard copy
passgen --no-copy

# Suppress stdout output (clipboard only)
passgen --no-print

# Exclude specific characters
passgen --exclude "abc123"

# Combine options
passgen 30 --symbols 2 --no-copy
```

### Environment Variables

Override defaults with environment variables. Priority: **CLI flag > Environment variable > Default**.

```bash
export PASSGEN_LENGTH=30
export PASSGEN_SYMBOLS=2
passgen  # uses length=30, symbols=2
```

## Character Set

All character sets exclude ambiguous characters to avoid confusion when reading passwords:

| Category | Characters | Count | Excluded |
|----------|-----------|-------|----------|
| Lower    | `abcdefghijkmnopqrstuvwxyz` | 25 | `l` |
| Upper    | `ABCDEFGHJKLMNPQRSTUVWXYZ` | 24 | `I`, `O` |
| Digits   | `23456789` | 8 | `0`, `1` |
| Symbols  | `-_.~` | 4 | — |

**Total: 61 characters**

## Flags

| Flag | Short | Default | Env | Description |
|------|-------|---------|-----|-------------|
| `[length]` | — | `20` | `PASSGEN_LENGTH` | Password length (positional argument) |
| `--symbols` | `-s` | `1` | `PASSGEN_SYMBOLS` | Minimum number of symbols |
| `--digits` | `-d` | `1` | `PASSGEN_DIGITS` | Minimum number of digits |
| `--upper` | `-u` | `1` | `PASSGEN_UPPER` | Minimum number of uppercase letters |
| `--lower` | `-l` | `1` | `PASSGEN_LOWER` | Minimum number of lowercase letters |
| `--exclude` | `-e` | `""` | `PASSGEN_EXCLUDE` | Characters to exclude |
| `--no-copy` | — | `false` | — | Disable clipboard copy |
| `--no-print` | — | `false` | — | Disable stdout output |
| `--version` | — | — | — | Show version information |

> **Note:** `--no-copy` and `--no-print` cannot be used together.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PASSGEN_LENGTH` | Password length | `20` |
| `PASSGEN_SYMBOLS` | Minimum number of symbols | `1` |
| `PASSGEN_DIGITS` | Minimum number of digits | `1` |
| `PASSGEN_UPPER` | Minimum number of uppercase letters | `1` |
| `PASSGEN_LOWER` | Minimum number of lowercase letters | `1` |
| `PASSGEN_EXCLUDE` | Characters to exclude | `""` |

## Shell Completion

### zsh

Add the following to your `.zshrc`:

```bash
eval "$(passgen completion zsh --short)"
```

Or generate the full completion script:

```bash
passgen completion zsh > /usr/local/share/zsh/site-functions/_passgen
```

## License

[MIT](LICENSE)
