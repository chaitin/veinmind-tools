package authz

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type harborWebhookServer struct {
	option *serverOption
}

func NewHarborWebhook(opts ...ServerOption) Server {
	s := &harborWebhookServer{
		option: new(serverOption),
	}
	return newDefaultServer(s, opts...)
}

func (s *harborWebhookServer) init() error {
	return nil
}

func (s *harborWebhookServer) start() error {
	multiWriter := io.MultiWriter(s.option.authLog, os.Stdout)

	logger := logrus.New()
	logger.Out = multiWriter

	log.SetDefaultLogger(log.NewLogrus(logger))
	gin.DefaultWriter = multiWriter

	engine := s.registerRouter()

	go func() {
		err := engine.Run(s.option.port)
		if err != nil {
			log.Error(err)
		}
	}()
	return nil
}

func (s *harborWebhookServer) wait() error {
	return nil
}

func (s *harborWebhookServer) close() {
	err := s.option.authLog.Close()
	if err != nil {
		log.Error(err)
	}

	err = s.option.pluginLog.Close()
	if err != nil {
		log.Error(err)
	}
}

func (s *harborWebhookServer) registerRouter() *gin.Engine {
	engine := gin.Default()

	engine.POST("/pushimage", func(c *gin.Context) {
		postData, err := route.ParseHarborwebhookPostdata(c)
		if err != nil {
			log.Error(err)
			return
		}
		val, ok := s.option.policies.Load(postData.Type)
		if !ok {
			log.Error(err)
			return
		}
		policy := val.(Policy)
		eventListCh, err := HandleWebhookImagePush(context.Background(), policy, postData)
		if err != nil {
			log.Error(err)
			return
		}
		go func() {
			handleReportEvents(eventListCh, policy, s.option.pluginLog)
		}()

	})

	engine.POST("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct{}{})
	})
	return engine
}
