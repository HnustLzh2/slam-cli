package app

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	authCmd "slam-cli/internal/command/auth"
	configCmd "slam-cli/internal/command/config"
	livepatchCmd "slam-cli/internal/command/livepatch"
	packageCmd "slam-cli/internal/command/package"
	shortcutCmd "slam-cli/internal/command/shortcut"
	taskCmd "slam-cli/internal/command/task"
	thirdRepoCmd "slam-cli/internal/command/thirdrepo"
	xflowCmd "slam-cli/internal/command/xflow"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "slam-cli",
		Short:         "slam-cli 的统一命令行入口",
		Long:          "slam-cli 是一个参考 lark-cli 分层设计构建的 SLAM 命令行工具骨架。",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Version = Version
	cmd.SetVersionTemplate(fmt.Sprintf("slam-cli version %s\n", Version))
	cmd.AddCommand(NewHelpCmd())
	cmd.AddCommand(NewVersionCmd())
	cmd.AddCommand(shortcutCmd.Commands()...)
	cmd.AddCommand(authCmd.NewAuthCommand())
	cmd.AddCommand(configCmd.NewConfigCommand())
	cmd.AddCommand(packageCmd.NewPackageCommand())
	cmd.AddCommand(livepatchCmd.NewLivePatchCommand())
	cmd.AddCommand(thirdRepoCmd.NewThirdRepoCommand())
	cmd.AddCommand(taskCmd.NewTaskCommand())
	cmd.AddCommand(xflowCmd.NewXflowCommand())
	cmd.AddCommand(NewCompletionCmd(cmd))

	return cmd
}

func Execute() error {
	return NewRootCmd().Execute()
}

//go:embed help.txt
var helpStr string

func NewHelpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "显示帮助信息",
		Long:  "显示 slam-cli 的帮助信息，以及各子命令的使用方式。",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(cmd.OutOrStdout(), helpStr)
			return nil
		},
	}
}

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示当前版本",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "slam-cli version %s\n", Version)
		},
	}
}

func NewCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "生成 shell 自动补全脚本",
		Long:  "为 bash、zsh、fish 或 powershell 生成自动补全脚本。\n\n示例：\n  source <(slam-cli completion zsh)\n  source <(slam-cli completion bash)",
	}

	cmd.AddCommand(&cobra.Command{
		Use:                   "bash",
		Short:                 "生成 bash 自动补全脚本",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenBashCompletionV2(os.Stdout, true)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:                   "zsh",
		Short:                 "生成 zsh 自动补全脚本",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenZshCompletion(os.Stdout)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:                   "fish",
		Short:                 "生成 fish 自动补全脚本",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenFishCompletion(os.Stdout, true)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:                   "powershell",
		Aliases:               []string{"ps"},
		Short:                 "生成 powershell 自动补全脚本",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		},
	})

	return cmd
}
