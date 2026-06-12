package DTO

import "encoding/json"

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

type FileReq struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type ListenPort struct {
	Protocol string `json:"protocol"`
	Port     int32  `json:"port"`
	Desc     string `json:"desc"`
}

type PublishDocs struct {
	ReviewDoc   string `json:"reviewDoc"`
	ReleaseNote string `json:"releaseNote"`
	TestDoc     string `json:"testDoc"`
}

type UpdatePackageReq struct {
	Name                string       `json:"name,omitempty"`
	Owner               []string     `json:"owner"`
	Type                string       `json:"type"`
	Repo                string       `json:"repo"`
	Branch              string       `json:"branch"`
	GoVersion           string       `json:"goVersion"`
	Files               []FileReq    `json:"files"`
	Region              []string     `json:"region"`
	Dist                []string     `json:"dist"`
	Dkms                bool         `json:"dkms"`
	PublishNotifyChatID string       `json:"publishNotifyChatID"`
	Agent               bool         `json:"agent"`
	Desc                string       `json:"desc"`
	AgentType           string       `json:"agentType"`
	ResourceLimit       string       `json:"resourceLimit"`
	OpenSource          bool         `json:"openSource"`
	CodeLanguage        []string     `json:"codeLanguage"`
	ListenPorts         []ListenPort `json:"listenPorts"`
	InvolvedHardtype    string       `json:"involvedHardtype"`
	HasNDA              bool         `json:"hasNDA"`
}

type CreateXflowReq struct {
	ActionUser   string           `json:"-"`
	Type         string           `json:"type"`
	Package      string           `json:"package"`
	RegisterOpts *RegisterOptions `json:"registerOpts,omitempty"`
	PublishOpts  *PublishOptions  `json:"publishOpts,omitempty"`
	SyncOpts     *SyncOptions     `json:"syncOpts,omitempty"`
}

type RegisterOptions struct {
	PackageType         string       `json:"packageType"`
	Owner               []string     `json:"owner"`
	Repo                string       `json:"repo"`
	Branch              string       `json:"branch"`
	GoVersion           string       `json:"goVersion"`
	Files               []FileReq    `json:"files"`
	Dist                []string     `json:"dist"`
	Region              []string     `json:"region"`
	Dkms                bool         `json:"dkms"`
	Desc                string       `json:"desc"`
	Agent               bool         `json:"agent"`
	AgentType           string       `json:"agentType"`
	ResourceLimit       string       `json:"resourceLimit"`
	PublishNotifyChatID string       `json:"publishNotifyChatID"`
	OpenSource          bool         `json:"openSource"`
	CodeLanguage        []string     `json:"codeLanguage"`
	ListenPorts         []ListenPort `json:"listenPorts"`
	InvolvedHardtype    string       `json:"involvedHardtype"`
	HasNDA              bool         `json:"hasNDA"`
}

type PublishOptions struct {
	DocsV2          *PublishDocs     `json:"docsV2,omitempty"`
	VersionNote     string           `json:"versionNote"`
	SmokeTestParams *SmokeTestParams `json:"smokeTestParams,omitempty"`
	ChangeOrder     json.RawMessage  `json:"changeOrder,omitempty"`
	ReasonNoChange  string           `json:"reasonNoChange"`
	Branch          string           `json:"branch"`
	Region          []string         `json:"region"`
	Dist            []string         `json:"dist"`
	Files           []FileReq        `json:"files"`
	HwKeys          []string         `json:"hwKeys"`
}

type SmokeTestParams struct {
	LinuxKernel []string `json:"linuxKernel"`
}

type SyncOptions struct {
	VersionID int64    `json:"versionID"`
	Region    []string `json:"region"`
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

type MGetTaskReq struct {
	PageReq
	ID      int64  `json:"id"`
	Package string `json:"package"`
	Type    string `json:"type"`
	State   string `json:"state"`
	Source  string `json:"source"`
	Creator string `json:"creator"`
	My      bool   `json:"my"`
}

type MGetXflowReq struct {
	PageReq
	XflowID int64  `json:"xflowID"`
	Package string `json:"package"`
	Type    string `json:"type"`
	State   string `json:"state"`
	Creator string `json:"creator"`
	My      bool   `json:"my"`
}
