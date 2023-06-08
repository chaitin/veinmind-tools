package utils

import (
	"reflect"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-privilege-escalation/rules"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-privilege-escalation/service"
)

var binPath = []string{"/bin", "/sbin", "/usr/bin", "/usr/sbin", "/usr/local/bin", "/usr/local/sbin"}

func ImagesScanRun(fs api.Image) []*event.EscalationDetail {
	var result = make([]*event.EscalationDetail, 0)
	// 加载规则
	Config, err := rules.GetRuleFromFile()
	if err != nil {
		return nil
	}
	// 获取文件系统中可能存在提权等风险的二进制文件
	for _, rule := range Config.Rules {
		for _, path := range binPath {
			content, err := fs.Stat(path + "/" + rule.Name)
			if err != nil {
				continue
			}

			// 根据文件对应的Tags,调用对应检查模块，判断是否有提权风险
			for _, tag := range rule.Tags {
				if checkFunc, ok := service.ImageCheckFuncMap[tag]; ok {
					// 判断是否存在有风险的命令
					risk, err := checkFunc(fs, content, path+"/"+rule.Name)
					if err != nil {
						continue
					}

					// 如果有风险，将风险信息添加到result中
					if risk {
						tag = convertToCamelCase(tag)
						exps := reflect.ValueOf(rule.Exps).FieldByName(tag).Interface().([]*rules.Exp)
						for _, exp := range exps {
							result = append(result, &event.EscalationDetail{
								BinName:     rule.Name,
								Description: rule.Description,
								FilePath:    path + "/" + rule.Name,
								Mod:         tag,
								Exp:         exp.Exp,
							})
						}
					}
				}
			}
		}
	}
	return result
}

func ContainersScanRun(fs api.Container) []*event.EscalationDetail {
	var result = make([]*event.EscalationDetail, 0)
	// 加载规则
	Config, err := rules.GetRuleFromFile()
	if err != nil {
		return nil
	}
	// 获取文件系统中可能存在提权等风险的二进制文件
	for _, rule := range Config.Rules {
		for _, path := range binPath {
			content, err := fs.Stat(path + "/" + rule.Name)
			if err != nil {
				continue
			}

			// 根据文件对应的Tags,调用对应检查模块，判断是否有提权风险
			for _, tag := range rule.Tags {
				if checkFunc, ok := service.ContainerCheckFuncMap[tag]; ok {
					// 判断是否存在有风险的命令
					risk, err := checkFunc(fs, content, path+"/"+rule.Name)
					if err != nil {
						continue
					}

					// 如果有风险，将风险信息添加到result中
					if risk {
						tag = convertToCamelCase(tag)
						exps := reflect.ValueOf(rule.Exps).FieldByName(tag).Interface().([]*rules.Exp)
						for _, exp := range exps {
							result = append(result, &event.EscalationDetail{
								BinName:     rule.Name,
								Description: rule.Description,
								FilePath:    path + "/" + rule.Name,
								Mod:         tag,
								Exp:         exp.Exp,
							})
						}
					}
				}
			}
		}
	}
	return result
}

// convertToCamelCase 将短横线命名转换为驼峰命名
func convertToCamelCase(s string) string {
	words := strings.Split(strings.ToLower(s), "-")

	for i, word := range words {
		if word == "suid" {
			words[i] = "SUID"
		} else {
			words[i] = strings.Title(word)
		}
	}

	return strings.Join(words, "")
}
