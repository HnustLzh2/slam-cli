package app

import (
	"fmt"

	"github.com/spf13/cobra"

	authCmd "slam-cli/internal/command/auth"
	configCmd "slam-cli/internal/command/config"
	packageCmd "slam-cli/internal/command/package"
	shortcutCmd "slam-cli/internal/command/shortcut"
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

	return cmd
}

func Execute() error {
	return NewRootCmd().Execute()
}

var helpStr = fmt.Sprintf("1. %s\n2. %s \n3. %s", "slam-cli v0.1.0 用于获取软件包详细信息和全量获取软件包信息",
	"slam package list [-s -i -a] 用于获取软件包列表, -s 代表页面大小 -i代表页码",
	"slam package <package-name>... 用于获取某个软件包的详细信息")

func NewHelpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "显示帮助信息",
		Long:  "显示 slam-cli 的帮助信息，以及各子命令的使用方式。",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), helpStr)
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
