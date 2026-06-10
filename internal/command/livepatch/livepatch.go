package livepatchCmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"slam-cli/internal/pkg/caller"
	dto "slam-cli/internal/pkg/caller/DTO"
)

const defaultPageSize int64 = 20

func NewLivePatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "livepatch <product-id>",
		Short: "livepatch 相关命令",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return runInfoCommand(cmd, args[0])
		},
	}

	var pageIndex int64
	var pageSize int64

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出 livepatch product",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, pageIndex, pageSize)
		},
	}
	listCmd.Flags().Int64VarP(&pageIndex, "index", "i", 1, "页码")
	listCmd.Flags().Int64VarP(&pageSize, "size", "s", defaultPageSize, "页面大小")
	cmd.AddCommand(listCmd)

	return cmd
}

func runInfoCommand(cmd *cobra.Command, productIDText string) error {
	productID, err := strconv.ParseInt(productIDText, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid livepatch product id %q: %w", productIDText, err)
	}

	client, err := caller.New()
	if err != nil {
		return err
	}

	resp, err := client.GetLivePatch(context.Background(), productID)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("livepatch product %d not found", productID)
	}

	return printJSON(cmd, resp)
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

	resp, err := client.SearchLivePatch(context.Background(), dto.SearchLivePatchReq{
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
