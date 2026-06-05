package authflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	persistent "github.com/juju/persistent-cookiejar"
)

const (
	hostBytedanceSSO                  = "https://sso.bytedance.com"
	pathBytedanceSSOLarkToken         = "/dingtalk/token"
	pathBytedanceSSOLarkTokenCheck    = "/dingtalk/check"
	pathBytedanceSSOLarkTokenValidate = "/dingtalk/validate"
	queryToken                        = "token"
	queryBusiness                     = "business"
	bytedanceSSOLoginPageKeyword      = "ByteDance SSO"

	pathByteCloudAuthJWT  = "/auth/api/v1/jwt"
	pathByteCloudUserInfo = "/auth/api/v1/userinfo"
	pathByteCloudLogin    = "/auth/api/v1/login"
	headerXJWTToken       = "x-jwt-token"
)

type Manager struct {
	cookieJarPath string
	tokenDir      string
	httpClient    *resty.Client
}

type CreateLarkTokenResponse struct {
	ReturnCode int64  `json:"ret"`
	Token      string `json:"token"`
	Business   string `json:"business"`
}

type CheckLarkTokenResponse struct {
	ReturnCode int64 `json:"ret"`
}

type BytedanceResponseDataBase struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type TokenInfo struct {
	Region  string    `json:"region"`
	Token   string    `json:"token"`
	SavedAt time.Time `json:"saved_at"`
}

type UserInfo struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Department  string `json:"department"`
	OpenID      string `json:"open_id"`
}

type bytedanceResponseEnvelope struct {
	Code  int             `json:"code"`
	Error string          `json:"error"`
	Data  json.RawMessage `json:"data"`
}

func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user home directory: %w", err)
	}

	baseDir := filepath.Join(homeDir, ".slam-cli")
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create auth directory: %w", err)
	}

	tokenDir := filepath.Join(baseDir, "auth")
	if err := os.MkdirAll(tokenDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create token directory: %w", err)
	}

	cookieJarPath := filepath.Join(baseDir, "cookies")
	jar, err := persistent.New(&persistent.Options{Filename: cookieJarPath})
	if err != nil {
		return nil, fmt.Errorf("failed to create persistent cookie jar: %w", err)
	}

	client := resty.New()
	client.SetCookieJar(jar)
	client.SetTimeout(30 * time.Second)

	return &Manager{
		cookieJarPath: cookieJarPath,
		tokenDir:      tokenDir,
		httpClient:    client,
	}, nil

}

func (m *Manager) saveCookies() error {
	if jar, ok := m.httpClient.GetClient().Jar.(*persistent.Jar); ok {
		return jar.Save()
	}
	return nil
}

func (m *Manager) tokenFilePath(region string) string {
	return filepath.Join(m.tokenDir, fmt.Sprintf("%s.json", normalizeRegion(region)))
}

func (m *Manager) SaveToken(region, token string) error {
	info := TokenInfo{
		Region:  normalizeRegion(region),
		Token:   token,
		SavedAt: time.Now(),
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode token info: %w", err)
	}

	if err := os.WriteFile(m.tokenFilePath(region), data, 0o600); err != nil {
		return fmt.Errorf("failed to save token file: %w", err)
	}

	return nil
}

func (m *Manager) LoadToken(region string) (*TokenInfo, error) {
	data, err := os.ReadFile(m.tokenFilePath(region))
	if err != nil {
		return nil, err
	}

	var info TokenInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return &info, nil
}

func (m *Manager) TokenFilePath(region string) string {
	return m.tokenFilePath(region)
}

func (m *Manager) CookieJarPath() string {
	return m.cookieJarPath
}

func (m *Manager) GetUserInfo(ctx context.Context, region, token string) (*UserInfo, error) {
	host := getByteCloudHostByRegion(region)

	resp, err := m.httpClient.R().
		SetContext(ctx).
		SetHeader(headerXJWTToken, token).
		Get(host + pathByteCloudUserInfo)
	if err != nil {
		return nil, fmt.Errorf("request userinfo failed: %w", err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, fmt.Errorf("userinfo unauthorized")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("userinfo HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	var envelope bytedanceResponseEnvelope
	if err := json.Unmarshal(resp.Body(), &envelope); err == nil && len(envelope.Data) > 0 {
		if envelope.Code != 0 {
			return nil, fmt.Errorf("userinfo API error (code %d): %s", envelope.Code, envelope.Error)
		}
		return parseUserInfo(envelope.Data)
	}

	return parseUserInfo(resp.Body())
}

func (m *Manager) Login(ctx context.Context, region string) (string, error) {
	region = normalizeRegion(region)

	token, err := m.GenerateJWT(ctx, region)
	if err == nil {
		if saveErr := m.SaveToken(region, token); saveErr != nil {
			return "", saveErr
		}
		if saveErr := m.saveCookies(); saveErr != nil {
			return "", saveErr
		}
		return token, nil
	}

	startURL := getByteCloudHostByRegion(region) + pathByteCloudLogin
	isRedirectToSSO, err := m.tryRequestStartURL(ctx, startURL)
	if err != nil {
		return "", fmt.Errorf("failed to check login status: %w", err)
	}

	if isRedirectToSSO {
		if err := m.bytedanceSSOLogin(ctx, startURL); err != nil {
			return "", fmt.Errorf("SSO login failed: %w", err)
		}
	}

	token, err = m.GenerateJWT(ctx, region)
	if err != nil {
		return "", fmt.Errorf("login verification failed: %w", err)
	}

	if err := m.SaveToken(region, token); err != nil {
		return "", err
	}
	if err := m.saveCookies(); err != nil {
		return "", err
	}

	return token, nil

}

func (m *Manager) tryRequestStartURL(ctx context.Context, startURL string) (bool, error) {
	resp, err := m.httpClient.R().SetContext(ctx).Get(startURL)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(resp.Body()), bytedanceSSOLoginPageKeyword), nil
}

func (m *Manager) bytedanceSSOLogin(ctx context.Context, startURL string) error {
	for {
		resp, err := m.httpClient.R().SetContext(ctx).Post(hostBytedanceSSO + pathBytedanceSSOLarkToken)
		if err != nil {
			return fmt.Errorf("failed to create lark token: %w", err)
		}

		var larkToken CreateLarkTokenResponse
		if err := json.Unmarshal(resp.Body(), &larkToken); err != nil {
			return fmt.Errorf("failed to parse lark token response: %w", err)
		}

		qrCodeURL, _ := url.Parse(hostBytedanceSSO + pathBytedanceSSOLarkTokenValidate)
		qrCodeQuery := qrCodeURL.Query()
		qrCodeQuery.Set(queryToken, larkToken.Token)
		qrCodeQuery.Set(queryBusiness, larkToken.Business)
		qrCodeURL.RawQuery = qrCodeQuery.Encode()

		fmt.Println("请使用飞书扫码完成登录：")
		fmt.Printf("URL: %s\n\n", qrCodeURL.String())
		PrintLoginQR(qrCodeURL.String())

		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return ctx.Err()
			case <-ticker.C:
				resp, err := m.httpClient.R().
					SetContext(ctx).
					SetQueryParam(queryBusiness, larkToken.Business).
					SetQueryParam(queryToken, larkToken.Token).
					Get(hostBytedanceSSO + pathBytedanceSSOLarkTokenCheck)
				if err != nil {
					ticker.Stop()
					return fmt.Errorf("failed to check login status: %w", err)
				}

				var checkResp CheckLarkTokenResponse
				if err := json.Unmarshal(resp.Body(), &checkResp); err != nil {
					ticker.Stop()
					return fmt.Errorf("failed to parse check response: %w", err)
				}

				switch checkResp.ReturnCode {
				case 0:
					continue
				case 1:
					ticker.Stop()
					fmt.Println("扫码成功，正在建立会话...")
					_, err := m.tryRequestStartURL(ctx, startURL)
					if err != nil {
						return fmt.Errorf("failed to establish session after SSO login: %w", err)
					}
					return nil
				case 2:
					ticker.Stop()
					fmt.Println("二维码已过期，正在重新生成...")
					goto NEXT_QR
				default:
					ticker.Stop()
					return fmt.Errorf("unknown return code: %d", checkResp.ReturnCode)
				}
			}
		}
	NEXT_QR:
	}
}

func (m *Manager) GenerateJWT(ctx context.Context, region string) (string, error) {
	host := getByteCloudHostByRegion(region)

	resp, err := m.httpClient.R().SetContext(ctx).Get(host + pathByteCloudAuthJWT)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return "", fmt.Errorf("unauthorized - please run 'slam-cli auth login' first")
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	var responseData BytedanceResponseDataBase
	if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if responseData.Code != 0 {
		return "", fmt.Errorf("API error (code %d): %s", responseData.Code, responseData.Error)
	}

	jwtToken := resp.Header().Get(headerXJWTToken)
	if jwtToken == "" {
		return "", fmt.Errorf("JWT token not found in response headers")
	}

	return jwtToken, nil
}

func getByteCloudHostByRegion(region string) string {
	switch normalizeRegion(region) {
	case "cn":
		return "https://cloud.bytedance.net"
	case "us":
		return "https://console.volcengine.com"
	case "sg":
		return "https://console.volces.com"
	default:
		return "https://cloud.bytedance.net"
	}
}

func normalizeRegion(region string) string {
	if region == "" {
		return "cn"
	}
	return strings.ToLower(region)
}

func parseUserInfo(data []byte) (*UserInfo, error) {
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo response: %w", err)
	}

	info := &UserInfo{
		Username:    firstNonEmptyString(findString(payload, "username"), findString(payload, "user_name"), findString(payload, "name")),
		Email:       firstNonEmptyString(findString(payload, "email"), findString(payload, "mail")),
		DisplayName: firstNonEmptyString(findString(payload, "display_name"), findString(payload, "nickname"), findString(payload, "nick_name"), findString(payload, "cn_name")),
		AvatarURL:   firstNonEmptyString(findString(payload, "avatar_url"), findString(payload, "avatar"), findString(payload, "avatarUrl")),
		Department:  firstNonEmptyString(findString(payload, "department_name"), findString(payload, "department"), findString(payload, "dept_name")),
		OpenID:      firstNonEmptyString(findString(payload, "open_id"), findString(payload, "openid"), findString(payload, "openId")),
	}

	if info.Username == "" && info.Email == "" && info.DisplayName == "" {
		return nil, fmt.Errorf("userinfo response does not contain recognizable user fields")
	}

	return info, nil
}

func findString(v any, key string) string {
	switch value := v.(type) {
	case map[string]any:
		for k, nested := range value {
			if strings.EqualFold(k, key) {
				if s, ok := nested.(string); ok {
					return s
				}
			}
			if s := findString(nested, key); s != "" {
				return s
			}
		}
	case []any:
		for _, item := range value {
			if s := findString(item, key); s != "" {
				return s
			}
		}
	}
	return ""
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
