package authz

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/action"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/plugins"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	multiWriter := io.MultiWriter(s.option.authLog, os.Stdout)

	logger := logrus.New()
	logger.Out = multiWriter

	log.SetDefaultLogger(log.NewLogrus(logger))
	gin.DefaultWriter = multiWriter

	engine := s.registerRouter()

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
	engine := gin.Default()

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
	var (
		eventListCh <-chan []reporter.ReportEvent
		result      bool
		err         error
	)
	switch dockerPluginAction {
	case action.ContainerCreate:
		eventListCh, result, err = handleContainerCreate(policy, req)
	case action.ImageCreate:
		eventListCh, result, err = handleImageCreate(policy, req)
	case action.ImagePush:
		eventListCh, result, err = handleImagePush(policy, req)
	default:
		eventListCh, result, err = handleDefaultAction()
	}
	go func() {
		handleReportEvents(eventListCh, policy, s.option.pluginLog)
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
