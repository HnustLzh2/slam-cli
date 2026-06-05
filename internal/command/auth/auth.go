package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	authflow "slam-cli/internal/pkg/authflow"
	configpkg "slam-cli/internal/pkg/config"
)

func NewAuthCommand() *cobra.Command {
	var region string

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "认证相关命令",
	}
	cmd.PersistentFlags().StringVarP(&region, "region", "r", "cn", "认证区域 (cn, us, sg)")

	cmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "登录",
		Long:  "通过飞书扫码方式完成 ByteDance SSO 认证，并获取/保存 JWT token。",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := authflow.NewManager()
			if err != nil {
				return err
			}

			ctx := context.Background()
			token, err := manager.Login(ctx, region)
			if err != nil {
				return err
			}

			userInfo, err := manager.GetUserInfo(ctx, region, token)
			if err != nil {
				return err
			}

			if err := configpkg.SaveAuth(configpkg.AuthConfig{
				LoggedIn:     true,
				Region:       region,
				TokenFile:    manager.TokenFilePath(region),
				CookieJar:    manager.CookieJarPath(),
				TokenPreview: configpkg.PreviewToken(token),
				User: configpkg.UserInfo{
					Username:    userInfo.Username,
					Email:       userInfo.Email,
					DisplayName: userInfo.DisplayName,
					AvatarURL:   userInfo.AvatarURL,
					Department:  userInfo.Department,
					OpenID:      userInfo.OpenID,
				},
				UpdatedAt: time.Now(),
			}); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "登录成功，已保存认证信息。")
			configPath, err := configpkg.FilePath()
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "region: %s\n", region)
			fmt.Fprintf(cmd.OutOrStdout(), "config file: %s\n", configPath)
			fmt.Fprintf(cmd.OutOrStdout(), "cookie jar: %s\n", manager.CookieJarPath())
			fmt.Fprintf(cmd.OutOrStdout(), "token file: %s\n", manager.TokenFilePath(region))
			fmt.Fprintf(cmd.OutOrStdout(), "username: %s\n", userInfo.Username)
			fmt.Fprintf(cmd.OutOrStdout(), "email: %s\n", userInfo.Email)
			fmt.Fprintln(cmd.OutOrStdout(), "token:")
			fmt.Fprintln(cmd.OutOrStdout(), token)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "查看当前认证状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			authCfg, err := configpkg.CurrentAuth()
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "region: %s\n", authCfg.Region)
			if !authCfg.IsLoggedIn() {
				fmt.Fprintln(cmd.OutOrStdout(), "status: 未登录")
				fmt.Fprintln(cmd.OutOrStdout(), "请先执行 'slam-cli auth login' 完成认证。")
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), "status: 已登录")
			fmt.Fprintf(cmd.OutOrStdout(), "updated at: %s\n", authCfg.UpdatedAt.Format("2006-01-02 15:04:05"))
			fmt.Fprintf(cmd.OutOrStdout(), "token file: %s\n", authCfg.TokenFile)
			fmt.Fprintf(cmd.OutOrStdout(), "token preview: %s\n", authCfg.TokenPreview)
			fmt.Fprintf(cmd.OutOrStdout(), "username: %s\n", authCfg.User.Username)
			fmt.Fprintf(cmd.OutOrStdout(), "email: %s\n", authCfg.User.Email)
			return nil
		},
	})

	return cmd
}
