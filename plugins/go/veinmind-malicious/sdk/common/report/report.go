package report

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/embed"
	"html/template"
	"os"
)

const (
	HTML = "html"
	JSON = "json"
	CSV  = "csv"
)

func OutputHTML(scanResult model.ReportData, name string) error {
	tpl, err := template.ParseFS(embed.EmbedFile, "template/template.html")

	if err != nil {
		return err
	}

	var tplFile *os.File
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		tplFile, err = os.Create(name)
		if err != nil {
			return err
		}
	} else {
		err := os.Remove(name)
		if err != nil {
			return err
		}

		tplFile, err = os.Create(name)
		if err != nil {
			return err
		}
	}

	err = tpl.Execute(tplFile, scanResult)
	if err != nil {
		return err
	}

	return nil
}

func OutputJSON(scanResult model.ReportData, name string) error {
	var jsonFile *os.File
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		jsonFile, err = os.Create(name)
		if err != nil {
			return err
		}
	} else {
		err := os.Remove(name)
		if err != nil {
			return err
		}

		jsonFile, err = os.Create(name)
		if err != nil {
			return err
		}
	}

	scanResultJson, err := json.Marshal(scanResult)
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(scanResultJson)
	if err != nil {
		return err
	}

	return nil
}

func OutputCSV(scanResult model.ReportData, name string) error {
	var csvFile *os.File
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		csvFile, err = os.Create(name)
		if err != nil {
			return err
		}
	} else {
		err := os.Remove(name)
		if err != nil {
			return err
		}

		csvFile, err = os.Create(name)
		if err != nil {
			return err
		}
	}

	w := csv.NewWriter(csvFile)
	//TODO: 处理CSV格式
	w.Write([]string{"ImageID", "LayerID", "RelativePath", "FileName", "FileSize", "FileMd5", "FileSha256", "FileCreated", "Description"})
	for _, i := range scanResult.ScanImageResult {
		for _, l := range i.Layers {
			for _, r := range l.MaliciousFileInfos {
				err := w.Write([]string{l.ImageID, l.LayerID, r.RelativePath, r.FileName, r.FileSize, r.FileMd5, r.FileSha256, r.FileCreated, r.Description})
				if err != nil {
					return err
				}
			}
		}
	}
	w.Flush()

	return nil
}

func CalculateScanReportCount(scanResult *model.ReportData) {
	for _, i := range scanResult.ScanImageResult {
		scanResult.ScanImageCount++
		scanResult.MaliciousFileCount = scanResult.MaliciousFileCount + i.MaliciousFileCount
		scanResult.ScanFileCount = scanResult.ScanFileCount + i.ScanFileCount
	}
}

func SortScanReport(scanResult *model.ReportData) {
	// 将恶意样本数量高的报告放在前面(冒泡排序)
	scanImageResultLength := len(scanResult.ScanImageResult)
	for i := 0; i < scanImageResultLength; i++ {
		for j := i; j < scanImageResultLength; j++ {
			if scanResult.ScanImageResult[i].MaliciousFileCount < scanResult.ScanImageResult[j].MaliciousFileCount {
				scanImageResultTemp := scanResult.ScanImageResult[i]
				scanResult.ScanImageResult[i] = scanResult.ScanImageResult[j]
				scanResult.ScanImageResult[j] = scanImageResultTemp
			}
		}
	}
}
