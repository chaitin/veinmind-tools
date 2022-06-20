package authz

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/registry"
	"github.com/docker/docker/pkg/authorization"
)

var pullAction actionMap
var regexIDSPattern = `^[0-9a-f]+$`
var regexImageNamePattern = `images/(.*?)/push`

type Authorizer interface {
	Init(config *PolicyConfig) error                             // Init initialize the handler
	AuthZReq(req *authorization.Request) *authorization.Response // AuthZReq handles the request from docker client
	AuthZRes(req *authorization.Request) *authorization.Response // AuthZRes handles the response from docker daemon to docker clien
}

type BasicAuthorizer struct {
	config       *PolicyConfig
	dockerClient registry.Client
	policyMap    map[string]Policy
}

func (b *BasicAuthorizer) Init(cfg *PolicyConfig) error {
	b.config = cfg
	b.policyMap = cfg.PolicysMap()
	dockerClient, err :=
		registry.NewRegistryDockerClient(registry.WithAuthField(cfg.DockerAuth))
	if err != nil {
		return err
	}
	b.dockerClient = dockerClient
	return nil
}

func (b *BasicAuthorizer) AuthZReq(req *authorization.Request) *authorization.Response {
	// get action from URL
	action := ParseRoute(req.RequestMethod, req.RequestURI)
	// make sure Report file was created
	// and open the Report file
	reportFile, err := GetLogFile(b.config.Log.ReportLogPath)
	if err != nil {
		return defaultAuthResponse()
	}
	// make sure user set the policy for the action
	policy, ok := b.policyMap[action]
	// if action in policy,then check
	// else pass the request
	if ok {
		switch action {
		case ActionContainerCreate:
			{
				imageName, err := getImageNameFromJson(req, "Image")
				if err != nil {
					log.Error(err)
					return defaultAuthResponse()
				}
				block, err := CheckImage(imageName, policy, reportFile)
				if err != nil {
					log.Error(err)
					return defaultAuthResponse()
				}
				if block {
					return defaultRejectResponse()
				} else {
					return defaultAuthResponse()
				}
			}
		case ActionImageCreate:
			{
				imageName, err := getImageNameFromTextPlain(req, "fromImage")
				if err != nil {
					log.Error(err)
					return defaultAuthResponse()
				}
				//make a ratelimit
				//if rate more than 10, return reject resp
				if pullAction.Count(imageName) > 100 {
					return defaultRejectResponse()
				}
				imageActionId := fmt.Sprintf("%s-%d", imageName, time.Now().UnixMicro())
				pullAction.Store(imageActionId, struct{}{})
				//start a goroutine to handle the policy

				go func() {
					defer pullAction.Delete(imageActionId)

					ticker := time.NewTicker(time.Second)
					veinmindRuntime, _ := docker.New()
					for {
						select {
						case <-time.After(time.Minute * 30):
							//timeout limit
							return
						case <-ticker.C:
							//check the image existed
							//if with block policy, delete the image and exit goroutine
							imageIds, err := veinmindRuntime.FindImageIDs(imageName)
							if err != nil {
								log.Error(err)
								break
							}
							if len(imageIds) > 0 {
								//scan image
								block, err := CheckImage(imageName, policy, reportFile)
								if err != nil {
									log.Error(err)
								}
								//policy handle
								if block {
									log.Infof("Risk found in %s, deleting...", imageName)
									err := b.dockerClient.Remove(imageName)
									if err != nil {
										log.Errorf("Remove image failed: %#v\n", imageName)
									}
									log.Infof("Remove image success: %#v\n", imageName)
								} else {
									log.Info("Pull image success!")
								}
								return
							}
						}
					}
				}()
				return defaultAuthResponse()
			}
		case ActionImagePush:
			{
				imageName, err := getImageNameFromUrlPath(req)
				if err != nil {
					log.Error(err)
					return defaultAuthResponse()
				}
				block, err := CheckImage(imageName, policy, reportFile)
				if err != nil {
					return defaultAuthResponse()
				}
				if block {
					return defaultRejectResponse()
				}
				return defaultAuthResponse()

			}
		default:
			{
				return defaultAuthResponse()
			}
		}

	}

	return defaultAuthResponse()

}
func getImageNameFromUrlPath(req *authorization.Request) (string, error) {
	uri := req.RequestURI
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	path := u.Path
	// post url like : /v1.41/images/{name}/push
	// ref:https://docs.docker.com/engine/api/v1.41/#tag/Image/operation/ImageHistory
	reg := regexp.MustCompile(regexImageNamePattern)
	imageNames := reg.FindStringSubmatch(path)
	if len(imageNames) != 2 {
		return "", errors.New("parse url path error")
	}
	kv, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	// tag as a query paramters
	tags, ok := kv["tag"]
	if !ok {
		return "", errors.New("bad formate request")
	}
	imageName := imageNames[1] + ":" + tags[0]
	return imageName, nil
}

func getImageNameFromTextPlain(req *authorization.Request, key string) (string, error) {
	uri := req.RequestURI
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	// post url like : /v1.41/images/create?fromImage=huzai9527%2Fbabyrop&tag=latest
	// ref:https://docs.docker.com/engine/api/v1.41/#tag/Image/operation/ImageCreate
	kv, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	imageNames, ok := kv[key]
	if !ok {
		return "", errors.New("bad formate request")
	}
	tags, ok := kv["tag"]
	if !ok {
		return "", errors.New("bad formate request")
	}
	imageName := imageNames[0] + ":" + tags[0]
	return imageName, nil

}
func getImageNameFromJson(req *authorization.Request, key string) (string, error) {
	// parse request body from URL and get image's name
	var body map[string]interface{}
	if req.RequestHeaders["Content-Type"] == "application/json" && len(req.RequestBody) > 0 {
		if err := json.Unmarshal(req.RequestBody, &body); err != nil {
			return "", err
		}
	}
	fmt.Println(req.RequestURI)
	// imageName should be taged
	// if doesn't has Image Key then return defaultResponse
	imageNameI, ok := body[key]
	if !ok {
		return "", errors.New("Bad request head format")
	}

	imageName, ok := imageNameI.(string)
	if !ok {
		return "", errors.New("Fail to parse image name with string")
	}
	// cause only get imagename of string type
	// it may be patial digest id, such as "92e0f4bd4b90"
	// besides, it may only a string like "tocmat"
	// in this way, docker will take "tocmat" as "tocmat:latest"
	// it also may "tocmat:1.9" with tag or "tomcat:1.9@shaxxx" with tag and id
	// specially, if the imagename is "abcd" then we just return this to vemind-sdk
	// vemind-sdk will find all relative images whether "abcd" is name or id
	// reference : https://github.com/docker/docker-ce/blob/5d94ad617b913e7eaa5adb65dd6260d0aa87f9c9/components/engine/daemon/images/image.go#L150
	var HexRegexpAnchored = regexp.MustCompile(regexIDSPattern)
	if !HexRegexpAnchored.MatchString(imageName) {
		if strings.Count("imageName", ":") == 1 {
			imageName = imageName + ":latest"
		}
	}
	return imageName, nil

}

func (b *BasicAuthorizer) AuthZRes(req *authorization.Request) *authorization.Response {
	return defaultAuthResponse()
}

func defaultAuthResponse() *authorization.Response {
	return &authorization.Response{
		Allow: true,
	}
}

func defaultRejectResponse() *authorization.Response {
	return &authorization.Response{
		Allow: false,
		Msg:   fmt.Sprintln("This image has risk!"),
	}
}
