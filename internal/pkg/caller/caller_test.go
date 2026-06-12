package caller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"

	dto "slam-cli/internal/pkg/caller/DTO"
)

func TestUpdatePackageUsesPutPackageName(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody dto.UpdatePackageReq
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(dto.ApiResponse[map[string]any]{Code: 0, Data: map[string]any{}})
	}))
	defer server.Close()

	client := &Client{baseURL: server.URL, http: newHTTPClientForTest()}
	err := client.UpdatePackage(context.Background(), "demo/pkg", dto.UpdatePackageReq{
		Owner:     []string{"alice"},
		Type:      "src",
		Repo:      "git@example.com/repo.git",
		Branch:    "main",
		GoVersion: "1.22",
		Region:    []string{"cn"},
		Dist:      []string{"Debian 12"},
		Desc:      "demo package",
	})
	if err != nil {
		t.Fatalf("UpdatePackage returned error: %v", err)
	}

	if gotMethod != http.MethodPut {
		t.Fatalf("method = %q, want %q", gotMethod, http.MethodPut)
	}
	if gotPath != "/package/demo%2Fpkg" {
		t.Fatalf("path = %q, want escaped package path", gotPath)
	}
	if gotBody.Name != "demo/pkg" {
		t.Fatalf("body name = %q, want path package name", gotBody.Name)
	}
	if gotBody.Owner[0] != "alice" || gotBody.Type != "src" {
		t.Fatalf("unexpected body: %+v", gotBody)
	}
}

func TestCreateXflowUsesPostAndActionUserHeader(t *testing.T) {
	var gotMethod, gotPath, gotActionUser string
	var gotBody dto.CreateXflowReq
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotActionUser = r.Header.Get("X-Action-User")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(dto.ApiResponse[int64]{Code: 0, Data: 12345})
	}))
	defer server.Close()

	client := &Client{baseURL: server.URL, http: newHTTPClientForTest()}
	resp, err := client.CreateXflow(context.Background(), dto.CreateXflowReq{
		ActionUser: "operator",
		Type:       "publish",
		Package:    "demo",
		PublishOpts: &dto.PublishOptions{
			VersionNote:    "release",
			ReasonNoChange: "no code change",
			Region:         []string{"cn"},
			Dist:           []string{"Debian 12"},
		},
	})
	if err != nil {
		t.Fatalf("CreateXflow returned error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %q, want %q", gotMethod, http.MethodPost)
	}
	if gotPath != "/xflow" {
		t.Fatalf("path = %q, want /xflow", gotPath)
	}
	if gotActionUser != "operator" {
		t.Fatalf("X-Action-User = %q, want operator", gotActionUser)
	}
	if gotBody.ActionUser != "" {
		t.Fatalf("actionUser should not be sent in JSON body, got %q", gotBody.ActionUser)
	}
	if gotBody.Type != "publish" || gotBody.Package != "demo" {
		t.Fatalf("unexpected body: %+v", gotBody)
	}
	if resp == nil || resp.XflowID != 12345 {
		t.Fatalf("response = %+v, want xflowID 12345", resp)
	}
}

func newHTTPClientForTest() *resty.Client {
	return resty.New()
}
