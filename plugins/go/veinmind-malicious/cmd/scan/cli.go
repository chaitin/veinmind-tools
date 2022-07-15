package main

import (
	"github.com/biter777/processex"
	"github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/config"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database/model"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/embed"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/scanner/malicious"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/sdk/av/clamav"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/sdk/common/report"
	reportService "github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	ps "github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/net"
	"github.com/spf13/cobra"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var reportData = model.ReportData{}
var reportLock sync.Mutex
var scanStart = time.Now()
var clamdPid = 0

var rootCmd = &cmd.Command{}
var extractCmd = &cmd.Command{
	Use:   "extract",
	Short: "Extract config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		embed.ExtractAll()
		return nil
	},
}
var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "Scan image malicious files",
	PreRun: func(cmd *cobra.Command, args []string) {
		clamavManualStart, err := cmd.Flags().GetBool("clamav-manually-start")
		if err != nil {
			log.Error(err)
		}
		clamavConf, err := cmd.Flags().GetString("clamav-conf")
		if err != nil {
			log.Error(err)
		}
		clamavExec, err := cmd.Flags().GetString("clamav-exec")
		if err != nil {
			log.Error(err)
		}
		port, err := cmd.Flags().GetString("clamav-port")
		if err != nil {
			log.Error(err.Error())
		}
		host, err := cmd.Flags().GetString("clamav-addr")
		if err != nil {
			log.Error(err)
		}

		// the flag of Manual run the clamAV
		if !clamavManualStart {
			if isLocalHost(host) {
				log.Info("automatically start clamAV")
				// make sure local is not running a clamAV
				_, process, err := processex.FindByName("clamd")
				if err == processex.ErrNotFound {
					log.Info("start local clamAV")
					clam := exec.Command("/bin/sh", "-c", clamavExec+" --config-file="+clamavConf)
					err := clam.Run()
					if err != nil {
						log.Error("clamAV failed to start: ", err)
						return
					}
					//check the server is really running
					clamdPid, err = getPidByPort(port)
					if err != nil {
						log.Error(err)
					}
					if clamdPid == 0 {
						log.Error("this port is not open")
					} else {
						log.Info("the clamAV is running at port: ", port, ", the pid is ", clamdPid)
					}
				} else if err == nil {
					log.Info("the local clamAV is running")
					clamdListenPid, err := getPidByPort(port)
					if err != nil {
						log.Error(err)
					}
					if clamdListenPid == 0 {
						log.Error("the port is not open, maybe the port is error")
						return
					}

					// make sure the port is running clamAV
					for i := 0; i < len(process); i++ {
						if process[i].PID == clamdListenPid {
							log.Info("working Pid is ", process[i].Pid)
							return
						}
					}
					log.Error("this port is not listen by clamAV")
				} else {
					log.Error("find clamAV ERROR: ", err)
				}
			} else {
				log.Error("the host is not local, can't not start it")
			}
		}
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		// 计算扫描数据
		spend := time.Since(scanStart)
		reportData.ScanSpendTime = spend.String()
		reportData.ScanStartTime = scanStart.Format("2006-01-02 15:04:05")
		report.CalculateScanReportCount(&reportData)
		report.SortScanReport(&reportData)

		format, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Error(err)
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error(err)
		}
		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			log.Error(err)
		}
		name = strings.Join([]string{name, format}, ".")
		fpath := path.Join(outputPath, name)
		port, err := cmd.Flags().GetString("clamav-port")
		if err != nil {
			log.Error(err)
		}

		switch format {
		case report.HTML:
			report.OutputHTML(reportData, fpath)
		case report.JSON:
			report.OutputJSON(reportData, fpath)
		case report.CSV:
			report.OutputCSV(reportData, fpath)
		}

		if clamdPid != 0 {
			log.Info("close clamAV server")
			runFinishPid, err := getPidByPort(port)
			if err != nil {
				log.Error(err)
			}
			//make sure the pid run in port equal start pid
			//log.Info(runFinishPid)
			if runFinishPid == 0 {
				log.Error("the port is close")
			} else {
				if runFinishPid == clamdPid {
					proc, err := ps.FindProcess(clamdPid)
					if err != nil {
						log.Error(err)
					}
					if proc.Executable() == "clamd" {
						//log.Info(clamdPid)
						err := syscall.Kill(-clamdPid, syscall.SIGKILL)
						if err != nil {
							log.Error(err)
						}
					} else {
						log.Error("the pid ", clamdPid, " is not belong to clamAV")
					}
				} else {
					log.Error("the listener of the port: ", port, " has been change")
				}
			}
		}
	},
}

func scan(c *cmd.Command, image api.Image) error {
	//default clamAV host and port

	clamavHost, err := c.Flags().GetString("clamav-addr")
	if err != nil {
		log.Error(err)
	}
	clamavPort, err := c.Flags().GetString("clamav-port")
	if err != nil {
		log.Error(err)
	}

	antiVirusAgent := malicious.AntiVirusService{ClamavAgent: clamav.New(clamavHost, clamavPort)}

	result, err := malicious.Scan(image, antiVirusAgent)
	if err != nil {
		log.Error(err)
		return nil
	}

	reportLock.Lock()
	reportData.ScanImageResult = append(reportData.ScanImageResult, result)
	reportLock.Unlock()

	// result event
	if result.MaliciousFileCount > 0 {
		details := []reportService.AlertDetail{}
		for _, l := range result.Layers {
			if len(l.MaliciousFileInfos) > 0 {
				for _, mr := range l.MaliciousFileInfos {
					f, err := image.Open(mr.RelativePath)
					if err != nil {
						log.Error(err)
						continue
					}

					fStat, err := f.Stat()
					if err != nil {
						log.Error(err)
						continue
					}
					fSys := fStat.Sys().(*syscall.Stat_t)

					details = append(details, reportService.AlertDetail{
						MaliciousFileDetail: &reportService.MaliciousFileDetail{
							Engine:        mr.Engine,
							MaliciousName: mr.Description,
							FileDetail: reportService.FileDetail{
								Path: mr.RelativePath,
								Perm: fStat.Mode(),
								Size: fStat.Size(),
								Gid:  int64(fSys.Gid),
								Uid:  int64(fSys.Uid),
								Ctim: fSys.Ctim.Sec,
								Mtim: fSys.Mtim.Sec,
								Atim: fSys.Atim.Sec,
							},
						},
					})
				}
			}
		}
		reportEvent := reportService.ReportEvent{
			ID:           image.ID(),
			Level:        reportService.High,
			DetectType:   reportService.Image,
			EventType:    reportService.Risk,
			AlertType:    reportService.MaliciousFile,
			AlertDetails: details,
		}
		err = reportService.DefaultReportClient().Report(reportEvent)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPidByPort(port string) (int, error) {
	netConnections, err := net.Connections("tcp")
	if err != nil {
		return 0, err
	}
	for idx := 0; idx < len(netConnections); idx++ {
		if strconv.Itoa(int(netConnections[idx].Laddr.Port)) == port {
			return int(netConnections[idx].Pid), nil
		}
	}
	return 0, nil
}

func isLocalHost(host string) bool {
	localIPs, err := getLocalIP()
	if err != nil {
		log.Error(err)
		return false
	}
	for i := 0; i < len(localIPs); i++ {
		if host == localIPs[i] {
			return true
		}
	}
	return false
}

func getLocalIP() ([]string, error) {
	addresses, err := net.Interfaces()
	partIP := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regIP := partIP + "\\." + partIP + "\\." + partIP + "\\." + partIP
	matchIP := regexp.MustCompile(regIP)
	if err != nil {
		return nil, err
	}
	IPs := make([]string, 0)
	for _, address := range addresses {
		for i := 0; i < len(address.Addrs); i++ {
			IPs = append(IPs, matchIP.FindString(address.Addrs[i].Addr))
		}
	}
	return IPs, nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanCmd, scan))
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-malicious",
		Author:      "veinmind-team",
		Description: "veinmind-malicious scanner image malicious file",
	}))
	scanCmd.Flags().StringP("format", "f", "html", "report format for scan report")
	scanCmd.Flags().StringP("name", "n", "report", "report name for scan report")
	scanCmd.Flags().StringP("output", "o", ".", "output path for report")
	scanCmd.Flags().BoolP("clamav-manually-start", "s", true, "whether need to manually start clamAV")
	scanCmd.Flags().StringP("clamav-addr", "a", "127.0.0.1", "host of ClamAV")
	scanCmd.Flags().StringP("clamav-port", "p", "3310", "port of ClamAV")
	scanCmd.Flags().StringP("clamav-exec", "e", "/usr/sbin/clamd", "execution file path of ClamAV")
	scanCmd.Flags().StringP("clamav-conf", "c", "/etc/clamav/clamd.conf ", "config file path of ClamAV")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
