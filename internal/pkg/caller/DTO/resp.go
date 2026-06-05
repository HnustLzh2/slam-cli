package DTO

type ApiResponse[T any] struct {
	Code  int    `json:"code"`
	Data  T      `json:"data"`
	Error string `json:"error"`
}

type PageInfo struct {
	Page      int64 `json:"page"`
	PageSize  int64 `json:"pageSize"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"totalPage"`
}

type MGetPackageResp struct {
	Packages []SimplePackage `json:"packages"`
	PageInfo *PageInfo       `json:"pageInfo"`
}

type GetAllPackagesResp struct {
	PrivatePackages   []SimplePackage    `json:"privatePackages"`
	ThirdRepoPackages []ThirdRepoPackage `json:"thirdRepoPackages"`
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
