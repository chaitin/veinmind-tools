package authz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/action"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/plugins"
	"github.com/gin-gonic/gin"
)

type dockerPluginServer serverOption

func NewDockerPlugin(opts ...ServerOption) Server {
	s := &dockerPluginServer{
		policies: newPolicyMap(),
	}
	return newDefaultServer(s, opts...)
}

func (my *dockerPluginServer) init() error {
	return nil
}

func (my *dockerPluginServer) start() error {
	engine := my.registerRouter()
	gin.DefaultWriter = io.MultiWriter(my.authLog, os.Stdout)
	go func() {
		err := engine.RunListener(my.listener)
		if err != nil {
			log.Error(err)
		}
	}()
	return nil
}

func (my *dockerPluginServer) wait() error {
	return nil
}

func (my *dockerPluginServer) handleAuthZReq(req *authorization.Request) *authorization.Response {
	dockerPluginAction := route.ParseDockerPluginAction(req)
	policy, ok := my.policies.Load(string(dockerPluginAction))
	if !ok {
		return my.allowAuthResp()
	}

	reportService := report.NewReportService()
	runnerReporter, _ := reporter.NewReporter()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runnerReporter.Listen()
	go startReportService(ctx, runnerReporter, reportService)

	var err error
	var result bool
	switch dockerPluginAction {
	case action.ContainerCreate:
		result, err = handleContainerCreate(policy, req, runnerReporter, reportService)
	case action.ImageCreate:
		result, err = handleImageCreate(policy, req, runnerReporter, reportService)
	case action.ImagePush:
		result, err = handleImagePush(policy, req, runnerReporter, reportService)
	}
	if err != nil {
		log.Error(err)
		return my.allowAuthResp()
	}

	handleReportAlert(policy, runnerReporter)
	handleReportLog(policy, my.pluginLog, runnerReporter)
	return my.retAuthResp(result, "", "")
}

func (my *dockerPluginServer) allowAuthResp() *authorization.Response {
	return my.retAuthResp(true, "", "")
}

func (my *dockerPluginServer) forbidAuthResp() *authorization.Response {
	return my.retAuthResp(false, "", "")
}

func (my *dockerPluginServer) retAuthResp(allow bool, msg, err string) *authorization.Response {
	return &authorization.Response{
		Allow: allow,
		Msg:   msg,
		Err:   err,
	}
}

func (my *dockerPluginServer) registerRouter() *gin.Engine {
	engine := gin.New()

	engine.Any("/Plugin.Activate", func(c *gin.Context) {
		c.JSON(200, plugins.Manifest{
			Implements: []string{authorization.AuthZApiImplements},
		})
	})

	engine.Any(fmt.Sprintf("/%s", authorization.AuthZApiRequest), func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error(err)
			c.JSON(200, my.allowAuthResp())
			return
		}

		req := authorization.Request{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Error(err)
			c.JSON(200, my.allowAuthResp())
			return
		}

		c.JSON(200, my.handleAuthZReq(&req))
	})

	engine.Any(fmt.Sprintf("/%s", authorization.AuthZApiResponse), func(c *gin.Context) {
		c.JSON(200, authorization.Response{
			Allow: true,
		})
	})

	return engine
}
