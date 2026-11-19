# gotohelp

I love the ergonomics of `go help`:

![go help](https://vhs.charm.sh/vhs-3Sau4u5QB6OVnmaOrj19MI.gif)

`gotohelp` skips this charade and shows the help. Give it the arguments you would give to `go` and add `-h` or `--help` before any `--` argument.

```sh
gotohelp build -o foo -tags bar -h  # prints: go help build
gotohelp mod edit --help            # prints: go help mod edit
gotohelp tool compile -h            # prints: go tool compile -h
```

Install with `go install github.com/ardnew/gotohelp@latest`. To build from the source, run `go build .`.

## How it works

By default, `gotohelp` prints the command. It does not run the command. The printed command starts with `go`. Arguments are quoted for the shell when necessary.

To run the command, set `GOTOHELPGO` to the Go executable. You can use an absolute path, a relative path, or a name in `PATH`:

```sh
GOTOHELPGO=/path/to/go gotohelp build -h  # runs: /path/to/go help build
GOTOHELPGO=go gotohelp build -h           # runs: go help build
```

The command uses the same input, output, and error streams as `gotohelp`. Its exit status is also passed back. If `GOTOHELPGO` cannot be found, `gotohelp` prints the command instead, starting with `go`.

Without `-h` or `--help` before `--`, `gotohelp` prints or runs the original arguments without changes.

## Use as wrapper for the `go` command

Add this function to your shell configuration to run `gotohelp` when you call `go`:

```sh
go() { GOTOHELPGO=go gotohelp "$@" ; }
```

For example, `go build -h` runs `go help build`. The function does not call itself: `gotohelp` finds the Go executable in `PATH`.

## Help requests

- `-h` and `--help` before `--` request help, even as values of other flags. Arguments after `--` do not request help.
- For Go commands, `gotohelp` removes other arguments and uses `go help`. It also recognizes `go mod` and `go work` subcommands.
- For `go tool`, it asks the named tool for help with `-h`. It keeps the tool context flags but removes `-n`, so the tool can show its help. The tool controls its help text and exit status.
- Before the tool name, use the documented `go tool` flags. An unknown flag with a separate value can be mistaken for the tool name.
- The global `-C` flag is preserved. Shell aliases and functions are not interpreted.
