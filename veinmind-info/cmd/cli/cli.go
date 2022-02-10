package main

import (
	"encoding/json"
	common "github.com/chaitin/veinmind-tools/veinmind-info/log"
	"github.com/chaitin/veinmind-tools/veinmind-info/scan"
	"github.com/urfave/cli/v2"
	"os"
)

var App = &cli.App{
	Name:  "veinmind-info",
	Usage: "veinmind-info is a image info scanner",
	Commands: []*cli.Command{
		{
			Name:  "scan",
			Usage: "scan image info, image e.g. ubuntu:latest",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "engine",
					Value:   "dockerd",
					Usage:   "scan engine e.g. dockerd",
					Aliases: []string{"e"},
				},
			},
			Action: func(c *cli.Context) error {
				results, err := scan.Scan(scan.ScanOption{
					EngineType: func() scan.EngineType {
						if v, ok := scan.EngineTypeMap[c.String("engine")]; ok {
							return v
						} else {
							common.Log.Fatal("Engine type doesn't match")
							return -1
						}
					}(),
					ImageName: c.Args().First(),
				})

				if err != nil {
					return err
				}

				resultJson, _ := json.Marshal(results)
				common.Log.Info(string(resultJson))

				return nil
			},
		},
	},
}

func main() {
	err := App.Run(os.Args)
	if err != nil {
		common.Log.Fatal(err)
	}
}
