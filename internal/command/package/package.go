package packageCmd

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

//go:embed package_update_long.txt
var packageUpdateLong string

//go:embed package_update_example.txt
var packageUpdateExample string

func NewPackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "软件包相关命令",
	}

	var pageIndex int64
	var pageSize int64
	var onlyMine bool

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出软件包",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, args, pageIndex, pageSize, onlyMine)
		},
	}
	listCmd.Flags().Int64VarP(&pageIndex, "index", "i", 1, "页码")
	listCmd.Flags().Int64VarP(&pageSize, "size", "s", defaultPageSize, "页面大小")
	listCmd.Flags().BoolVarP(&onlyMine, "my", "m", false, "仅查看当前用户负责的记录")
	cmd.AddCommand(listCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "info",
		Short: "显示软件包信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfoCommand(cmd, args)
		},
	})

	var updateFile string
	updateCmd := &cobra.Command{
		Use:     "update <name>",
		Short:   "更新软件包信息",
		Long:    packageUpdateLong,
		Args:    cobra.ExactArgs(1),
		Example: packageUpdateExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateCommand(cmd, args[0], updateFile)
		},
	}
	updateCmd.Flags().StringVarP(&updateFile, "file", "f", "", "请求 JSON 文件路径")
	cmd.AddCommand(updateCmd)

	return cmd
}

func runInfoCommand(cmd *cobra.Command, args []string) error {
	client, err := caller.New()
	if err != nil {
		return err
	}
	if len(args) < 1 {
		return fmt.Errorf("args length must Greater than or equal to 1")
	}
	details := []dto.PackageDetail{}
	for _, name := range args {
		resp, err := client.GetPackage(context.Background(), name)
		if err != nil {
			return fmt.Errorf("failed to get package %q: %w", name, err)
		}
		if resp == nil {
			return fmt.Errorf("package %q not found", name)
		}
		details = append(details, *resp)
	}
	if err != nil {
		return err
	}
	return printJSON(cmd, details)
}

func runListCommand(cmd *cobra.Command, args []string, pageIndex, pageSize int64, onlyMine bool) error {
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

	req, err := buildMGetPackageReq(args, pageIndex, pageSize)
	if err != nil {
		return err
	}
	if onlyMine {
		req.Owner = client.Auth().User.Username
	}

	resp, err := client.MGetPackage(context.Background(), req)
	if err != nil {
		return err
	}

	return printJSON(cmd, resp)
}

func runUpdateCommand(cmd *cobra.Command, name, filePath string) error {
	client, err := caller.New()
	if err != nil {
		return err
	}

	var req dto.UpdatePackageReq
	if err := readJSONFile(filePath, &req); err != nil {
		return err
	}

	if err := client.UpdatePackage(context.Background(), name, req); err != nil {
		return err
	}

	return printJSON(cmd, map[string]string{
		"status":  "ok",
		"package": name,
	})
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

func buildMGetPackageReq(args []string, pageIndex, pageSize int64) (dto.MGetPackageReq, error) {
	if len(args)%2 != 0 {
		return dto.MGetPackageReq{}, fmt.Errorf("arguments must be provided as field/value pairs, for example: slam-cli package list name demo owner alice")
	}

	req := dto.MGetPackageReq{
		PageReq: dto.PageReq{
			Page:     pageIndex,
			PageSize: pageSize,
		},
	}

	for i := 0; i < len(args); i += 2 {
		field := strings.ToLower(strings.TrimSpace(args[i]))
		value := strings.TrimSpace(args[i+1])
		if value == "" {
			return dto.MGetPackageReq{}, fmt.Errorf("value for field %q cannot be empty", args[i])
		}

		switch field {
		case "name", "keyword", "keywords":
			req.Keywords = value
		case "owner":
			req.Owner = value
		case "dkms":
			parsed, err := parseBoolArg(value)
			if err != nil {
				return dto.MGetPackageReq{}, fmt.Errorf("invalid dkms value %q: %w", value, err)
			}
			req.Dkms = &parsed
		case "agent":
			parsed, err := parseBoolArg(value)
			if err != nil {
				return dto.MGetPackageReq{}, fmt.Errorf("invalid agent value %q: %w", value, err)
			}
			req.Agent = &parsed
		case "type":
			normalizedType := strings.ToLower(value)
			if normalizedType != "src" && normalizedType != "bin" {
				return dto.MGetPackageReq{}, fmt.Errorf("invalid type value %q: only src or bin are supported", value)
			}
			req.Type = normalizedType
		default:
			return dto.MGetPackageReq{}, fmt.Errorf("unsupported filter field %q", args[i])
		}
	}

	return req, nil
}

func parseBoolArg(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "t", "yes", "y", "1":
		return true, nil
	case "false", "f", "no", "n", "0":
		return false, nil
	default:
		return false, fmt.Errorf("supported values: true/false/yes/no/1/0")
	}
}

func printJSON(cmd *cobra.Command, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
