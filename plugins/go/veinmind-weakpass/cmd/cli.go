package main

import (
	"os"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-common-go/service/report/event"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/dict"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/dict/embed"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/hash"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/utils"
)

var serviceName []string
var threads int
var username string
var dictpath string

var reportService = &report.Service{}
var rootCmd = &cmd.Command{}
var extractCmd = &cmd.Command{
	Use:   "extract",
	Short: "extract dict file to disk",
	Run: func(cmd *cmd.Command, args []string) {
		embed.ExtractAll()
	},
}
var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "scan mode",
}

var scanImageCmd = &cmd.Command{
	Use:   "image",
	Short: "scan image weakpass",
}

var scanContainerCmd = &cmd.Command{
	Use:   "container",
	Short: "scan container weakpass",
}

func scanImage(c *cmd.Command, image api.Image) (err error) {
	config := model.Config{Thread: threads, Username: username, Dictpath: dictpath}
	for _, service := range serviceName {
		ModuleResult, err := utils.StartModule(config, image, service, map[string]string{
			"module_name": service,
			"image_name": func() string {
				refs, _ := image.RepoRefs()
				if len(refs) > 0 {
					return refs[0]
				}

				return ""
			}(),
		})
		if err != nil {
			log.Error(err)
			continue
		}
		err = utils.GenerateImageReport(ModuleResult, image, reportService)
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return nil
}

func scanContainer(c *cmd.Command, container api.Container) (err error) {
	config := model.Config{Thread: threads, Username: username, Dictpath: dictpath}
	for _, service := range serviceName {
		ModuleResult, err := utils.StartModule(config, container, service, map[string]string{
			"module_name": service,
			"image_name":  "",
		})
		if err != nil {
			log.Error(err)
			continue
		}
		err = utils.GenerateContainerReport(ModuleResult, container, reportService)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	// match environment weak password.
	oci, err := container.OCISpec()
	if err != nil {
		return nil
	}
	p := hash.Plain{}
	if oci.Process != nil && oci.Process.Env != nil {
		for _, env := range oci.Process.Env {
			for _, d := range dict.DictMap["base"] {
				envs := strings.Split(env, "=")
				if len(envs) != 2 {
					continue
				}
				matched, _ := p.Match(envs[1], d)
				if matched {
					err = utils.GenerateContainerReport([]model.WeakpassResult{
						{
							Password:    envs[1],
							ServiceType: event.Env,
							Filepath:    env,
						},
					}, container, reportService)
					if err != nil {
						log.Error(err)
						continue
					}
				}
			}
		}
	}

	return nil
}

func init() {

	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-weakpass",
		Author:      "veinmind-team",
		Description: "veinmind-weakpass scanner image weakpass",
	}))
	scanCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 10, "password brute threads")
	scanCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username e.g. root")
	scanCmd.PersistentFlags().StringVarP(&dictpath, "dictpath", "d", "", "dict path e.g. ./mypass.dict")
	scanCmd.PersistentFlags().StringSliceVarP(&serviceName, "serviceName", "s", []string{"mysql", "tomcat", "redis", "ssh", "ftp"}, "find weakpass in these service e.g. ssh")
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapImageCommand(scanImageCmd, scanImage), reportService))
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapContainerCommand(scanContainerCmd, scanContainer), reportService))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
