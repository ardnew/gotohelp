// gotohelp translates Go command lines into help commands.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func helpArgs(args []string) []string {
	clean := make([]string, 0, len(args))
	help := false
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "-h" || arg == "--help" {
			help = true
		} else {
			clean = append(clean, arg)
		}
	}
	if !help {
		return args
	}

	// -C must remain before the command so Go resolves module tools there.
	var prefix []string
	if len(clean) > 1 && clean[0] == "-C" {
		prefix, clean = clean[:2:2], clean[2:]
	} else if len(clean) > 0 && strings.HasPrefix(clean[0], "-C=") {
		prefix, clean = clean[:1:1], clean[1:]
	}
	if len(clean) == 0 {
		return append(prefix, "help")
	}
	command, rest := clean[0], clean[1:]
	if command == "help" {
		if len(rest) == 0 || strings.HasPrefix(rest[0], "-") {
			return append(prefix, "help")
		}
		result := append(prefix, "help", rest[0])
		if (rest[0] == "mod" || rest[0] == "work") && len(rest) > 1 && !strings.HasPrefix(rest[1], "-") {
			result = append(result, rest[1])
		}
		return result
	}
	if command == "tool" {
		toolFlags := []string{}
		for len(rest) > 0 && strings.HasPrefix(rest[0], "-") {
			arg := rest[0]
			rest = rest[1:]
			if arg == "-n" {
				continue
			}
			toolFlags = append(toolFlags, arg)
			if (arg == "-C" || arg == "-modfile" || arg == "-overlay") && len(rest) > 0 {
				toolFlags = append(toolFlags, rest[0])
				rest = rest[1:]
			}
		}
		if len(rest) == 0 {
			return append(prefix, "help", "tool")
		}
		result := append(prefix, "tool")
		result = append(result, toolFlags...)
		return append(result, rest[0], "-h")
	}
	result := append(prefix, "help", command)
	if command == "mod" || command == "work" {
		// These command groups accept -C before their subcommand.
		if len(rest) > 1 && rest[0] == "-C" {
			rest = rest[2:]
		} else if len(rest) > 0 && strings.HasPrefix(rest[0], "-C=") {
			rest = rest[1:]
		}
		if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
			result = append(result, rest[0])
		}
	}
	return result
}

func quote(s string) string {
	if s != "" && strings.IndexFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("_@%+=:,./-", r))
	}) == -1 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func main() {
	args := helpArgs(os.Args[1:])
	goexe := os.Getenv("GOTOHELPGO")
	if goexe != "" {
		_, lookErr := exec.LookPath(goexe)
		if _, err := os.Stat(goexe); err == nil || lookErr == nil {
			cmd := exec.Command(goexe, args...)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				var exit *exec.ExitError
				if errors.As(err, &exit) && exit.ExitCode() >= 0 {
					os.Exit(exit.ExitCode())
				}
				fmt.Fprintln(os.Stderr, "gotohelp:", err)
				os.Exit(1)
			}
			return
		}
	}
	words := append([]string{"go"}, args...)
	for i, word := range words {
		words[i] = quote(word)
	}
	fmt.Println(strings.Join(words, " "))
}
