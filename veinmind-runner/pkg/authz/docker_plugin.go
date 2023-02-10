package authz

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"io"
	"net"
	"os"
	"sync"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/plugins"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/action"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
)

type DockerServerOption func(option *dockerServerOption) error

type dockerServerOption struct {
	authLog   io.WriteCloser
	pluginLog io.WriteCloser
	policies  sync.Map
	listener  net.Listener
}

func WithPolicy(policies ...Policy) DockerServerOption {
	return func(option *dockerServerOption) error {
		for _, policy := range policies {
			option.policies.Store(policy.Action, policy)
		}

		return nil
	}
}

func WithAuthLog(path string) DockerServerOption {
	return func(option *dockerServerOption) error {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			_, err = os.Create(path)
			if err != nil {
				return err
			}
		}

		fp, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		option.authLog = fp
		return nil
	}
}

func WithPluginLog(path string) DockerServerOption {
	return func(option *dockerServerOption) error {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			_, err = os.Create(path)
			if err != nil {
				return err
			}
		}

		fp, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		option.pluginLog = fp
		return nil
	}
}

func WithListenerUnix(addr string) DockerServerOption {
	return func(option *dockerServerOption) error {
		listener, err := net.ListenUnix("unix", &net.UnixAddr{Net: "unix", Name: addr})
		if err != nil {
			return err
		}

		option.listener = listener
		return nil
	}
}

func WithServerOptions(options ...DockerServerOption) DockerServerOption {
	return func(s *dockerServerOption) error {
		for _, option := range options {
			if err := option(s); err != nil {
				return err
			}
		}

		return nil
	}
}

type dockerPluginServer struct {
	defaultServer
	dockerOpt *dockerServerOption
}

func NewDockerPluginServer(opts ...DockerServerOption) (dockerPluginServer, error) {
	option := &dockerServerOption{}
	for _, opt := range opts {
		if err := opt(option); err != nil {
			return dockerPluginServer{}, err
		}
	}
	return dockerPluginServer{dockerOpt: option}, nil

}

func (s *dockerPluginServer) Init() error {
	opts := make([]DockerServerOption, 0)
	if s.dockerOpt.authLog == nil {
		opts = append(opts, WithAuthLog(defaultAuthLogPath))
	}
	if s.dockerOpt.pluginLog == nil {
		opts = append(opts, WithPluginLog(defaultPluginPath))
	}
	if s.dockerOpt.listener == nil {
		opts = append(opts, WithListenerUnix(defaultSockListenAddr))
	}
	if err := WithServerOptions(opts...)(s.dockerOpt); err != nil {
		return err
	}

	return nil
}

func (s *dockerPluginServer) Start() error {
	multiWriter := io.MultiWriter(s.dockerOpt.authLog, os.Stdout)

	logger := logrus.New()
	logger.Out = multiWriter

	log.SetDefaultLogger(log.NewLogrus(logger))
	gin.DefaultWriter = multiWriter

	engine := s.registerRouter()

	go func() {
		err := engine.RunListener(s.dockerOpt.listener)
		if err != nil {
			log.Error(err)
		}
	}()
	return nil
}

func (s *dockerPluginServer) Close() error {
	err := s.dockerOpt.authLog.Close()
	if err != nil {
		return err
	}

	err = s.dockerOpt.pluginLog.Close()
	if err != nil {
		return err
	}
	return nil
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
	val, ok := s.dockerOpt.policies.Load(string(dockerPluginAction))
	if !ok {
		return s.allowAuthResp()
	}
	policy := val.(Policy)
	var (
		eventListCh <-chan []*event.Event
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
		handleDockerPluginReportEvents(eventListCh, policy, s.dockerOpt.pluginLog)
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
