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
