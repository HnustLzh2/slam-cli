package packagecmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"slam-cli/internal/pkg/caller"
)

const defaultPageSize int64 = 20

func NewPackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "软件包相关命令",
	}

	var pageIndex int64
	var pageSize int64

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出软件包",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, args, pageIndex, pageSize)
		},
	}
	listCmd.Flags().Int64VarP(&pageIndex, "index", "i", 1, "页码")
	listCmd.Flags().Int64VarP(&pageSize, "size", "s", defaultPageSize, "页面大小")
	cmd.AddCommand(listCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "info",
		Short: "显示软件包信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfoCommand(cmd, args)
		},
	})

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
	details := []caller.PackageDetail{}
	for _, name := range args {
		resp, err := client.GetPackage(context.Background(), name)
		if err != nil{
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

func runListCommand(cmd *cobra.Command, args []string, pageIndex, pageSize int64) error {
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

	resp, err := client.MGetPackage(context.Background(), req)
	if err != nil {
		return err
	}

	return printJSON(cmd, resp)
}

func buildMGetPackageReq(args []string, pageIndex, pageSize int64) (caller.MGetPackageReq, error) {
	if len(args)%2 != 0 {
		return caller.MGetPackageReq{}, fmt.Errorf("arguments must be provided as field/value pairs, for example: slam-cli package list name demo owner alice")
	}

	req := caller.MGetPackageReq{
		PageReq: caller.PageReq{
			Page:     pageIndex,
			PageSize: pageSize,
		},
	}

	for i := 0; i < len(args); i += 2 {
		field := strings.ToLower(strings.TrimSpace(args[i]))
		value := strings.TrimSpace(args[i+1])
		if value == "" {
			return caller.MGetPackageReq{}, fmt.Errorf("value for field %q cannot be empty", args[i])
		}

		switch field {
		case "name", "keyword", "keywords":
			req.Keywords = value
		case "owner":
			req.Owner = value
		case "dkms":
			parsed, err := parseBoolArg(value)
			if err != nil {
				return caller.MGetPackageReq{}, fmt.Errorf("invalid dkms value %q: %w", value, err)
			}
			req.Dkms = &parsed
		case "agent":
			parsed, err := parseBoolArg(value)
			if err != nil {
				return caller.MGetPackageReq{}, fmt.Errorf("invalid agent value %q: %w", value, err)
			}
			req.Agent = &parsed
		case "type":
			normalizedType := strings.ToLower(value)
			if normalizedType != "src" && normalizedType != "bin" {
				return caller.MGetPackageReq{}, fmt.Errorf("invalid type value %q: only src or bin are supported", value)
			}
			req.Type = normalizedType
		default:
			return caller.MGetPackageReq{}, fmt.Errorf("unsupported filter field %q", args[i])
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
