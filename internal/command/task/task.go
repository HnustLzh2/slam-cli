package taskCmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"slam-cli/internal/pkg/caller"
	dto "slam-cli/internal/pkg/caller/DTO"
)

const defaultPageSize int64 = 20

func NewTaskCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "task 相关命令",
	}

	var pageIndex int64
	var pageSize int64
	var onlyMine bool
	var taskID int64
	var packageName string
	var taskType string
	var taskState string
	var source string
	var creator string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出 task",
		Example: "  slam-cli task list --type build --state running\n" +
			"  slam-cli task list --package kernel-agent --source xflow\n" +
			"  slam-cli task list --creator alice --my",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, dto.MGetTaskReq{
				PageReq: dto.PageReq{
					Page:     pageIndex,
					PageSize: pageSize,
				},
				ID:      taskID,
				Package: packageName,
				Type:    taskType,
				State:   taskState,
				Source:  source,
				Creator: creator,
				My:      onlyMine,
			})
		},
	}
	listCmd.Flags().Int64VarP(&pageIndex, "index", "i", 1, "页码")
	listCmd.Flags().Int64VarP(&pageSize, "size", "s", defaultPageSize, "页面大小")
	listCmd.Flags().Int64Var(&taskID, "id", 0, "按 task ID 筛选")
	listCmd.Flags().StringVar(&packageName, "package", "", "按软件包名称筛选")
	listCmd.Flags().StringVar(&taskType, "type", "", "按任务类型筛选")
	listCmd.Flags().StringVar(&taskState, "state", "", "按任务状态筛选")
	listCmd.Flags().StringVar(&source, "source", "", "按任务来源筛选")
	listCmd.Flags().StringVar(&creator, "creator", "", "按创建人筛选")
	listCmd.Flags().BoolVarP(&onlyMine, "my", "m", false, "仅查看当前用户负责的记录")
	cmd.AddCommand(listCmd)

	return cmd
}

func runListCommand(cmd *cobra.Command, req dto.MGetTaskReq) error {
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

	resp, err := client.MGetTask(context.Background(), req)
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
