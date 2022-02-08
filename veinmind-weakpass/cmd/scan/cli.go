//+build community

package main

import (
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/cmd/scan/common_cli"
	common "github.com/chaitin/veinmind-tools/veinmind-weakpass/log"
	_ "net/http/pprof"
	"os"
)

func main() {
	err := common_cli.App.Run(os.Args)
	if err != nil {
		common.Log.Fatal(err)
	}
}
