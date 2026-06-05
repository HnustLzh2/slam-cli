package caller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	configpkg "slam-cli/internal/pkg/config"
)

const (
	defaultSLAMAPIBaseURL = "https://slam.byted.org/api/slam/v2"
	headerXJWTToken       = "X-Jwt-Token"
)

type Client struct {
	baseURL string
	token   string
	http    *resty.Client
	auth    configpkg.AuthConfig
}

type tokenFile struct {
	Region  string    `json:"region"`
	Token   string    `json:"token"`
	SavedAt time.Time `json:"saved_at"`
}

type PageReq struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

type apiResponse[T any] struct {
	Code  int    `json:"code"`
	Data  T      `json:"data"`
	Error string `json:"error"`
}

type MGetPackageReq struct {
	Keywords string `json:"keywords"`
	Owner    string `json:"owner"`
	PageReq
	Dkms  *bool  `json:"dkms"`
	Agent *bool  `json:"agent"`
	Type  string `json:"type"`
}

type PageInfo struct {
	Page      int64 `json:"page"`
	PageSize  int64 `json:"pageSize"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"totalPage"`
}

type SimplePackage struct {
	Name       string   `json:"name"`
	Owner      []string `json:"owner"`
	Type       string   `json:"type"`
	Repo       string   `json:"repo"`
	Branch     string   `json:"branch"`
	GoVersion  string   `json:"goVersion"`
	Region     []string `json:"region"`
	Dist       []string `json:"dist"`
	Dkms       bool     `json:"dkms"`
	Agent      bool     `json:"agent"`
	IsComplete bool     `json:"isComplete"`
}

type PackageVersion struct {
	ID          int64    `json:"id"`
	Version     string   `json:"version"`
	Type        string   `json:"type"`
	Repo        string   `json:"repo"`
	Branch      string   `json:"branch"`
	Commit      string   `json:"commit"`
	GoVersion   string   `json:"goVersion"`
	Region      []string `json:"region"`
	Dist        []string `json:"dist"`
	Applicant   string   `json:"applicant"`
	VersionNote string   `json:"versionNote"`
	Docs        []string `json:"docs"`
	URL         []string `json:"url"`
	PublishTime string   `json:"publishTime"`
}

type PackageDetail struct {
	Name                string           `json:"name"`
	Owner               []string         `json:"owner"`
	Type                string           `json:"type"`
	IsComplete          bool             `json:"isComplete"`
	Repo                string           `json:"repo"`
	Branch              string           `json:"branch"`
	GoVersion           string           `json:"goVersion"`
	Region              []string         `json:"region"`
	Dist                []string         `json:"dist"`
	Dkms                bool             `json:"dkms"`
	Agent               bool             `json:"agent"`
	PublishNotifyChatID string           `json:"publishNotifyChatID"`
	Desc                string           `json:"desc"`
	AgentType           string           `json:"agentType"`
	ResourceLimit       string           `json:"resourceLimit"`
	OpenSource          bool             `json:"openSource"`
	CodeLanguage        []string         `json:"codeLanguage"`
	InvolvedHardtype    string           `json:"involvedHardtype"`
	HasNDA              bool             `json:"hasNDA"`
	Versions            []PackageVersion `json:"versions"`
	Raw                 map[string]any   `json:"-"`
}

type MGetPackageResp struct {
	Packages []SimplePackage `json:"packages"`
	PageInfo *PageInfo       `json:"pageInfo"`
	Raw      map[string]any  `json:"-"`
}

type ThirdRepoPackage struct {
	RepoName           string   `json:"repoName"`
	Owner              []string `json:"owner"`
	Level              int64    `json:"level"`
	Type               string   `json:"type"`
	OperatorWhitelists []string `json:"operatorWhitelists"`
	URL                string   `json:"url"`
	Distribution       []string `json:"distribution"`
	Component          []string `json:"component"`
	Architecture       []string `json:"architecture"`
	XflowProcess       string   `json:"xflowProcess"`
}

type GetAllPackagesResp struct {
	PrivatePackages   []SimplePackage    `json:"privatePackages"`
	ThirdRepoPackages []ThirdRepoPackage `json:"thirdRepoPackages"`
	Raw               map[string]any     `json:"-"`
}

func New() (*Client, error) {
	authCfg, err := configpkg.RequireLogin()
	if err != nil {
		return nil, err
	}

	token, err := loadTokenFromFile(authCfg.TokenFile)
	if err != nil {
		return nil, err
	}

	baseURL := resolveBaseURL()
	httpClient := resty.New()
	httpClient.SetTimeout(30 * time.Second)
	httpClient.SetHeader(headerXJWTToken, token)
	httpClient.SetHeader("Accept", "application/json")

	return &Client{
		baseURL: baseURL,
		token:   token,
		http:    httpClient,
		auth:    authCfg,
	}, nil
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) Auth() configpkg.AuthConfig {
	return c.auth
}

func (c *Client) GetPackage(ctx context.Context, name string) (*PackageDetail, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("package name is required")
	}

	resp, err := c.http.R().SetContext(ctx).Get(c.urlFor("package", url.PathEscape(name)))
	if err != nil {
		return nil, fmt.Errorf("failed to call get package API: %w", err)
	}
	result, err := decodeAPIData[PackageDetail](resp)
	if err != nil {
		return nil, err
	}
	result.Raw = rawMap(result)
	return result, nil
}

func (c *Client) MGetPackage(ctx context.Context, req MGetPackageReq) (*MGetPackageResp, error) {
	r := c.http.R().SetContext(ctx)
	if req.Keywords != "" {
		setOptionalQuery(r, "keywords", req.Keywords)
	}
	if req.Agent != nil {
		setOptionalBoolQuery(r, "agent", req.Agent)
	}
	if req.Dkms != nil {
		setOptionalBoolQuery(r, "dkms", req.Dkms)
	}
	if req.Type != "" {
		setOptionalQuery(r, "type", req.Type)

	}
	if req.Owner != "" {
		setOptionalQuery(r, "owner", req.Owner)

	}
	setOptionalIntQuery(r, "page", req.Page)
	setOptionalIntQuery(r, "pageSize", req.PageSize)

	resp, err := r.Get(c.urlFor("package"))
	if err != nil {
		return nil, fmt.Errorf("failed to call mget package API: %w", err)
	}
	result, err := decodeAPIData[MGetPackageResp](resp)
	if err != nil {
		return nil, err
	}
	result.Raw = rawMap(result)
	return result, nil
}

func (c *Client) GetAllPackages(ctx context.Context) (*GetAllPackagesResp, error) {
	resp, err := c.http.R().SetContext(ctx).Get(c.urlFor("package", "all"))
	if err != nil {
		return nil, fmt.Errorf("failed to call get all packages API: %w", err)
	}
	result, err := decodeAPIData[GetAllPackagesResp](resp)
	if err != nil {
		return nil, err
	}
	result.Raw = rawMap(result)
	return result, nil
}

func (c *Client) urlFor(segments ...string) string {
	parts := []string{strings.TrimRight(c.baseURL, "/")}
	for _, segment := range segments {
		if strings.TrimSpace(segment) == "" {
			continue
		}
		parts = append(parts, strings.Trim(segment, "/"))
	}
	return strings.Join(parts, "/")
}

func resolveBaseURL() string {
	if value := strings.TrimSpace(os.Getenv("SLAM_API_BASE_URL")); value != "" {
		return strings.TrimRight(value, "/")
	}
	if host := strings.TrimSpace(os.Getenv("SLAM_API_HOST")); host != "" {
		return strings.TrimRight(host, "/") + "/api/slam/v2"
	}
	return defaultSLAMAPIBaseURL
}

func loadTokenFromFile(filePath string) (string, error) {
	if strings.TrimSpace(filePath) == "" {
		return "", fmt.Errorf("token file path is empty")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read token file %s: %w", filePath, err)
	}

	var tokenInfo tokenFile
	if err := json.Unmarshal(data, &tokenInfo); err != nil {
		return "", fmt.Errorf("failed to parse token file %s: %w", filePath, err)
	}
	if strings.TrimSpace(tokenInfo.Token) == "" {
		return "", fmt.Errorf("token file %s does not contain token", filePath)
	}
	return tokenInfo.Token, nil
}

func checkResponse(resp *resty.Response) error {
	if resp == nil {
		return fmt.Errorf("empty HTTP response")
	}
	if resp.IsSuccess() {
		return nil
	}
	return fmt.Errorf("slam backend returned HTTP %d: %s", resp.StatusCode(), strings.TrimSpace(resp.String()))
}

func decodeAPIData[T any](resp *resty.Response) (*T, error) {
	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var envelope apiResponse[T]
	if err := json.Unmarshal(resp.Body(), &envelope); err != nil {
		return nil, fmt.Errorf("failed to decode slam response: %w", err)
	}
	if envelope.Code != 0 {
		message := strings.TrimSpace(envelope.Error)
		if message == "" {
			message = fmt.Sprintf("slam backend returned business code %d", envelope.Code)
		}
		return nil, fmt.Errorf("%s", message)
	}

	result := envelope.Data
	return &result, nil
}

func rawMap(value any) map[string]any {
	data, err := json.Marshal(value)
	if err != nil {
		return map[string]any{}
	}

	result := map[string]any{}
	if err := json.Unmarshal(data, &result); err != nil {
		return map[string]any{}
	}
	return result
}

func setOptionalQuery(r *resty.Request, key, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	r.SetQueryParam(key, value)
}

func setOptionalIntQuery(r *resty.Request, key string, value int64) {
	if value <= 0 {
		return
	}
	r.SetQueryParam(key, fmt.Sprintf("%d", value))
}

func setOptionalBoolQuery(r *resty.Request, key string, value *bool) {
	if value == nil {
		return
	}
	if *value {
		r.SetQueryParam(key, "true")
	} else {
		r.SetQueryParam(key, "false")
	}
}
