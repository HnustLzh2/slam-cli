package xflowCmd

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"slam-cli/internal/pkg/caller"
	dto "slam-cli/internal/pkg/caller/DTO"
)

const defaultPageSize int64 = 20

//go:embed xflow_create_long.txt
var xflowCreateLong string

//go:embed xflow_create_example.txt
var xflowCreateExample string

func NewXflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "xflow",
		Short: "xflow 相关命令",
	}

	var pageIndex int64
	var pageSize int64
	var onlyMine bool
	var xflowID int64
	var packageName string
	var xflowType string
	var state string
	var creator string
	var createFile string
	var actionUser string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出 xflow 工单",
		Example: "  slam-cli xflow list --type publish --state running\n" +
			"  slam-cli xflow list --package kernel-agent --creator alice\n" +
			"  slam-cli xflow list --id 123456 --my",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, dto.MGetXflowReq{
				PageReq: dto.PageReq{
					Page:     pageIndex,
					PageSize: pageSize,
				},
				XflowID: xflowID,
				Package: packageName,
				Type:    xflowType,
				State:   state,
				Creator: creator,
				My:      onlyMine,
			})
		},
	}
	listCmd.Flags().Int64VarP(&pageIndex, "index", "i", 1, "页码")
	listCmd.Flags().Int64VarP(&pageSize, "size", "s", defaultPageSize, "页面大小")
	listCmd.Flags().Int64Var(&xflowID, "id", 0, "按 xflow ID 筛选")
	listCmd.Flags().StringVar(&packageName, "package", "", "按软件包名称筛选")
	listCmd.Flags().StringVar(&xflowType, "type", "", "按工单类型筛选")
	listCmd.Flags().StringVar(&state, "state", "", "按工单状态筛选")
	listCmd.Flags().StringVar(&creator, "creator", "", "按创建人筛选")
	listCmd.Flags().BoolVarP(&onlyMine, "my", "m", false, "仅查看当前用户负责的记录")
	cmd.AddCommand(listCmd)

	createCmd := &cobra.Command{
		Use:     "create",
		Short:   "创建 xflow 工单",
		Long:    xflowCreateLong,
		Example: xflowCreateExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateCommand(cmd, createFile, actionUser)
		},
	}
	createCmd.Flags().StringVarP(&createFile, "file", "f", "", "请求 JSON 文件路径")
	createCmd.Flags().StringVar(&actionUser, "action-user", "", "写入 X-Action-User header，用于服务账号代用户发起场景")
	cmd.AddCommand(createCmd)

	return cmd
}

func runListCommand(cmd *cobra.Command, req dto.MGetXflowReq) error {
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

	resp, err := client.MGetXflow(context.Background(), req)
	if err != nil {
		return err
	}

	return printJSON(cmd, resp)
}

func runCreateCommand(cmd *cobra.Command, filePath, actionUser string) error {
	client, err := caller.New()
	if err != nil {
		return err
	}

	var req dto.CreateXflowReq
	if err := readJSONFile(filePath, &req); err != nil {
		return err
	}
	if strings.TrimSpace(actionUser) != "" {
		req.ActionUser = actionUser
	}

	resp, err := client.CreateXflow(context.Background(), req)
	if err != nil {
		return err
	}
	return printJSON(cmd, resp)
}

func readJSONFile(filePath string, dst any) error {
	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("--file is required")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %s: %w", filePath, err)
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("failed to parse JSON file %s: %w", filePath, err)
	}
	return nil
}

func printJSON(cmd *cobra.Command, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
