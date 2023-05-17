package parser

import (
	"io"
	"os"
	"strings"

	api "github.com/chaitin/libveinmind/go/iac"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"gopkg.in/yaml.v3"
)

type parseHandle func(file *os.File, path string) (interface{}, error)

func NewParser(iacFile api.IAC) parseHandle {
	switch iacFile.Type {
	case api.Dockerfile:
		return dockerfile
	case api.Kubernetes:
		return kubernetes
	}
	return nil
}

// DockerFileInput is the struct for each dockerfile command
// The input must be the most atomic part of the file
type DockerFileInput struct {
	Cmd       string
	SubCmd    string
	Flags     []string
	Value     []string
	Original  string
	JSON      bool
	Stage     int
	Path      string
	StartLine int
	EndLine   int
}

func dockerfile(file *os.File, path string) (interface{}, error) {

	docker, err := parser.Parse(file)
	if err != nil {
		return nil, err
	}

	var ret []DockerFileInput

	var stageIndex int
	var fromValue = "args"

	for _, child := range docker.AST.Children {

		child.Value = strings.ToLower(child.Value)
		instr, err := instructions.ParseInstruction(child)

		if err != nil {
			return nil, err
		}

		if _, ok := instr.(*instructions.Stage); ok {
			if fromValue != "args" {
				stageIndex++
			}
			fromValue = strings.TrimSpace(strings.TrimPrefix(child.Original, "FROM "))
		}

		cmd := DockerFileInput{
			Cmd:       child.Value,
			Original:  child.Original,
			Flags:     child.Flags,
			Stage:     stageIndex,
			Path:      path,
			StartLine: child.StartLine,
			EndLine:   child.EndLine,
		}

		// Only happens for ONBUILD
		if child.Next != nil && len(child.Next.Children) > 0 {
			cmd.SubCmd = child.Next.Children[0].Value
			child = child.Next.Children[0]
		}

		cmd.JSON = child.Attributes["json"]
		for n := child.Next; n != nil; n = n.Next {
			cmd.Value = append(cmd.Value, n.Value)
		}
		ret = append(ret, cmd)
	}

	return ret, nil
}

type KubernetesInput struct {
	ApiVersion      string            `yaml:"apiVersion" json:"apiVersion"`
	Path            string            `yaml:"path" json:"Path"`
	Kind            string            `yaml:"kind" json:"kind"`
	Meta            interface{}       `yaml:"metadata" json:"metadata"`
	Spec            interface{}       `yaml:"spec" json:"spec"`
	RoleRef         interface{}       `yaml:"roleRef" json:"roleRef"`
	Status          interface{}       `yaml:"status" json:"status"`
	Subjects        interface{}       `yaml:"subjects" json:"subjects"`
	Authentication  interface{}       `yaml:"authentication" json:"authentication"`
	Authorization   interface{}       `yaml:"authorization" json:"authorization"`
	Template        interface{}       `yaml:"template" json:"template"`
	Containers      interface{}       `yaml:"containers" json:"containers"`
	Args            interface{}       `yaml:"args" json:"args"`
	Command         interface{}       `yaml:"command" json:"command"`
	SecurityContext interface{}       `yaml:"securityContext" json:"securityContext"`
	Privileged      interface{}       `yaml:"privileged" json:"privileged"`
	Capabilities    interface{}       `yaml:"capabilities" json:"capabilities"`
	Add             interface{}       `yaml:"add" json:"add"`
	Volumes         interface{}       `yaml:"volumes" json:"volumes"`
	HostPath        interface{}       `yaml:"hostPath" json:"hostPath"`
	HostPID         bool              `yaml:"hostPID" json:"hostPID"`
	Data            map[string]string `yaml:"data" json:"data"`
}

func kubernetes(file *os.File, path string) (interface{}, error) {

	res := make([]*KubernetesInput, 0)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	kubernetesInput := &KubernetesInput{}
	err = yaml.Unmarshal(data, &kubernetesInput)
	if err != nil {
		return nil, err
	}
	kubernetesInput.Path = path

	// check is config file
	if kubernetesInput.Data != nil {
		for _, value := range kubernetesInput.Data {
			// check: is config?
			if strings.HasPrefix(value, "apiVersion") {
				kubernetesTempInput := &KubernetesInput{}
				err = yaml.Unmarshal(data, &kubernetesTempInput)
				if err != nil {
					log.Warnf("parse kubernetes config err: %s", err)
					continue
				}
				kubernetesTempInput.Path = path
				res = append(res, kubernetesTempInput)
			}
		}
	} else {
		res = append(res, kubernetesInput)
	}

	return res, nil
}
