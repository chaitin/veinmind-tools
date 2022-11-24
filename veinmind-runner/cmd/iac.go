package main

import (
	"context"
	"github.com/chaitin/libveinmind/go/cmd"
	iacApi "github.com/chaitin/libveinmind/go/iac"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/git"
	"os"
	"path"
	"regexp"
)

var scanIaCCmd = &cmd.Command{
	Use:   "scan-iac",
	Short: "perform iac file scan",
}

var scanLocalCmd = &cmd.Command{
	Use:      "local",
	Short:    "perform local iac file scan",
	PreRunE:  scanPreRunE,
	PostRunE: scanPostRunE,
}

var scanGitRepoCmd = &cmd.Command{
	Use:      "git",
	Short:    "perform git repo iac file scan",
	PreRunE:  scanPreRunE,
	RunE:     scanGitRepoIaCFile,
	PostRunE: scanPostRunE,
}

func scanIaCFile(c *cmd.Command, iac iacApi.IAC) error {
	log.Infof("scan Iac: %s, filetype: %s", iac.Path, iac.Type)
	if err := cmd.ScanIAC(ctx, ps, iac,
		plugin.WithExecInterceptor(func(
			ctx context.Context, plug *plugin.Plugin, c *plugin.Command, next func(context.Context, ...plugin.ExecOption) error,
		) error {
			// Register Service
			reg := service.NewRegistry()
			reg.AddServices(log.WithFields(log.Fields{
				"plugin":  plug.Name,
				"command": path.Join(c.Path...),
			}))
			reg.AddServices(reportService)

			// Next Plugin
			return next(ctx, reg.Bind())
		})); err != nil {
		return err
	}
	return nil
}

func scanGitRepoIaCFile(c *cmd.Command, args []string) error {
	key, err := c.Flags().GetString("ssh-pubkey")
	if err != nil {
		return err
	}

	insecure, err := c.Flags().GetBool("insecure-skip")
	if err != nil {
		return err
	}

	proxy, err := c.Flags().GetString("proxy")
	if err != nil {
		return err
	}

	if proxy != "" {
		os.Setenv("ALL_PROXY", proxy)
	}

	for _, arg := range args {
		isGitUrl, err := regexp.MatchString("^(http(s)?://([^/]+?/){2}|git@[^:]+:[^/]+?/).*?.git$", arg)
		if err != nil {
			continue
		}
		if isGitUrl {
			func() {
				var opt []iacApi.DiscoverOption
				tempDir := path.Join("/tmp", git.RandStr(12))
				err = git.Clone(tempDir, arg, key, insecure)
				if err != nil {
					log.Errorf("git download failed: %s", err)
					// nil point fix
					return
				}
				defer os.RemoveAll(tempDir)
				discovered, err := iacApi.DiscoverIACs(tempDir, opt...)
				if err != nil {
					log.Errorf("git discovered failed: %s", err)
				}
				for _, iac := range discovered {
					_ = scanIaCFile(c, iac)
				}
			}()
		}
	}

	return nil
}

func init() {
	scanIaCCmd.AddCommand(cmd.MapIACCommand(scanLocalCmd, scanIaCFile))
	scanIaCCmd.AddCommand(scanGitRepoCmd)

	scanIaCCmd.PersistentFlags().String("iac-type", "", "dedicate iac type for iac files")
	scanIaCCmd.PersistentFlags().Int("threads", 5, "threads for scan action")
	scanIaCCmd.PersistentFlags().StringP("output", "o", "report.json", "output filepath of report")
	scanIaCCmd.PersistentFlags().StringP("glob", "g", "", "specifies the pattern of plugin file to find")

	scanGitRepoCmd.Flags().String("proxy", "", "proxy to git like: https://xxxxx or socks5://xxxx")
	scanGitRepoCmd.Flags().String("ssh-pubkey", "", "auth to git if use by ssh proto")
	scanGitRepoCmd.Flags().Bool("insecure-skip", false, "skip tls config")
}
