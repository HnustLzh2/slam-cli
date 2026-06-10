package DTO

type PageReq struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}
type MGetPackageReq struct {
	Keywords string `json:"keywords"`
	Owner    string `json:"owner"`
	PageReq
	Dkms  *bool  `json:"dkms"`
	Agent *bool  `json:"agent"`
	Type  string `json:"type"`
}

type SearchLivePatchReq struct {
	PageReq
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	Version       string `json:"version"`
	Distribution  string `json:"distribution"`
	OsVersion     string `json:"os_version"`
	KernelVersion string `json:"kernel_version"`
	Arch          string `json:"arch"`
	Creator       string `json:"creator"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
}

type MGetThirdRepoReq struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Owner string `json:"owner"`
	PageReq
}
