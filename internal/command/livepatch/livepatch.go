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
	var onlyMine bool
	var productID int64
	var name string
	var status string
	var version string
	var distribution string
	var osVersion string
	var kernelVersion string
	var arch string
	var creator string
	var startTime string
	var endTime string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出 livepatch product",
		Example: "  slam-cli livepatch list --status testing --creator alice\n" +
			"  slam-cli livepatch list --distribution debian --os-version 12 --arch amd64\n" +
			"  slam-cli livepatch list --id 1001 --my",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, dto.SearchLivePatchReq{
				PageReq: dto.PageReq{
					Page:     pageIndex,
					PageSize: pageSize,
				},
				ID:            productID,
				Name:          name,
				Status:        status,
				Version:       version,
				Distribution:  distribution,
				OsVersion:     osVersion,
				KernelVersion: kernelVersion,
				Arch:          arch,
				Creator:       creator,
				StartTime:     startTime,
				EndTime:       endTime,
			})
		},
	}
	listCmd.Flags().Int64VarP(&pageIndex, "index", "i", 1, "页码")
	listCmd.Flags().Int64VarP(&pageSize, "size", "s", defaultPageSize, "页面大小")
	listCmd.Flags().Int64Var(&productID, "id", 0, "按 livepatch product ID 筛选")
	listCmd.Flags().StringVar(&name, "name", "", "按产品名称筛选")
	listCmd.Flags().StringVar(&status, "status", "", "按状态筛选")
	listCmd.Flags().StringVar(&version, "version", "", "按版本筛选")
	listCmd.Flags().StringVar(&distribution, "distribution", "", "按发行版筛选")
	listCmd.Flags().StringVar(&osVersion, "os-version", "", "按 OS 版本筛选")
	listCmd.Flags().StringVar(&kernelVersion, "kernel-version", "", "按内核版本筛选")
	listCmd.Flags().StringVar(&arch, "arch", "", "按架构筛选")
	listCmd.Flags().StringVar(&creator, "creator", "", "按创建人筛选")
	listCmd.Flags().StringVar(&startTime, "start-time", "", "按开始时间筛选，对应原始 API startTime")
	listCmd.Flags().StringVar(&endTime, "end-time", "", "按结束时间筛选，对应原始 API endTime")
	listCmd.Flags().BoolVarP(&onlyMine, "my", "m", false, "仅查看当前用户负责的记录")
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

func runListCommand(cmd *cobra.Command, req dto.SearchLivePatchReq) error {
	client, err := caller.New()
	if err != nil {
		return err
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = defaultPageSize
	}
	if cmd.Flags().Changed("my") && cmd.Flag("my").Value.String() == "true" {
		req.Creator = client.Auth().User.Username
	}

	resp, err := client.SearchLivePatch(context.Background(), req)
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
