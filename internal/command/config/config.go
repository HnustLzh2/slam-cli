package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	configpkg "slam-cli/internal/pkg/config"
)

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "配置相关命令",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printCurrentConfig(cmd)
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "列出当前配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printCurrentConfig(cmd)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "获取当前配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printCurrentConfig(cmd)
		},
	})

	return cmd
}

func printCurrentConfig(cmd *cobra.Command) error {
	cfg, err := configpkg.Load()
	if err != nil {
		return err
	}

	configPath, err := configpkg.FilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "config file: %s\n", configPath)
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
