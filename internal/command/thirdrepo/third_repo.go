package thirdRepCmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"slam-cli/internal/pkg/caller"
	dto "slam-cli/internal/pkg/caller/DTO"
)

const defaultPageSize int64 = 20

func NewThirdRepoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "third-repo",
		Short: "third repo 相关命令",
	}

	var pageIndex int64
	var pageSize int64

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出 third repo package",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, pageIndex, pageSize)
		},
	}
	listCmd.Flags().Int64VarP(&pageIndex, "index", "i", 1, "页码")
	listCmd.Flags().Int64VarP(&pageSize, "size", "s", defaultPageSize, "页面大小")
	cmd.AddCommand(listCmd)

	return cmd
}

func runListCommand(cmd *cobra.Command, pageIndex, pageSize int64) error {
	client, err := caller.New()
	if err != nil {
		return err
	}

	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	resp, err := client.MGetThirdRepo(context.Background(), dto.MGetThirdRepoReq{
		PageReq: dto.PageReq{
			Page:     pageIndex,
			PageSize: pageSize,
		},
	})
	if err != nil {
		return err
	}

	return printJSON(cmd, resp)
}

func printJSON(cmd *cobra.Command, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
