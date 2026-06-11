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

type SearchLivePatchResp struct {
	PageInfo *PageInfo          `json:"page_info"`
	Data     []LivePatchProduct `json:"data"`
}

type MGetThirdRepoResp struct {
	ConfigList []ThirdRepoInfo `json:"configList"`
	PageInfo   *PageInfo       `json:"pageInfo"`
}

type MGetTaskResp struct {
	PageInfo *PageInfo  `json:"pageInfo"`
	Tasks    []TaskInfo `json:"tasks"`
}

type MGetXflowResp struct {
	XflowList []SimpleXflowInfo `json:"xflowList"`
	PageInfo  *PageInfo         `json:"pageInfo"`
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

type ThirdRepoOperations struct {
	Add               string `json:"add"`
	Remove            string `json:"remove"`
	Release           string `json:"release"`
	Update            string `json:"update"`
	Status            string `json:"status"`
	DistributionsOper string `json:"distributionsOper"`
	ComponentsOper    string `json:"componentsOper"`
}

type ThirdRepoInfo struct {
	ID                 int64                `json:"id"`
	CreatedAt          string               `json:"createdAt"`
	UpdatedAt          string               `json:"updatedAt"`
	Name               string               `json:"name"`
	Owner              []string             `json:"owner"`
	Level              int64                `json:"level"`
	Type               string               `json:"type"`
	OperatorWhitelists []string             `json:"operatorWhitelists"`
	URL                string               `json:"url"`
	Distribution       []string             `json:"distribution"`
	Component          []string             `json:"component"`
	Architecture       []string             `json:"architecture"`
	XflowProcess       string               `json:"xflowProcess"`
	Operatings         *ThirdRepoOperations `json:"operatings"`
}

type LivePatchStatusTime struct {
	Pending            string `json:"pending"`
	InProduction       string `json:"in_production"`
	ProductionComplete string `json:"production_complete"`
	WaitTest           string `json:"wait_test"`
	WaitLongTest       string `json:"wait_long_test"`
	Testing            string `json:"testing"`
	LongTesting        string `json:"long_testing"`
	TestCompletion     string `json:"test_completion"`
	LongTestCompletion string `json:"long_test_completion"`
	BuildFailed        string `json:"build_failed"`
	TestFailed         string `json:"test_failed"`
	LongTestFailed     string `json:"long_test_failed"`
	UploadReqComplete  string `json:"upload_req_complete"`
}

type LivePatchProduct struct {
	ID             int64                `json:"id"`
	CreatedAt      string               `json:"created_at"`
	UpdatedAt      string               `json:"updated_at"`
	Name           string               `json:"name"`
	Status         string               `json:"status"`
	Distribution   string               `json:"distribution"`
	Version        string               `json:"version"`
	OsVersion      string               `json:"os_version"`
	KernelVersion  string               `json:"kernel_version"`
	Arch           string               `json:"arch"`
	Creator        string               `json:"creator"`
	PatchNumber    int64                `json:"patch_number"`
	TestNumber     int64                `json:"test_number"`
	LongTestNumber int64                `json:"long_test_number"`
	FileLink       string               `json:"file_link"`
	PatchLink      string               `json:"patch_link"`
	PatchMD5       string               `json:"patch_md5"`
	Log            string               `json:"log"`
	TestLog        string               `json:"test_log"`
	LongTestLog    string               `json:"long_test_log"`
	Comment        string               `json:"comment"`
	StatusTime     *LivePatchStatusTime `json:"status_time"`
	UploadXflowID  int64                `json:"upload_xflow_id"`
}

type TaskInfo struct {
	ID                   int64                 `json:"id"`
	Package              string                `json:"package"`
	Type                 string                `json:"type"`
	State                string                `json:"state"`
	Source               string                `json:"source"`
	Creator              string                `json:"creator"`
	URL                  []string              `json:"url"`
	BuildTaskOpts        *BuildTaskOpts        `json:"buildTaskOpts,omitempty"`
	SmokeTestTaskOpts    *SmokeTestTaskOpts    `json:"smokeTestTaskOpts,omitempty"`
	PackageCheckTaskOpts *PackageCheckTaskOpts `json:"packageCheckTaskOpts,omitempty"`
	CreatedAt            string                `json:"createdAt"`
	UpdatedAt            string                `json:"updatedAt"`
	Error                string                `json:"error"`
	PackageInfo          *SimplePackage        `json:"packageInfo"`
}

type BuildTaskOpts struct {
	Branch string   `json:"branch"`
	Dist   []string `json:"dist"`
}

type SmokeTestTaskOpts struct {
	LinuxKernel []string `json:"linuxKernel"`
	Dist        []string `json:"dist"`
	TosObjName  []string `json:"tosObjName,omitempty"`
}

type PackageCheckTaskOpts struct {
	Dist        []string `json:"dist"`
	Region      []string `json:"region"`
	LinuxKernel []string `json:"linuxKernel"`
}

type SimpleXflowInfo struct {
	XflowID   int64    `json:"xflowID"`
	Type      string   `json:"type"`
	Package   string   `json:"package"`
	Applicant string   `json:"applicant"`
	Reviewers []string `json:"reviewers"`
	State     string   `json:"state"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
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
