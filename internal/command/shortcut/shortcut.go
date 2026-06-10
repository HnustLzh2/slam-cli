package shortcut

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	zshCompletionLine  = "source <(slam-cli completion zsh)"
	bashCompletionLine = "source <(slam-cli completion bash)"
)

func Commands() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   "+doctor",
			Short: "检查本地 slam-cli 运行环境",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Fprintln(cmd.OutOrStdout(), "slam-cli doctor: ok (skeleton mode)")
			},
		},
		{
			Use:   "+init",
			Short: "初始化 slam-cli 本地工作环境",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runInit(cmd)
			},
		},
	}
}

func runInit(cmd *cobra.Command) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	zshrcPath := filepath.Join(homeDir, ".zshrc")
	bashrcPath := filepath.Join(homeDir, ".bashrc")

	zshAdded, err := ensureShellCompletion(zshrcPath, zshCompletionLine)
	if err != nil {
		return err
	}
	bashAdded, err := ensureShellCompletion(bashrcPath, bashCompletionLine)
	if err != nil {
		return err
	}

	if err := sourceShellRC("zsh", zshrcPath); err != nil {
		return err
	}
	if err := sourceShellRC("bash", bashrcPath); err != nil {
		return err
	}

	if zshAdded {
		fmt.Fprintf(cmd.OutOrStdout(), "已写入 %s\n", zshrcPath)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "%s 已存在 slam-cli 补全配置\n", zshrcPath)
	}
	if bashAdded {
		fmt.Fprintf(cmd.OutOrStdout(), "已写入 %s\n", bashrcPath)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "%s 已存在 slam-cli 补全配置\n", bashrcPath)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "已执行 source ~/.zshrc 和 source ~/.bashrc")

	return nil
}

func ensureShellCompletion(filePath, sourceLine string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	text := string(content)
	if strings.Contains(text, sourceLine) {
		return false, nil
	}

	block := "\n# ---------- slam-cli completion ----------\nif command -v slam-cli >/dev/null 2>&1; then\n  " + sourceLine + "\nfi\n"
	if strings.TrimSpace(text) == "" {
		block = strings.TrimLeft(block, "\n")
	}

	updated := text + block
	if err := os.WriteFile(filePath, []byte(updated), 0644); err != nil {
		return false, fmt.Errorf("failed to write %s: %w", filePath, err)
	}
	return true, nil
}

func sourceShellRC(shell, rcPath string) error {
	command := exec.Command(shell, "-lc", fmt.Sprintf("source %q", rcPath))
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to source %s with %s: %w, output: %s", rcPath, shell, err, strings.TrimSpace(string(output)))
	}
	return nil
}
