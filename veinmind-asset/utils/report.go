package utils

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/chaitin/veinmind-tools/veinmind-asset/model"
)

const (
	JSON   = "json"
	CSV    = "csv"
	STDOUT = "stdout"
)

var tab = 24 // stdout输出视图大小控制
var tabwTitle = tabwriter.NewWriter(os.Stdout, 3+(tab+2)*4, 3+(tab+2)*4, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
var tabwBody = tabwriter.NewWriter(os.Stdout, tab+2, (tab+2)*4, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)

type StdStream func() error

func OutputStream(spend time.Duration, results []model.ScanImageResult, stdstream StdStream) error {
	fmt.Fprintln(tabwTitle, "#", Repeat("=", (tab+2)*4), "#")
	fmt.Fprintln(tabwTitle, "| Scan Image Total: ", strconv.Itoa(len(results)), "\t")
	fmt.Fprintln(tabwTitle, "| Spend Time: ", spend.String(), "\t")
	fmt.Fprintln(tabwTitle, "+", Repeat("-", (tab+2)*4), "+")

	err := stdstream()
	if err != nil {
		return err
	}

	fmt.Fprintln(tabwBody, "#", Repeat("=", (tab+2)*4), "#")
	tabwBody.Flush()
	return nil
}

func OutputStdout(verbose bool, pkgType string, results []model.ScanImageResult) error {
	fmt.Fprintln(tabwBody, "| Image ID", "\t", "Image Name", "\t", "Package Total", "\t", "Application Total", "\t")
	for _, r := range results {
		fmt.Fprintln(tabwBody, "|", strings.Replace(r.ImageID, "sha256:", "", -1)[0:12], "\t", Limit(r.ImageName, tab), "\t", strconv.Itoa(r.PackageTotal), "\t", strconv.Itoa(r.ApplicationTotal), "\t")
		if verbose {
			fmt.Fprintln(tabwBody, "+", Repeat("-", (tab+2)*4), "+")
			fmt.Fprintln(tabwBody, "| Package Type", "\t", "Package Name", "\t", "Package Version", "\t", "Package File", "\t")
			if pkgType == "all" || pkgType == "os" {
				for _, pkg := range r.Packages {
					fmt.Fprintln(tabwBody, "|", "os-pkg", "\t", Limit(pkg.Name, tab), "\t", Limit(pkg.Version, tab), "\t", Limit(pkg.FilePath, tab), "\t")
				}
			}
			for _, info := range r.Applications {
				if pkgType == "all" || pkgType == info.Type || strings.Contains(info.Type, pkgType) {
					for _, lib := range info.Libraries {
						fmt.Fprintln(tabwBody, "|", Limit(info.Type, tab), "\t", Limit(lib.Name, tab), "\t", Limit(lib.Version, tab), "\t", Limit(func() string {
							if lib.FilePath != "" {
								return lib.FilePath
							} else {
								return info.FilePath
							}
						}(), tab), "\t")
					}
				}
			}
		}
	}
	return nil
}

func OutputJSON(res []model.ScanImageResult) error {
	var jsonFile *os.File
	var name = "report.json"
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

	scanResultJson, err := json.Marshal(res)
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(scanResultJson)
	if err != nil {
		return err
	}

	tabw := tabwriter.NewWriter(os.Stdout, 3+(tab+2)*4, 3+(tab+2)*4, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
	out := "| The Detail Results Export report.json"
	fmt.Fprintln(tabw, out, Repeat(" ", (tab+2)*4-len(out)+1), "|")
	return nil
}

func OutputCSV(res []model.ScanImageResult) error {
	var csvFile *os.File
	var name = "report.csv"
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

	w.Write([]string{"ImageID", "ImageName", "OSFamily", "OSName", "Type", "Name", "Version", "FilePath"})
	for _, r := range res {
		for _, pkg := range r.Packages {
			err := w.Write([]string{r.ImageID, r.ImageName, r.ImageInfo.Family, r.ImageInfo.Name, "os-pkg", pkg.Name, pkg.Version, pkg.FilePath})
			if err != nil {
				return err
			}
		}
		for _, info := range r.Applications {
			for _, lib := range info.Libraries {
				err := w.Write([]string{r.ImageID, r.ImageName, r.ImageInfo.Family, r.ImageInfo.Name, info.Type, lib.Name, lib.Version, lib.FilePath})
				if err != nil {
					return err
				}
			}
		}
	}
	w.Flush()

	tabw := tabwriter.NewWriter(os.Stdout, 3+(tab+2)*4, 3+(tab+2)*4, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
	out := "| The Detail Results Export report.csv"
	fmt.Fprintln(tabw, out, Repeat(" ", (tab+2)*4-len(out)+1), "|")
	return nil
}
