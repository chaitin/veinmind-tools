package pkg

import (
	"bufio"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-escape/rules"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/pelletier/go-toml"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
)

var cveList = make([]CVE, 0)

var taskList = make([]tasks, 0)

type versionCheck struct {
	BeginVersion         string
	EndVersion           string
	BeginEqual, EndEqual bool
}

type CVE struct {
	CVENumber        string
	versionchecklist []versionCheck
}

type tasks struct {
	CVEinfo CVE
	para    interface{}
}

// 需要自定义函数进行检测时可以使用该函数获取*CVE对象，再使用setFunc设置自定义检测函数
func getCVEObject(name string) *CVE {
	for _, value := range cveList {
		if value.CVENumber == name {
			return &value
		}
	}
	return &CVE{}
}

func (task *tasks) check() (*event.EscapeDetail, error) {
	versions := task.para.([]string)
	KernelVersion, err := strconv.Atoi(versions[1])
	if err != nil {
		log.Error(err)
		return nil, err
	}
	MajorRevision, err := strconv.Atoi(versions[2])
	if err != nil {
		log.Error(err)
		return nil, err
	}
	MinorRevision, err := strconv.Atoi(versions[3])
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var versionString string
	for _, value := range versions {
		versionString += value
	}

	for _, ver := range task.CVEinfo.versionchecklist {
		cveBeginVersions := strings.Split(ver.BeginVersion, ".")
		cveEndVersions := strings.Split(ver.EndVersion, ".")
		if len(cveBeginVersions) == 0 || len(cveEndVersions) == 0 {
			continue
		}
		//若ver为5.1格式，则将其规范为5.1.0,并将其中的string转为int
		cveBeginVersionsInt := make([]int, 0)
		cveEndVersionsInt := make([]int, 0)
		for i := 0; i < 3; i++ {
			if i < len(cveBeginVersions) {
				intVersion, err := strconv.Atoi(cveBeginVersions[i])
				if err != nil {
					continue
				} else {
					cveBeginVersionsInt = append(cveBeginVersionsInt, intVersion)
				}
			} else {
				cveBeginVersionsInt = append(cveBeginVersionsInt, 0)
			}

			if i < len(cveEndVersions) {
				intVersion, err := strconv.Atoi(cveEndVersions[i])
				if err != nil {
					continue
				} else {
					cveEndVersionsInt = append(cveEndVersionsInt, intVersion)
				}
			} else {
				cveEndVersionsInt = append(cveEndVersionsInt, 0)
			}
		}
		if morethan(cveBeginVersionsInt, []int{KernelVersion, MajorRevision, MinorRevision}, ver.BeginEqual) &&
			lessthan(cveEndVersionsInt, []int{KernelVersion, MajorRevision, MinorRevision}, ver.EndEqual) {
			return &event.EscapeDetail{
				Target: "KERNEL VERSION",
				Reason: CVEREASON + task.CVEinfo.CVENumber,
				Detail: "UnSafeKernelVersion " + versionString,
			}, nil
		}
	}
	return nil, nil
}

func getCVEFromFile() error {
	f, err := rules.Open("rule.toml")
	if err != nil {
		log.Error(err)
		return err
	}
	config, err := toml.LoadReader(f)
	if err != nil {
		log.Error(err)
		return err
	}
	cf := config.Get("cve").([]*toml.Tree)
	for _, value := range cf {
		tmpCVE := CVE{}
		tmpCVE.CVENumber = value.Get("cveNumber").(string)
		tmpCVE.versionchecklist = []versionCheck{}
		for _, version := range value.Get("version").([]interface{}) {
			tmpCVE.versionchecklist = append(tmpCVE.versionchecklist, parserVersion(version.(string)))
		}
		cveList = append(cveList, tmpCVE)
	}
	return nil
}

func parserVersion(version string) versionCheck {
	res := strings.Split(version, "ver")
	versionCheck1 := versionCheck{}
	if res[0] != "" && res[1] != "" { // eg: 1.1.1<ver<=2.2.2
		var opera1, opera2 string
		for i := 0; i < len(res[0]); i++ { // res[0] : 1.1.1<
			if string(res[0][i]) == "=" || string(res[0][i]) == "<" || string(res[0][i]) == ">" {
				opera1 = res[0][i:] // opera1 : <
				res[0] = res[0][:i] // res[0] : 1.1.1
				break
			}
		}
		for i := 0; i < len(res[1]); i++ { // res[1] : <=2.2.2
			if unicode.IsNumber(rune(res[1][i])) {
				opera2 = res[1][:i] // opera2 : <=
				res[1] = res[1][i:] // res[1] : 2.2.2
				break
			}
		}
		opera1 = strings.TrimSpace(opera1)
		opera2 = strings.TrimSpace(opera2)
		res[0] = strings.TrimSpace(res[0])
		res[1] = strings.TrimSpace(res[1])
		switch opera1 {
		case "<=":
			versionCheck1.BeginVersion = res[0]
			versionCheck1.BeginEqual = true
		case "<":
			versionCheck1.BeginVersion = res[0]
			versionCheck1.BeginEqual = false
		case ">":
			versionCheck1.EndVersion = res[0]
			versionCheck1.EndEqual = false
		case ">=":
			versionCheck1.EndVersion = res[0]
			versionCheck1.EndEqual = true
		}
		switch opera2 {
		case ">=":
			versionCheck1.BeginVersion = res[1]
			versionCheck1.BeginEqual = true
		case ">":
			versionCheck1.BeginVersion = res[1]
			versionCheck1.BeginEqual = false
		case "<":
			versionCheck1.EndVersion = res[1]
			versionCheck1.EndEqual = false
		case "<=":
			versionCheck1.EndVersion = res[1]
			versionCheck1.EndEqual = true
		}
	} else { // eg: ver>=3.3.3
		var result string
		if res[0] != "" {
			result = res[0]
		} else if res[1] != "" {
			result = res[1]
		} else {
			return versionCheck1
		}
		var opera string
		result = strings.TrimSpace(result)
		for i := 0; i < len(result); i++ { //res[0] : >=3.3.3
			if unicode.IsNumber(rune(result[i])) {
				opera = result[:i]  // opera >=
				result = result[i:] // res[0] 3.3.3
				break
			}
		}
		opera = strings.TrimSpace(opera)
		result = strings.TrimSpace(result)
		switch opera {
		case "<":
			versionCheck1.EndVersion = result
			versionCheck1.EndEqual = false
			versionCheck1.BeginVersion = "-1.-1.-1"
			versionCheck1.BeginEqual = false
		case "<=":
			versionCheck1.EndVersion = result
			versionCheck1.EndEqual = true
			versionCheck1.BeginVersion = "-1.-1.-1"
			versionCheck1.BeginEqual = false
		case ">":
			versionCheck1.BeginVersion = result
			versionCheck1.BeginEqual = false
			versionCheck1.EndVersion = "-1.-1.-1"
			versionCheck1.EndEqual = false
		case ">=":
			versionCheck1.BeginVersion = result
			versionCheck1.BeginEqual = true
			versionCheck1.EndVersion = "-1.-1.-1"
			versionCheck1.EndEqual = false

		}
	}
	return versionCheck1
}

func getVersion() ([]string, error) {
	content, err := os.Open("/proc/version")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer content.Close()
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		complie := regexp.MustCompile(KERNELPATTERN)
		res1 := complie.FindStringSubmatch(scanner.Text())
		if len(res1) == 4 {
			return res1, nil
		}
	}
	return nil, nil
}

func morethan(cveVersion []int, inputVersion []int, equal bool) bool {
	if cveVersion[0] == -1 && cveVersion[1] == -1 && cveVersion[2] == -1 {
		return true
	}
	if equal && cveVersion[0] == inputVersion[0] && cveVersion[1] == inputVersion[1] && cveVersion[2] == inputVersion[2] {
		return true
	} else {
		if inputVersion[0] > cveVersion[0] {
			return true
		} else if inputVersion[0] == cveVersion[0] {
			if inputVersion[1] > cveVersion[1] {
				return true
			} else if inputVersion[1] == cveVersion[1] {
				if inputVersion[2] > cveVersion[2] {
					return true
				}
			}
		}
	}
	return false
}

func lessthan(cveVersion []int, inputVersion []int, equal bool) bool {
	if cveVersion[0] == -1 && cveVersion[1] == -1 && cveVersion[2] == -1 {
		return true
	}
	if equal && cveVersion[0] == inputVersion[0] && cveVersion[1] == inputVersion[1] && cveVersion[2] == inputVersion[2] {
		return true
	} else {
		if inputVersion[0] < cveVersion[0] {
			return true
		} else if inputVersion[0] == cveVersion[0] {
			if inputVersion[1] < cveVersion[1] {
				return true
			} else if inputVersion[1] == cveVersion[1] {
				if inputVersion[2] < cveVersion[2] {
					return true
				}
			}
		}
	}
	return false
}

// ContainerCVECheck 此处传入fs api.FileSystem只是为了和其他检测函数统一格式，实际并无作用
func ContainerCVECheck(fs api.FileSystem) ([]*event.EscapeDetail, error) {
	var res = make([]*event.EscapeDetail, 0)
	for _, task := range taskList {
		check, err := task.check()
		if err == nil && check != nil {
			res = append(res, check)
		}
	}
	return res, nil
}

func init() {
	err := getCVEFromFile()
	if err != nil {
		return
	}
	for _, value := range cveList {
		para, err := getVersion()
		if err != nil {
			log.Error(err)
			return
		}
		task := tasks{
			CVEinfo: value,
			para:    para,
		}
		taskList = append(taskList, task)
	}

	ContainerCheckList = append(ContainerCheckList, ContainerCVECheck)
}
