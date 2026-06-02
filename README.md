# binpath

A small CLI utility to query and manipulate your `PATH` environment variable: add entries, 
remove them (regex supported), and check if it already contains what you ask for.  
Mainly for use in scripts and shell configs. Linux and macOS only. (don't have the use for 
window and unable to tests)

## How it works

A process cannot change its parent shell's `PATH`, so `binpath add`/`remove` print
the **new** `PATH` to stdout. Apply it by assigning the output back:

```sh
export PATH="$(binpath remove /opt/mytool)"
```

`binpath contains` prints nothing and reports the result through its exit code, so
it drops into shell conditionals like `test(1)`:

```sh
if binpath contains /usr/bin; then echo "present"; fi

# Or following `command` examples, you also can just redirect the uneeded output

if binpath contains /usr/bin &> /dev/null; then echo "present"; fi

```

## Install

### Install script (prebuilt binary)

```sh
curl -fsSL https://raw.githubusercontent.com/loosestack-dev/binpathutil/main/install.sh | sh
```

Installs the latest release to `~/.local/bin` (you will see a warning if this path is not already in your `PATH` (ironically, not using this utility, using the plain old syntax, which remind me exactly why I wanted something simpler)). Options:

```sh
# a specific release tag
curl -fsSL https://raw.githubusercontent.com/loosestack-dev/binpathutil/main/install.sh | sh -s -- v1.0.0

# a different install directory
BINPATH_INSTALL_DIR=/usr/local/bin curl -fsSL https://raw.githubusercontent.com/loosestack-dev/binpathutil/main/install.sh | sh
```

### With Go

```sh
go install github.com/loosestack-dev/binpathutil/cmd/binpath@latest
```

### From source

This project use `Taskfile` instead of `Makefile` because I still attach importance to my sanity. If you don't have it, I can only recommend to install it [using their website](https://taskfile.dev/). If you don't want to, you can always rely on the `go build` manually 

```sh
git clone https://github.com/loosestack-dev/binpathutil
cd binpathutil
task build      # or: go build -o binpath ./cmd/binpath
```

## Usage

```sh
binpath add <entry>        [--first|--last] [--if-absent]   # default: --first (prepend)
binpath remove <entry>     [--if-present] [--all] [--regex]
binpath contains <entry>   [--regex]
binpath --version
```

Examples:

```sh
export PATH="$(binpath add /opt/bin)"            # prepend
export PATH="$(binpath add --last /opt/bin)"     # append
export PATH="$(binpath add -i /opt/bin)"         # add only if absent

export PATH="$(binpath remove /opt/bin)"         # remove first occurrence (error if absent)
export PATH="$(binpath remove -i /opt/bin)"      # no error if absent
export PATH="$(binpath remove -a /opt/bin)"      # remove every occurrence
export PATH="$(binpath remove -r '/python[0-9.]*/')"      # remove first match of a regex
export PATH="$(binpath remove -r -a '/\.cache/')"          # remove all regex matches

binpath contains /usr/bin                        # exit 0 if present
binpath contains -r 'go/bin$'                    # regex membership test
```

Run `binpath <command> --help` for the full flag reference.

## Exit codes

`contains` (and strict `remove`) follow the `grep`/`test` convention:

| Code | Meaning                                        |
|------|------------------------------------------------|
| `0`  | present / success                              |
| `1`  | absent (the queried condition was negative)    |
| `2`  | error (PATH unset, malformed regex, bad flags) |

## Development

```sh
task test          # go test ./...
task test:count    # count passing tests/subtests
task build         # build the binary
```

## License

[0BSD](LICENSE) (Zero-Clause BSD) — do whatever you want with it.
