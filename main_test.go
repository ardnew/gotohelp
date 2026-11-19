package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestHelpArgs(t *testing.T) {
	tests := []struct{ in, want []string }{
		{nil, nil},
		{[]string{"build", "-o", "foo"}, []string{"build", "-o", "foo"}},
		{[]string{"-h"}, []string{"help"}},
		{[]string{"build", "-o", "foo", "-tags", "bar", "-h"}, []string{"help", "build"}},
		{[]string{"--help", "test", "-run", "TestX"}, []string{"help", "test"}},
		{[]string{"run", "main.go", "--", "-h", "--help"}, []string{"run", "main.go", "--", "-h", "--help"}},
		{[]string{"--", "-h"}, []string{"--", "-h"}},
		{[]string{"build", "-h", "--", "--help"}, []string{"help", "build"}},
		{[]string{"tool", "-h", "--", "compile"}, []string{"help", "tool"}},
		{[]string{"-C", "/tmp", "mod", "edit", "-h"}, []string{"-C", "/tmp", "help", "mod", "edit"}},
		{[]string{"work", "-C", "/tmp", "sync", "--help"}, []string{"help", "work", "sync"}},
		{[]string{"help", "mod", "edit", "-x", "-h"}, []string{"help", "mod", "edit"}},
		{[]string{"tool", "compile", "-o", "foo", "--help"}, []string{"tool", "compile", "-h"}},
		{[]string{"tool", "-n", "-C", "/tmp", "-modfile=alt.mod", "example.com/tool", "-h"}, []string{"tool", "-C", "/tmp", "-modfile=alt.mod", "example.com/tool", "-h"}},
		{[]string{"tool", "-h"}, []string{"help", "tool"}},
		{[]string{"doc", "fmt", "-h"}, []string{"help", "doc"}},
		{[]string{"vet", "-h"}, []string{"help", "vet"}},
		{[]string{"tool", "vet", "-h"}, []string{"tool", "vet", "-h"}},
	}
	for _, tt := range tests {
		if got := helpArgs(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("helpArgs(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestQuote(t *testing.T) {
	for input, want := range map[string]string{"go": "go", "": "''", "a b": "'a b'", "a'b": "'a'\\''b'", "$(x)": "'$(x)'"} {
		if got := quote(input); got != want {
			t.Errorf("quote(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestGoExecutable(t *testing.T) {
	if os.Getenv("GOTOHELP_TEST_CHILD") == "1" {
		os.Args = []string{"gotohelp", "build", "-h"}
		main()
		os.Exit(0)
	}
	goexe, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	relative, err := filepath.Rel(cwd, goexe)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.ContainsRune(relative, os.PathSeparator) {
		relative = "." + string(os.PathSeparator) + relative
	}
	for _, tc := range []struct {
		name, executable, want string
	}{
		{"unset", "", "go help build\n"},
		{"path", "go", "usage: go build"},
		{"relative", relative, "usage: go build"},
		{"absolute", goexe, "usage: go build"},
		{"missing", filepath.Join(t.TempDir(), "missing"), "go help build\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestGoExecutable$")
			cmd.Env = append(os.Environ(), "GOTOHELP_TEST_CHILD=1", "GOTOHELPGO="+tc.executable)
			out, err := cmd.CombinedOutput()
			if err != nil || !strings.Contains(string(out), tc.want) {
				t.Fatalf("output = %q, error = %v; want %q", out, err, tc.want)
			}
			if (tc.name == "unset" || tc.name == "missing") && string(out) != tc.want {
				t.Fatalf("output = %q; want %q", out, tc.want)
			}
		})
	}
}
