// Package report provides report service for veinmind-runner
// and veinmind-plugin
package report

import (
	"os"
	"time"
)

type Level uint32

const (
	Low Level = iota
	Medium
	High
	Critical
	None
)

type DetectType uint32

const (
	Image DetectType = iota
	Container
)

type EventType uint32

const (
	Risk EventType = iota
	Invasion
	Info
)

type AlertType uint32

const (
	Vulnerability AlertType = iota
	MaliciousFile
	Backdoor
	Sensitive
	AbnormalHistory
	Weakpass
	Asset
	Basic
)

type WeakpassService uint32

const (
	SSH WeakpassService = iota
)

type AlertDetail struct {
	MaliciousFileDetail *MaliciousFileDetail `json:"malicious_file_detail,omitempty"`
	WeakpassDetail      *WeakpassDetail      `json:"weakpass_detail,omitempty"`
	BackdoorDetail      *BackdoorDetail      `json:"backdoor_detail,omitempty"`
	SensitiveFileDetail *SensitveFileDetail  `json:"sensitive_file_detail,omitempty"`
	SensitiveEnvDetail  *SensitiveEnvDetail  `json:"sensitive_env_detail,omitempty"`
	HistoryDetail       *HistoryDetail       `json:"history_detail,omitempty"`
	AssetDetail         *AssetDetail         `json:"asset_detail,omitempty"`
	BasicDetail         *BasicDetail         `json:"basic_detail,omitempty"`
}

type FileDetail struct {
	Path  string      `json:"path"`
	Perm  os.FileMode `json:"perm"`
	Size  int64       `json:"size"`
	Gname string      `json:"gname"`
	Gid   int64       `json:"gid"`
	Uid   int64       `json:"uid"`
	Uname string      `json:"uname"`
	Ctim  int64       `json:"ctim"`
	Mtim  int64       `json:"mtim"`
	Atim  int64       `json:"atim"`
}

type MaliciousFileDetail struct {
	FileDetail
	Engine        string `json:"engine"`
	MaliciousType string `json:"malicious_type"`
	MaliciousName string `json:"malicious_name"`
}

type WeakpassDetail struct {
	Username string          `json:"username"`
	Password string          `json:"password"`
	Service  WeakpassService `json:"service"`
}

type BackdoorDetail struct {
	FileDetail
	Description string `json:"description"`
}

type SensitveFileDetail struct {
	FileDetail
	RuleID          int64 `json:"rule_id"`
	RuleName        string `json:"rule_name"`
	RuleDescription string `json:"rule_description"`
}

type SensitiveEnvDetail struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	RuleID          int64 `json:"rule_id"`
	RuleName        string `json:"rule_name"`
	RuleDescription string `json:"rule_description"`
}

type HistoryDetail struct {
	Instruction string `json:"instruction"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

type AssetDetail struct {
	OS           AssetOSDetail             `json:"os"`
	PackageInfos []AssetPackageDetails     `json:"package_infos"`
	Applications []AssetApplicationDetails `json:"applications"`
}

type AssetOSDetail struct {
	Family string `json:"family"`
	Name   string `json:"name"`
	Eosl   bool   `json:"EOSL,omitempty"`
}

type AssetPackageDetails struct {
	FilePath string               `json:"file_path"`
	Packages []AssetPackageDetail `json:"packages"`
}

type AssetApplicationDetails struct {
	Type     string               `json:"type"`
	FilePath string               `json:"file_path,omitempty"`
	Packages []AssetPackageDetail `json:"packages"`
}

type AssetPackageDetail struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Release         string `json:"release"`
	Epoch           int    `json:"epoch"`
	Arch            string `json:"arch"`
	SrcName         string `json:"srcName"`
	SrcVersion      string `json:"srcVersion"`
	SrcRelease      string `json:"srcRelease"`
	SrcEpoch        int    `json:"srcEpoch"`
	Modularitylabel string `json:"modularitylabel"`
	Indirect        bool   `json:"indirect"`
	License         string `json:"license"`
	Layer           string `json:"layer"`
}

type BasicDetail struct {
	References  []string `json:"references"`
	CreatedTime int64    `json:"created_time"`
	Env         []string `json:"env"`
	Entrypoint  []string `json:"entrypoint"`
	Cmd         []string `json:"cmd"`
	WorkingDir  string   `json:"working_dir"`
	Author      string   `json:"author"`
}

type ReportEvent struct {
	ID           string        `json:"id"`
	Time         time.Time     `json:"time"`
	Level        Level         `json:"level"`
	DetectType   DetectType    `json:"detect_type"`
	EventType    EventType     `json:"event_type"`
	AlertType    AlertType     `json:"alert_type"`
	AlertDetails []AlertDetail `json:"alert_details"`
}
