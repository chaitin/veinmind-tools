package test

import (
	"context"
	"strconv"
	"testing"

	"github.com/chaitin/libveinmind/go/iac"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-iac/pkg/scanner"
	"github.com/open-policy-agent/opa/ast"
)

const BasicPath = "./docker"
const DockerfileTestNum = 30

func Test_DockerfileRules(t *testing.T) {
	scanner := &scanner.Scanner{
		QueryPre: "data.brightMirror.",
		Policies: make(map[string]*ast.Module),
	}
	err := scanner.LoadLibs()
	if err != nil {
		t.Errorf("scanner load libs err:%v", err)
	}

	index := 1
	for index <= DockerfileTestNum {
		var n string
		switch {
		case index < 10:
			n = "00" + strconv.Itoa(index)
		case index < 100:
			n = "0" + strconv.Itoa(index)
		default:
			n = strconv.Itoa(index)
		}

		res, _ := scanner.Scan(context.Background(), iac.IAC{
			Type: iac.Dockerfile,
			Path: BasicPath + "/DF-" + n + "/DF-" + n + "-noncompliant.Dockerfile",
		})

		if res[0].Rule.Id == "DF-"+n {
			t.Logf("DF-%s test pass", n)
		}

		index++
	}

}
