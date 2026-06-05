package shortcut

import (
	"fmt"

	"github.com/spf13/cobra"
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
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Fprintln(cmd.OutOrStdout(), "slam-cli init: not implemented yet")
			},
		},
	}
}
