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

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/model"
)

const (
	JSON   = "json"
	CSV    = "csv"
	STDOUT = "stdout"
)

var tab = 24 // stdout输出视图大小控制
var col = 4
var tabwTitle = tabwriter.NewWriter(os.Stdout, 3+(tab+2)*col, 3+(tab+2)*col, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
var tabwBody = tabwriter.NewWriter(os.Stdout, tab+2, (tab+2)*col, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)

type StdStream func() error

func OutputStream(spend time.Duration, results []model.ScanResult, stdstream StdStream, asset bool) error {
	log.Info("For more vulnerabilities info, Use `-f json` and see details in `./data/report.json`")
	if !asset {
		col = 5
	}
	vulTotal := 0
	for _, r := range results {
		if r.CveTotal > 0 {
			vulTotal += r.CveTotal
		}
	}

	fmt.Fprintln(tabwTitle, "#", Repeat("=", (tab+2)*col), "#")
	fmt.Fprintln(tabwTitle, "| Scan Image Total: ", strconv.Itoa(len(results)), Repeat(" ", (tab+2)*(col-1)), "\t")
	fmt.Fprintln(tabwTitle, "| Spend Time: ", spend.String(), Repeat(" ", (tab+2)*(col-1)), "\t")
	fmt.Fprintln(tabwTitle, "+", Repeat("-", (tab+2)*col), "+")

	err := stdstream()
	if err != nil {
		return err
	}

	fmt.Fprintln(tabwBody, "#", Repeat("=", (tab+2)*col), "#")
	tabwBody.Flush()
	return nil
}

func OutputStdout(verbose bool, asset bool, pkgType string, results []model.ScanResult) error {
	for index, r := range results {
		if index != 0 {
			fmt.Fprintln(tabwBody, "+", Repeat("-", (tab+2)*col), "+")
		}
		if !asset {
			fmt.Fprintln(tabwBody, "| Image ID", "\t", "Image Name", "\t", "Package Total", "\t", "Application Total", "\t", "Vulns Total", "\t")
			fmt.Fprintln(tabwBody, "|", strings.Replace(r.ID, "sha256:", "", -1)[0:12], "\t", Limit(r.Name, tab), "\t", strconv.Itoa(r.PackageTotal), "\t", strconv.Itoa(r.ApplicationTotal), "\t", strconv.Itoa(r.CveTotal), "\t")
			if verbose {
				fmt.Fprintln(tabwBody, "+", Repeat("-", (tab+2)*col), "+")
				fmt.Fprintln(tabwBody, "| Package Type", "\t", "Package Name", "\t", "Package Version", "\t", "Package File", "\t", "Vulnerable", "\t")
				fmt.Fprintln(tabwBody, "|", Repeat("-", 1+(tab+2)*col), "|")
				if pkgType == "all" || pkgType == "os" {
					for _, pkgInfo := range r.PackageInfos {
						for pi, pkg := range pkgInfo.Packages {
							if len(pkg.Vulnerabilities) > 0 {
								for vi, vul := range pkg.Vulnerabilities {
									if vi == 0 {
										fmt.Fprintln(tabwBody, "|", "os-pkg", "\t", Limit(pkg.Name, tab), "\t", Limit(pkg.Version, tab), "\t", Limit(pkg.FilePath, tab), "\t", vul.GetAliases(), "\t")
									} else {
										fmt.Fprintln(tabwBody, "|", " ", "\t", " ", "\t", " ", "\t", " ", "\t", vul.GetAliases(), "\t")
									}
								}
							} else {
								fmt.Fprintln(tabwBody, "|", "os-pkg", "\t", Limit(pkg.Name, tab), "\t", Limit(pkg.Version, tab), "\t", Limit(pkg.FilePath, tab), "\t", " ", "\t")
							}
							if pi != len(pkgInfo.Packages)-1 {
								fmt.Fprintln(tabwBody, "|", Repeat("-", 1+(tab+2)*col), "|")
							}
						}
					}
				}
				for _, info := range r.Applications {
					if pkgType == "all" || pkgType == info.Type || strings.Contains(info.Type, pkgType) {
						for li, lib := range info.Libraries {
							if len(lib.Vulnerabilities) > 0 {
								for vi, vul := range lib.Vulnerabilities {
									if vi == 0 {
										fmt.Fprintln(tabwBody, "|", Limit(info.Type, tab), "\t", Limit(lib.Name, tab), "\t", Limit(lib.Version, tab), "\t", Limit(func() string {
											if lib.FilePath != "" {
												return lib.FilePath
											} else {
												return info.FilePath
											}
										}(), tab), "\t", vul.Aliases, "\t")
									} else {
										fmt.Fprintln(tabwBody, "|", " ", "\t", " ", "\t", " ", "\t", " ", "\t", vul.Aliases, "\t")
									}
								}
							} else {
								fmt.Fprintln(tabwBody, "|", Limit(info.Type, tab), "\t", Limit(lib.Name, tab), "\t", Limit(lib.Version, tab), "\t", Limit(func() string {
									if lib.FilePath != "" {
										return lib.FilePath
									} else {
										return info.FilePath
									}
								}(), tab), "\t", " ", "\t")
							}
							if li != len(info.Libraries)-1 {
								fmt.Fprintln(tabwBody, "|", Repeat("-", 1+(tab+2)*col), "|")
							}
						}
					}
				}
			}
		} else {
			fmt.Fprintln(tabwBody, "| Image ID", "\t", "Image Name", "\t", "Package Total", "\t", "Application Total", "\t")
			fmt.Fprintln(tabwBody, "|", strings.Replace(r.ID, "sha256:", "", -1)[0:12], "\t", Limit(r.Name, tab), "\t", strconv.Itoa(r.PackageTotal), "\t", strconv.Itoa(r.ApplicationTotal), "\t")
			if verbose {
				fmt.Fprintln(tabwBody, "+", Repeat("-", (tab+2)*col), "+")
				fmt.Fprintln(tabwBody, "| Package Type", "\t", "Package Name", "\t", "Package Version", "\t", "Package File", "\t")
				if pkgType == "all" || pkgType == "os" {
					for _, pkgInfo := range r.PackageInfos {
						for _, pkg := range pkgInfo.Packages {
							fmt.Fprintln(tabwBody, "|", "os-pkg", "\t", Limit(pkg.Name, tab), "\t", Limit(pkg.Version, tab), "\t", Limit(pkg.FilePath, tab), "\t")
						}
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
	}
	return nil
}

func OutputJSON(res []model.ScanResult) error {
	var jsonFile *os.File
	var name = "./data/report.json"
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

	tabw := tabwriter.NewWriter(os.Stdout, 3+(tab+2)*col, 3+(tab+2)*col, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
	out := "| The Detail Results Export report.json"
	fmt.Fprintln(tabw, out, Repeat(" ", (tab+2)*col-len(out)+1), "|")
	return nil
}

func OutputCSV(res []model.ScanResult) error {
	var csvFile *os.File
	var name = "./data/report.csv"
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

	w.Write([]string{"ImageID", "ImageName", "OSFamily", "OSName", "Type", "Name", "Version", "FilePath", "CVE"})
	for _, r := range res {
		for _, pkgInfo := range r.PackageInfos {
			for _, pkg := range pkgInfo.Packages {
				err := w.Write([]string{r.ID, r.Name, r.OSInfo.Family, r.OSInfo.Name, "os-pkg", pkg.Name, pkg.Version, pkg.FilePath, func() string {
					vulns := ""
					for vi, v := range pkg.Vulnerabilities {
						for di, d := range v.GetAliases() {
							vulns += d
							if di != len(v.GetAliases())-1 {
								vulns += "|"
							}
						}
						if vi != len(pkg.Vulnerabilities)-1 {
							vulns += "|"
						}
					}
					return vulns
				}()})
				if err != nil {
					return err
				}
			}
		}
		for _, info := range r.Applications {
			for _, lib := range info.Libraries {
				err := w.Write([]string{r.ID, r.Name, r.OSInfo.Family, r.OSInfo.Name, info.Type, lib.Name, lib.Version, lib.FilePath, func() string {
					vulns := ""
					for _, v := range lib.Vulnerabilities {
						if len(v.Aliases) > 0 {
							vulns += strings.Join(v.Aliases, "|")
						} else {
							vulns = strings.Join([]string{vulns, v.ID}, "|")
						}
					}
					return vulns
				}()})
				if err != nil {
					return err
				}
			}
		}
	}
	w.Flush()

	tabw := tabwriter.NewWriter(os.Stdout, 3+(tab+2)*col, 3+(tab+2)*col, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
	out := "| The Detail Results Export report.csv"
	fmt.Fprintln(tabw, out, Repeat(" ", (tab+2)*col-len(out)+1), "|")
	return nil
}
