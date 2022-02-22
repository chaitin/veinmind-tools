package common_cli

import (
	"fmt"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/embed"
	common "github.com/chaitin/veinmind-tools/veinmind-weakpass/log"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/scan"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"text/tabwriter"
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

				tabw := tabwriter.NewWriter(os.Stdout, 95, 95, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
				fmt.Fprintln(tabw, "# ============================================================================================ #")
				spend := time.Since(scanStart)
				fmt.Fprintln(tabw, "| Scan Total: ", strconv.Itoa(len(results)), "\t")
				fmt.Fprintln(tabw, "| Spend Time: ", spend.String(), "\t")
				var weakpassImageTotal = 0
				var weakpassTotal = 0
				for _, r := range results {
					if len(r.WeakpassResults) > 0 {
						weakpassImageTotal++
						weakpassTotal += len(r.WeakpassResults)
					}
				}
				fmt.Fprintln(tabw, "| Weakpass Image Total: ", strconv.Itoa(weakpassTotal), "\t")
				fmt.Fprintln(tabw, "| Weakpass Total: ", strconv.Itoa(weakpassTotal), "\t")
				fmt.Fprintln(tabw, "+----------------------------------------------------------------------------------------------+")

				for _, r := range results {
					if len(r.WeakpassResults) > 0 {
						fmt.Fprintln(tabw, "| ImageName: ", r.ImageName, "\t")
						fmt.Fprintln(tabw, "| Status: Unsafe", "\t")
						for _, w := range r.WeakpassResults {
							fmt.Fprintln(tabw, "| Username: ", w.Username, "\t")
							fmt.Fprintln(tabw, "| Password: ", w.Password, "\t")
						}
						fmt.Fprintln(tabw, "+----------------------------------------------------------------------------------------------+")
					}
				}
				fmt.Fprintln(tabw, "# ============================================================================================ #\n")
				tabw.Flush()

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
