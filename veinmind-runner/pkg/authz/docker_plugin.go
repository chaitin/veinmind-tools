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

type dockerPluginServer struct {
	option *serverOption
}

func NewDockerPlugin(opts ...ServerOption) Server {
	s := &dockerPluginServer{
		option: new(serverOption),
	}
	return newDefaultServer(s, opts...)
}

func (s *dockerPluginServer) init() error {
	return nil
}

func (s *dockerPluginServer) start() error {
	engine := s.registerRouter()
	gin.DefaultWriter = io.MultiWriter(s.option.authLog, os.Stdout)
	go func() {
		err := engine.RunListener(s.option.listener)
		if err != nil {
			log.Error(err)
		}
	}()
	return nil
}

func (s *dockerPluginServer) wait() error {
	return nil
}

func (s *dockerPluginServer) close() {
	err := s.option.authLog.Close()
	if err != nil {
		log.Error(err)
	}

	err = s.option.pluginLog.Close()
	if err != nil {
		log.Error(err)
	}
}

func (s *dockerPluginServer) registerRouter() *gin.Engine {
	engine := gin.New()

	engine.POST("/Plugin.Activate", func(c *gin.Context) {
		c.JSON(200, plugins.Manifest{
			Implements: []string{authorization.AuthZApiImplements},
		})
	})

	engine.POST(fmt.Sprintf("/%s", authorization.AuthZApiRequest), func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error(err)
			c.JSON(200, s.allowAuthResp())
			return
		}

		req := &authorization.Request{}
		err = json.Unmarshal(body, req)
		if err != nil {
			log.Error(err)
			c.JSON(200, s.allowAuthResp())
			return
		}

		c.JSON(200, s.handleAuthZReq(req))
	})

	engine.POST(fmt.Sprintf("/%s", authorization.AuthZApiResponse), func(c *gin.Context) {
		c.JSON(200, s.allowAuthResp())
	})

	return engine
}

func (s *dockerPluginServer) handleAuthZReq(req *authorization.Request) *authorization.Response {
	dockerPluginAction := route.ParseDockerPluginAction(req)
	val, ok := s.option.policies.Load(string(dockerPluginAction))
	if !ok {
		return s.allowAuthResp()
	}
	policy := val.(Policy)

	reportService := report.NewReportService()
	runnerReporter, _ := reporter.NewReporter()
	ctx, cancel := context.WithCancel(context.Background())
	go runnerReporter.Listen()
	go startReportService(ctx, runnerReporter, reportService)

	reportClose := func() {
		cancel()
		runnerReporter.StopListen()
	}

	var (
		result      bool
		err         error
		eventListCh <-chan []reporter.ReportEvent
	)
	switch dockerPluginAction {
	case action.ContainerCreate:
		eventListCh, result, err = handleContainerCreate(policy, req, runnerReporter, reportService)
	case action.ImageCreate:
		eventListCh, result, err = handleImageCreate(policy, req, runnerReporter, reportService)
	case action.ImagePush:
		eventListCh, result, err = handleImagePush(policy, req, runnerReporter, reportService)
	}

	go func() {
		defer reportClose()
		handleReportEvents(eventListCh, policy, s.option.pluginLog, runnerReporter)
	}()

	if err != nil {
		log.Error(err)
		return s.allowAuthResp()
	}

	return s.retAuthResp(result, "", "")
}

func (s *dockerPluginServer) allowAuthResp() *authorization.Response {
	return s.retAuthResp(true, "", "")
}

func (s *dockerPluginServer) forbidAuthResp() *authorization.Response {
	return s.retAuthResp(false, "", "")
}

func (s *dockerPluginServer) retAuthResp(allow bool, msg, err string) *authorization.Response {
	return &authorization.Response{
		Allow: allow,
		Msg:   msg,
		Err:   err,
	}
}
