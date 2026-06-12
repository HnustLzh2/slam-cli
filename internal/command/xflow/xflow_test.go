package xflowCmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCreateHelpIncludesCreateXflowJSONGuidance(t *testing.T) {
	cmd := NewXflowCommand()
	cmd.SetArgs([]string{"create", "--help"})
	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("xflow create help failed: %v", err)
	}

	help := out.String()
	required := []string{
		"create-xflow.json",
		"--action-user",
		"publishOpts",
		"registerOpts",
		"syncOpts",
		"reviewDoc",
		"versionID",
	}
	for _, want := range required {
		if !strings.Contains(help, want) {
			t.Fatalf("xflow create help missing %q\n%s", want, help)
		}
	}
}

func TestCreateHelpTextComesFromFiles(t *testing.T) {
	longText, err := os.ReadFile("xflow_create_long.txt")
	if err != nil {
		t.Fatalf("read xflow create long help: %v", err)
	}
	exampleText, err := os.ReadFile("xflow_create_example.txt")
	if err != nil {
		t.Fatalf("read xflow create example help: %v", err)
	}

	createCmd := findXflowSubcommand(t, "create")

	if createCmd.Long != string(longText) {
		t.Fatalf("xflow create long help should match text file\nwant:\n%s\ngot:\n%s", string(longText), createCmd.Long)
	}
	if createCmd.Example != string(exampleText) {
		t.Fatalf("xflow create example help should match text file\nwant:\n%s\ngot:\n%s", string(exampleText), createCmd.Example)
	}
}

func findXflowSubcommand(t *testing.T, name string) *cobra.Command {
	t.Helper()

	for _, cmd := range NewXflowCommand().Commands() {
		if strings.Fields(cmd.Use)[0] == name {
			return cmd
		}
	}
	t.Fatalf("xflow subcommand %q not found", name)
	return nil
}
