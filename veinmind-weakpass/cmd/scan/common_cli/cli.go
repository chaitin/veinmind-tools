package common_cli

import (
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/embed"
	common "github.com/chaitin/veinmind-tools/veinmind-weakpass/log"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/scan"
	"github.com/urfave/cli/v2"
	"time"
)

var App = &cli.App{
	Name:  "veinmind-weakpass",
	Usage: "veinmind-weakpass is a image weakpass scanner",
	Commands: []*cli.Command{
		{
			Name:  "scan",
			Usage: "scan image weakpass, image e.g. ubuntu:latest",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "threads",
					Value:   10,
					Usage:   "password brute threads",
					Aliases: []string{"t"},
				},
				&cli.StringFlag{
					Name:    "engine",
					Value:   "dockerd",
					Usage:   "scan engine e.g. dockerd",
					Aliases: []string{"e"},
				},
				&cli.StringFlag{
					Name:    "dictpath",
					Value:   "",
					Usage:   "dict path e.g ./mypass.dict",
					Aliases: []string{"d"},
				},
				&cli.StringFlag{
					Name:    "username",
					Value:   "",
					Usage:   "username e.g. root",
					Aliases: []string{"u"},
				},
			},
			Action: func(c *cli.Context) error {
				// 记录扫描开始时间
				scanStart := time.Now()

				p := scan.SSHScanPlugin{}
				results, err := p.Scan(scan.ScanOption{
					EngineType: func() scan.EngineType {
						if v, ok := scan.EngineTypeMap[c.String("engine")]; ok {
							return v
						} else {
							common.Log.Fatal("Engine type doesn't match")
							return -1
						}
					}(),
					ImageName:   c.Args().First(),
					ScanThreads: c.Int("threads"),
					Username:    c.String("username"),
					Dictpath:    c.String("dictpath"),
				})

				if err != nil {
					return err
				}

				for _, r := range results {
					if len(r.WeakpassResults) == 0 {
						common.Log.Info("ImageName: ", r.ImageName)
						common.Log.Info("Status: Safe")
						common.Log.Info("====================================================================================")
					} else {
						common.Log.Warn("ImageName: ", r.ImageName)
						common.Log.Warn("Status: Unsafe")
						for _, w := range r.WeakpassResults {
							common.Log.Warn("Username: ", w.Username)
							common.Log.Warn("Password: ", w.Password)
						}
						common.Log.Info("====================================================================================")
					}
				}
				spend := time.Since(scanStart)
				common.Log.Info("Scan Spend Time: ", spend.String())
				common.Log.Info("====================================================================================")

				return nil
			},
		},
		{
			Name:  "extract",
			Usage: "extract dict file to disk",
			Action: func(c *cli.Context) error {
				embed.ExtractAll()
				return nil
			},
		},
	},
}
