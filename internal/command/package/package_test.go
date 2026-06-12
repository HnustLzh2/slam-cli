package packageCmd

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestUpdateHelpTextComesFromFiles(t *testing.T) {
	longText, err := os.ReadFile("package_update_long.txt")
	if err != nil {
		t.Fatalf("read package update long help: %v", err)
	}
	exampleText, err := os.ReadFile("package_update_example.txt")
	if err != nil {
		t.Fatalf("read package update example help: %v", err)
	}

	updateCmd := findPackageSubcommand(t, "update")

	if updateCmd.Long != string(longText) {
		t.Fatalf("package update long help should match text file\nwant:\n%s\ngot:\n%s", string(longText), updateCmd.Long)
	}
	if updateCmd.Example != string(exampleText) {
		t.Fatalf("package update example help should match text file\nwant:\n%s\ngot:\n%s", string(exampleText), updateCmd.Example)
	}
}

func findPackageSubcommand(t *testing.T, name string) *cobra.Command {
	t.Helper()

	for _, cmd := range NewPackageCommand().Commands() {
		if strings.Fields(cmd.Use)[0] == name {
			return cmd
		}
	}
	t.Fatalf("package subcommand %q not found", name)
	return nil
}
