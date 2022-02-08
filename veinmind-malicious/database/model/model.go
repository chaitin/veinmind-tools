package model

import (
	"gorm.io/gorm"
)

type MaliciousFileInfo struct {
	gorm.Model
	ImageID      string
	LayerID      string
	RelativePath string
	FileName     string
	FileSize     string
	FileMd5      string
	FileSha256   string
	FileCreated  string
	Description  string
}

type ReportData struct {
	gorm.Model
	ScanImageCount     int64
	MaliciousFileCount int64
	ScanSpendTime      string
	ScanStartTime      string
	ScanFileCount      int
	ScanImageResult    []ReportImage `gorm:"foreignKey:ImageID"`
}

type ReportImage struct {
	gorm.Model
	EngineType         string
	ImageName          string
	ImageID            string
	MaliciousFileCount int64
	ScanFileCount      int
	ImageCreatedAt     string
	MaliciousFileInfos []MaliciousFileInfo `gorm:"foreignKey:ImageID;references:ImageID"`
	Layers             []ReportLayer       `gorm:"foreignKey:ImageID;references:ImageID"`
}

type ReportLayer struct {
	gorm.Model
	ImageID            string
	LayerID            string
	MaliciousFileInfos []MaliciousFileInfo `gorm:"foreignKey:LayerID;references:LayerID"`
}

func (self *ReportImage) IsMalicious() bool {
	if self.MaliciousFileCount > 0 {
		return true
	} else {
		return false
	}
}
