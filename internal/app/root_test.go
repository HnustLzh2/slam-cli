package app

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestHelpStaysConciseAndRoutesToCommandHelp(t *testing.T) {
	cmd := NewHelpCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("help command failed: %v", err)
	}

	help := out.String()
	required := []string{
		"slam-cli package update <name> --file update-package.json",
		"slam-cli xflow create --file create-xflow.json",
		"复杂命令可执行对应子命令 --help 查看 JSON 示例",
	}
	for _, want := range required {
		if !strings.Contains(help, want) {
			t.Fatalf("help output missing %q\n%s", want, help)
		}
	}

	forbidden := []string{
		`"publishOpts"`,
		`"registerOpts"`,
		`"syncOpts"`,
		`"reviewDoc"`,
	}
	for _, notWant := range forbidden {
		if strings.Contains(help, notWant) {
			t.Fatalf("root help should stay concise, but contains %q\n%s", notWant, help)
		}
	}
}

func TestHelpOutputMatchesTextFile(t *testing.T) {
	want, err := os.ReadFile("help.txt")
	if err != nil {
		t.Fatalf("read help text file: %v", err)
	}

	cmd := NewHelpCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("help command failed: %v", err)
	}

	if out.String() != string(want) {
		t.Fatalf("help output should match help.txt\nwant:\n%s\ngot:\n%s", string(want), out.String())
	}
}
