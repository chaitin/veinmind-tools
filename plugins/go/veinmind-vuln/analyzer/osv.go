package analyzer

import (
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/sdk/osv"
)

var AppMap = map[string]string{
	"gobinary":        "Go",
	"gomod":           "Go",
	"python-pkg":      "PyPI",
	"pip":             "PyPI",
	"pipenv":          "PyPI",
	"jar":             "Maven",
	"pom":             "Maven",
	"gradle-lockfile": "Maven",
	"npm":             "npm",
	"node-pkg":        "npm",
	"yarn":            "npm",
	"pnpm":            "npm",
}

func ScanOSV(res *model.ScanResult, verbose bool) {
	bqs := transferQuery(res)
	for _, bq := range bqs {
		resp, err := osv.MakeRequest(bq)
		if err != nil {
			log.Error(err)
			continue
		}
		if verbose {
			hydratedResp, err := osv.Hydrate(resp)
			if err != nil {
				log.Error(err)
				continue
			}
			for hi, hr := range hydratedResp.Results {
				for _, v := range hr.Vulns {
					req := bq.Queries[hi]
					switch req.Type {
					case "os":
						for pi, p := range res.PackageInfos {
							if p.FilePath == req.Package.Path {
								res.PackageInfos[pi].Packages[req.Package.Index].Vulnerabilities = append(res.PackageInfos[pi].Packages[req.Package.Index].Vulnerabilities, v)
								res.CveTotal += 1
								break
							}
						}
					case "app":
						for ai, a := range res.Applications {
							if a.FilePath == req.Package.Path {
								res.Applications[ai].Libraries[req.Package.Index].Vulnerabilities = append(res.Applications[ai].Libraries[req.Package.Index].Vulnerabilities, v)
								res.CveTotal += 1
								break
							}
						}
					}
				}
			}
		} else {
			for hi, hr := range resp.Results {
				for _, v := range hr.Vulns {
					req := bq.Queries[hi]
					switch req.Type {
					case "os":
						for pi, p := range res.PackageInfos {
							if p.FilePath == req.Package.Path {
								res.PackageInfos[pi].Packages[req.Package.Index].Vulnerabilities = append(res.PackageInfos[pi].Packages[req.Package.Index].Vulnerabilities, osv.Vulnerability{
									ID: v.ID,
								})
								res.CveTotal += 1
								break
							}
						}
					case "app":
						for ai, a := range res.Applications {
							if a.FilePath == req.Package.Path {
								res.Applications[ai].Libraries[req.Package.Index].Vulnerabilities = append(res.Applications[ai].Libraries[req.Package.Index].Vulnerabilities, osv.Vulnerability{
									ID: v.ID,
								})
								res.CveTotal += 1
								break
							}
						}
					}
				}
			}
		}
	}
}

func transferQuery(res *model.ScanResult) []osv.BatchedQuery {
	var bquerys []osv.BatchedQuery
	// os package
	// osv only support Debian/Alpine
	if res.OSInfo.Family == "debian" {
		bq := &osv.BatchedQuery{
			Queries: make([]*osv.Query, 0),
		}
		for _, pkgs := range res.PackageInfos {
			if len(pkgs.Packages) == 0 {
				continue
			}
			for _, pkg := range pkgs.Packages {
				if len(bq.Queries) > 1000 {
					bquerys = append(bquerys, *bq)
					bq = &osv.BatchedQuery{
						Queries: make([]*osv.Query, 0),
					}
				}
				bq.Queries = append(bq.Queries, &osv.Query{
					Type:    "os",
					Version: pkg.Version,
					Package: osv.Package{
						Name:      pkg.Name,
						Ecosystem: "Debian",
					},
				})
			}
		}
		bquerys = append(bquerys, *bq)
	} else if res.OSInfo.Family == "alpine" {
		bq := &osv.BatchedQuery{
			Queries: make([]*osv.Query, 0),
		}
		for _, pkgs := range res.PackageInfos {
			if len(pkgs.Packages) == 0 {
				continue
			}
			for index, pkg := range pkgs.Packages {
				if len(bq.Queries) > 1000 {
					bquerys = append(bquerys, *bq)
					bq = &osv.BatchedQuery{
						Queries: make([]*osv.Query, 0),
					}
				}
				bq.Queries = append(bq.Queries, &osv.Query{
					Type:    "os",
					Version: pkg.Version,
					Package: osv.Package{
						Index:     index,
						Path:      pkgs.FilePath,
						Name:      pkg.Name,
						Ecosystem: "Alpine",
					},
				})
			}
		}
		bquerys = append(bquerys, *bq)
	}

	// app
	appBq := &osv.BatchedQuery{
		Queries: make([]*osv.Query, 0),
	}
	for _, apps := range res.Applications {
		if t, ok := AppMap[apps.Type]; ok {
			for index, app := range apps.Libraries {
				if len(appBq.Queries) > 1000 {
					bquerys = append(bquerys, *appBq)
					appBq = &osv.BatchedQuery{
						Queries: make([]*osv.Query, 0),
					}
				}
				appBq.Queries = append(appBq.Queries, &osv.Query{
					Type:    "app",
					Version: app.Version,
					Package: osv.Package{
						Index:     index,
						Path:      apps.FilePath,
						Name:      app.Name,
						Ecosystem: t,
					},
				})
			}
		}
	}
	bquerys = append(bquerys, *appBq)

	return bquerys
}
