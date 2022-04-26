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
)

type AlertType uint32

const (
	Vulnerability AlertType = iota
	MaliciousFile
	Backdoor
	Sensitive
	AbnormalHistory
	Weakpass
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
}

type FileDetail struct {
	Path string      `json:"path"`
	Perm os.FileMode `json:"perm"`
	Size int64       `json:"size"`
	Gid  int64       `json:"gid"`
	Uid  int64       `json:"uid"`
	Ctim int64       `json:"ctim"`
	Mtim int64       `json:"mtim"`
	Atim int64       `json:"atim"`
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
	Description string `json:"description"`
}

type SensitiveEnvDetail struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type HistoryDetail struct {
	Instruction string `json:"instruction"`
	Content     string `json:"content"`
	Description string `json:"description"`
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
