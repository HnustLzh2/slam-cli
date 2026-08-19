package caller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	dto "slam-cli/internal/pkg/caller/DTO"

	"github.com/go-resty/resty/v2"

	configpkg "slam-cli/internal/pkg/config"
)

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

func (c *Client) GetPackage(ctx context.Context, name string) (*dto.PackageDetail, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("package name is required")
	}

	resp, err := c.http.R().
		SetContext(ctx).
		SetQueryParam("name", name).
		Get(c.urlFor("package", "detail"))
	if err != nil {
		return nil, fmt.Errorf("failed to call get package API: %w", err)
	}
	result, err := decodeAPIData[dto.PackageDetail](resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdatePackage(ctx context.Context, name string, req dto.UpdatePackageReq) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("package name is required")
	}
	req.Name = name

	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(req).
		Put(c.urlFor("package", url.PathEscape(name)))
	if err != nil {
		return fmt.Errorf("failed to call update package API: %w", err)
	}
	_, err = decodeAPIData[map[string]any](resp)
	return err
}

func (c *Client) MGetPackage(ctx context.Context, req dto.MGetPackageReq) (*dto.MGetPackageResp, error) {
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
	result, err := decodeAPIData[dto.MGetPackageResp](resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetAllPackages(ctx context.Context) (*dto.GetAllPackagesResp, error) {
	resp, err := c.http.R().SetContext(ctx).Get(c.urlFor("package", "all"))
	if err != nil {
		return nil, fmt.Errorf("failed to call get all packages API: %w", err)
	}
	result, err := decodeAPIData[dto.GetAllPackagesResp](resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) SearchLivePatch(ctx context.Context, req dto.SearchLivePatchReq) (*dto.SearchLivePatchResp, error) {
	r := c.http.R().SetContext(ctx)
	setOptionalIntQuery(r, "page", req.Page)
	setOptionalIntQuery(r, "pageSize", req.PageSize)
	setOptionalIntQuery(r, "id", req.ID)
	setOptionalQuery(r, "name", req.Name)
	setOptionalQuery(r, "status", req.Status)
	setOptionalQuery(r, "version", req.Version)
	setOptionalQuery(r, "distribution", req.Distribution)
	setOptionalQuery(r, "os_version", req.OsVersion)
	setOptionalQuery(r, "kernel_version", req.KernelVersion)
	setOptionalQuery(r, "arch", req.Arch)
	setOptionalQuery(r, "creator", req.Creator)
	setOptionalQuery(r, "startTime", req.StartTime)
	setOptionalQuery(r, "endTime", req.EndTime)

	resp, err := r.Get(c.urlFor("livepatch", "product"))
	if err != nil {
		return nil, fmt.Errorf("failed to call search livepatch API: %w", err)
	}
	result, err := decodeAPIData[dto.SearchLivePatchResp](resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetLivePatch(ctx context.Context, productID int64) (*dto.LivePatchProduct, error) {
	if productID <= 0 {
		return nil, fmt.Errorf("product id must be greater than 0")
	}

	resp, err := c.SearchLivePatch(ctx, dto.SearchLivePatchReq{
		PageReq: dto.PageReq{
			Page:     1,
			PageSize: 1,
		},
		ID: productID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Data) == 0 {
		return nil, nil
	}
	return &resp.Data[0], nil
}

func (c *Client) MGetThirdRepo(ctx context.Context, req dto.MGetThirdRepoReq) (*dto.MGetThirdRepoResp, error) {
	r := c.http.R().SetContext(ctx)
	setOptionalQuery(r, "name", req.Name)
	setOptionalQuery(r, "type", req.Type)
	setOptionalQuery(r, "owner", req.Owner)
	setOptionalIntQuery(r, "page", req.Page)
	setOptionalIntQuery(r, "pageSize", req.PageSize)

	resp, err := r.Get(c.urlFor("third", "config"))
	if err != nil {
		return nil, fmt.Errorf("failed to call mget third repo API: %w", err)
	}
	result, err := decodeAPIData[dto.MGetThirdRepoResp](resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) MGetTask(ctx context.Context, req dto.MGetTaskReq) (*dto.MGetTaskResp, error) {
	r := c.http.R().SetContext(ctx)
	setOptionalIntQuery(r, "page", req.Page)
	setOptionalIntQuery(r, "pageSize", req.PageSize)
	setOptionalIntQuery(r, "id", req.ID)
	setOptionalQuery(r, "package", req.Package)
	setOptionalQuery(r, "type", req.Type)
	setOptionalQuery(r, "state", req.State)
	setOptionalQuery(r, "source", req.Source)
	setOptionalQuery(r, "creator", req.Creator)
	if req.My {
		r.SetQueryParam("my", "true")
	}

	resp, err := r.Get(c.urlFor("task"))
	if err != nil {
		return nil, fmt.Errorf("failed to call mget task API: %w", err)
	}
	result, err := decodeAPIData[dto.MGetTaskResp](resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) MGetXflow(ctx context.Context, req dto.MGetXflowReq) (*dto.MGetXflowResp, error) {
	r := c.http.R().SetContext(ctx)
	setOptionalIntQuery(r, "page", req.Page)
	setOptionalIntQuery(r, "pageSize", req.PageSize)
	setOptionalIntQuery(r, "xflowID", req.XflowID)
	setOptionalQuery(r, "package", req.Package)
	setOptionalQuery(r, "type", req.Type)
	setOptionalQuery(r, "state", req.State)
	setOptionalQuery(r, "creator", req.Creator)
	if req.My {
		r.SetQueryParam("my", "true")
	}

	resp, err := r.Get(c.urlFor("xflow"))
	if err != nil {
		return nil, fmt.Errorf("failed to call mget xflow API: %w", err)
	}
	result, err := decodeAPIData[dto.MGetXflowResp](resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateXflow(ctx context.Context, req dto.CreateXflowReq) (*dto.CreateXflowResp, error) {
	r := c.http.R().SetContext(ctx)
	if strings.TrimSpace(req.ActionUser) != "" {
		r.SetHeader("X-Action-User", req.ActionUser)
	}

	resp, err := r.SetBody(req).Post(c.urlFor("xflow"))
	if err != nil {
		return nil, fmt.Errorf("failed to call create xflow API: %w", err)
	}

	xflowID, err := decodeAPIData[int64](resp)
	if err != nil {
		return nil, err
	}
	return &dto.CreateXflowResp{XflowID: *xflowID}, nil
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

	var envelope dto.ApiResponse[T]
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
