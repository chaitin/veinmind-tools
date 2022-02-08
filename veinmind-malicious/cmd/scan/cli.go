//+build community

package main

import (
	"github.com/chaitin/veinmind-tools/veinmind-malicious/cmd/scan/common_cli"
	_ "github.com/chaitin/veinmind-tools/veinmind-malicious/config"
	_ "github.com/chaitin/veinmind-tools/veinmind-malicious/database"
	_ "github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/common"
	_ "net/http/pprof"
	"os"
)

func main() {
	err := common_cli.App.Run(os.Args)
	if err != nil {
		common.Log.Fatal(err)
	}
}
