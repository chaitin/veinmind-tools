package authz

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/pkg/auth"
	"github.com/chaitin/veinmind-common-go/runtime"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/action"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HarborWebhookOption func(hwopt *harborWebhookOption) error

type harborWebhookOption struct {
	authLog       io.WriteCloser
	pluginLog     io.WriteCloser
	policies      sync.Map
	authInfo      auth.Auth
	mailConf      MailConf
	WebhookServer WebhookServer
}

func WithMailServer(mailconf MailConf) HarborWebhookOption {
	return func(hwopt *harborWebhookOption) error {
		hwopt.mailConf = mailconf
		return nil
	}
}
func WithAuthInfo(auth auth.Auth) HarborWebhookOption {
	return func(hwopt *harborWebhookOption) error {
		hwopt.authInfo = auth
		return nil
	}
}

func WithHarborPolicy(hwpolicies ...HarborPolicy) HarborWebhookOption {
	return func(option *harborWebhookOption) error {
		for _, hwpolicy := range hwpolicies {
			option.policies.Store(hwpolicy.Action, hwpolicy)
		}

		return nil
	}
}

func WithHarborAuthLog(path string) HarborWebhookOption {
	return func(option *harborWebhookOption) error {
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

func WithHarborPluginLog(path string) HarborWebhookOption {
	return func(option *harborWebhookOption) error {
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

func WithWebhookServer(webhookServer WebhookServer) HarborWebhookOption {
	return func(hwopt *harborWebhookOption) error {
		hwopt.WebhookServer = webhookServer
		return nil
	}
}
func WithHarborServerOptions(options ...HarborWebhookOption) HarborWebhookOption {
	return func(s *harborWebhookOption) error {
		for _, option := range options {
			if err := option(s); err != nil {
				return err
			}
		}
		return nil
	}
}

type harborWebhookServer struct {
	defaultServer
	hwoption *harborWebhookOption
}

func NewHarborWebhookServer(options ...HarborWebhookOption) (harborWebhookServer, error) {
	option := &harborWebhookOption{}
	for _, opt := range options {
		if err := opt(option); err != nil {
			return harborWebhookServer{}, err
		}
	}
	return harborWebhookServer{hwoption: option}, nil
}

func (s *harborWebhookServer) Init() error {
	opts := make([]HarborWebhookOption, 0)
	if s.hwoption.authLog == nil {
		opts = append(opts, WithHarborAuthLog(defaultAuthLogPath))
	}
	if s.hwoption.pluginLog == nil {
		opts = append(opts, WithHarborPluginLog(defaultPluginPath))
	}
	if s.hwoption.WebhookServer == (WebhookServer{}) {
		opts = append(opts, WithWebhookServer(defaultWebHookServer))
	}
	if err := WithHarborServerOptions(opts...)(s.hwoption); err != nil {
		return err
	}
	return nil
}

func (s *harborWebhookServer) Start() error {
	multiWriter := io.MultiWriter(s.hwoption.authLog, os.Stdout)

	logger := logrus.New()
	logger.Out = multiWriter

	log.SetDefaultLogger(log.NewLogrus(logger))
	gin.DefaultWriter = multiWriter

	engine := s.registerRouter()

	go func() {
		port := ":" + strconv.Itoa(s.hwoption.WebhookServer.Port)
		err := engine.Run(port)
		if err != nil {
			log.Error(err)
		}
	}()
	return nil
}

func (s *harborWebhookServer) Close() error {
	err := s.hwoption.authLog.Close()
	if err != nil {
		return err
	}

	err = s.hwoption.pluginLog.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *harborWebhookServer) registerRouter() *gin.Engine {
	engine := gin.Default()
	engine.POST("/api", func(c *gin.Context) {
		apiHandler(c, *s.hwoption)
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

func apiHandler(c *gin.Context, option harborWebhookOption) {
	if err := checkPassword(c, option.WebhookServer.Authorization); err != nil {
		log.Error(err)
		return
	}

	postData, err := parseHarborwebhookPostdata(c)
	if err != nil {
		log.Error(err)
		return
	}

	if postData.Operator == "webhook" {
		return
	}
	imageNames, err := getImageNames(postData)
	if err != nil {
		log.Error(err)
		return
	}
	err = pullImagesFromHarbor(option.authInfo, imageNames)
	if err != nil {
		log.Error(err)
		return
	}
	val, ok := option.policies.Load(postData.Type)
	if !ok {
		log.Error(err)
		return
	}
	hpolicy := val.(HarborPolicy)
	var eventListCh chan []reporter.ReportEvent
	switch postData.Type {
	case action.PUSH_ARTIFACT:
		eventListCh, err = HandleWebhookImagePush(context.Background(), hpolicy.Policy, imageNames)
	// TODO: other type's process
	default:
		return
	}
	if err != nil {
		return
	}
	go func() {
		handleHarborWebhookReportEvents(eventListCh, hpolicy,
			option.pluginLog, option.mailConf)
	}()
}

// get secrect from Authorization field and check
func checkPassword(c *gin.Context, password string) error {
	if password == "" {
		return nil
	}
	if c.Request.Header.Get("Authorization") == password {
		return nil
	}
	return errors.New("error passowrd")
}

func parseHarborwebhookPostdata(c *gin.Context) (action.PullandPushData, error) {
	postData := &action.PullandPushData{}
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(data, &postData); err != nil {
		return action.PullandPushData{}, err
	}
	return *postData, nil
}
func getImageNames(data action.PullandPushData) ([]string, error) {
	resources := data.EventData.Resources
	if len(resources) < 1 {
		return []string{}, errors.New("no image choosed")
	}
	var imagenames []string
	for _, resource := range resources {
		imagenames = append(imagenames, resource.ResourceURL)
	}
	return imagenames, nil
}

// download relevant images
func pullImagesFromHarbor(authentity auth.Auth, imageNames []string) error {
	authConfig := auth.AuthConfig{
		Auths: []auth.Auth{authentity}}
	dockerclient, err := runtime.NewDockerClient(runtime.WithAuth(authConfig))
	if err != nil {
		return err
	}
	for _, img := range imageNames {
		_, err := dockerclient.Pull(img)
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return nil
}
