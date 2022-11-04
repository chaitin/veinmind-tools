package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-iac/pkg/scanner"
)

const (
	JSON   = "json"
	STDOUT = "stdout"
)

var tab = 30 // stdout输出视图大小控制
var tabwTitle = tabwriter.NewWriter(os.Stdout, 3+(tab+2)*4, 3+(tab+2)*4, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
var tabwBody = tabwriter.NewWriter(os.Stdout, tab+2, (tab+2)*4, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)

type StdStream func() error

func Stream(spend time.Duration, scanTotal int, stdstream StdStream) error {
	fmt.Fprintln(tabwTitle, "#", Repeat("=", (tab+2)*4), "#")
	fmt.Fprintln(tabwTitle, "| Scan File Total: ", scanTotal, "\t")
	fmt.Fprintln(tabwTitle, "| Spend Time: ", spend.String(), "\t")

	err := stdstream()
	if err != nil {
		return err
	}

	fmt.Fprintln(tabwBody, "#", Repeat("=", (tab+2)*4), "#")
	tabwBody.Flush()
	return nil
}

func Stdout(results []scanner.Result) error {
	for _, r := range results {
		fmt.Fprintln(tabwTitle, "+", Repeat("-", (tab+2)*4), "+")
		fmt.Fprintln(tabwBody, Align("| RuleID", r.Id, tab), "\t", Align("Rule Name", r.Name, tab), "\t", Align("Rule Level", r.Severity, tab), "\t", Align("Rule Type", r.Type, tab), "\t")
		fmt.Fprintln(tabwBody, "+", Repeat("-", (tab+2)*4), "#")
		fmt.Fprintln(tabwBody, "| Start Line", "\t", "End Line", "\t", "File", "\t", "Code", "\t")
		for _, risk := range r.Risks {
			fmt.Fprintln(tabwBody, "|", risk.StartLine, "\t", risk.StartLine, "\t", Limit(risk.FilePath, tab), "\t", Limit(risk.Original, tab), "\t")
		}
		fmt.Fprintln(tabwBody, "+", Repeat("-", (tab+2)*4), "#")
	}
	return nil
}

func Json(res []scanner.Result) error {
	var jsonFile *os.File
	var name = "report.json"
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		jsonFile, err = os.Create(name)
		if err != nil {
			return err
		}
	} else {
		err := os.Remove(name)
		if err != nil {
			return err
		}

		jsonFile, err = os.Create(name)
		if err != nil {
			return err
		}
	}

	scanResultJson, err := json.Marshal(res)
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(scanResultJson)
	if err != nil {
		return err
	}

	tabw := tabwriter.NewWriter(os.Stdout, 3+(tab+2)*4, 3+(tab+2)*4, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
	out := "| The Detail Results Export report.json"
	fmt.Fprintln(tabw, out, Repeat(" ", (tab+2)*4-len(out)+1), "|")
	return nil
}

func Limit(s string, num int) string {
	if len(s) > num {
		if num > 3 {
			return s[0:num-3] + "..."
		} else {
			return "..."
		}
	} else {
		return s
	}
}

func Repeat(s string, num int) string {
	res := ""
	for i := 0; i < num; i++ {
		res += s
	}
	return res
}

func Align(key string, value string, tab int) string {
	if len(key)+len(value)+2 < tab {
		return key + ":" + Repeat(" ", tab-len(key)-len(value)-2) + value
	} else {
		return key + ": " + Limit(value, tab-len(key)-2)
	}
}
