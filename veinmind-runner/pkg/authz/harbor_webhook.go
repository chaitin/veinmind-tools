package authz

import (
	"context"
	"io"
	"os"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/action"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
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
		err := engine.Run(":" + s.option.port)
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
	engine.POST("/api", func(c *gin.Context) {
		if err := action.CheckPassword(c, s.option.password); err != nil {
			log.Error(err)
			return
		}
		postData, err := route.ParseHarborwebhookPostdata(c)
		if err != nil {
			log.Error(err)
			return
		}
		imageNames, err := route.GetImageNames(postData)
		if err != nil {
			log.Error(err)
			return
		}
		err = action.GetImagesFromHarbor(s.option.authInfo, imageNames)
		if err != nil {
			log.Error(err)
			return
		}
		val, ok := s.option.policies.Load(postData.Type)
		if !ok {
			log.Error(err)
			return
		}
		hpolicy := val.(HarborPolicy)
		var eventListCh chan []reporter.ReportEvent
		switch postData.Type {
		case "PUSH_ARTIFACT":
			eventListCh, err = HandleWebhookImagePush(context.Background(), hpolicy.Policy, postData)
		// TODO: other type's process
		default:
			return
		}
		if err != nil {
			log.Error(err)
			return
		}
		go func() {
			handleHarborWebhookReportEvents(eventListCh, hpolicy,
				s.option.pluginLog, s.option.mailConf)
		}()

	})
	// //get post data content from this url
	// engine.POST("/", func(c *gin.Context) {
	// 	var body map[string]interface{}
	// 	data, _ := ioutil.ReadAll(c.Request.Body)
	// 	if err := json.Unmarshal(data, &body); err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println("body data => ", string(data))
	// 	for k, v := range c.Request.Header {
	// 		fmt.Println(k, v)
	// 	}
	// 	c.JSON(http.StatusOK, struct{}{})
	// })
	return engine
}
