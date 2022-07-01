package harborhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/registry"
	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
)

type WebhookServer struct {
	logPath      string
	dockerClient registry.Client
	authorizer   Authorizer
	// harborClient *apiv2.RESTClient
}

func NewWebhookServer(cfg WebhookConfig) (webhookServer WebhookServer, err error) {
	webhookServer.logPath = cfg.Log.WebhookLogPath
	authConfig := &registry.AuthConfig{
		Auths: []registry.Auth{cfg.DockerAuth}}
	webhookServer.authorizer.cfg = cfg
	webhookServer.dockerClient, err =
		registry.NewRegistryDockerClient(registry.WithAuthConfig(authConfig))
	if err != nil {
		return WebhookServer{}, err
	}
	return webhookServer, nil
}

func (t *WebhookServer) Run() {
	router := gin.Default()
	output := t.logPath
	if _, err := os.Stat(output); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(output)
		if err != nil {
			log.Error(err)
		}
	}
	f, err := os.OpenFile(output, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error(err)
	} else {
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		defer f.Close()
	}
	router.POST("/pullimage", func(c *gin.Context) {
		postData := &PullImage{}
		data, _ := ioutil.ReadAll(c.Request.Body)
		if err := json.Unmarshal(data, &postData); err != nil {
			fmt.Println(err)
		}
		fmt.Println("API TYPE => ", postData.Type)
		fmt.Println("iamge name => ", postData.EventData.Resources[0].ResourceURL)
		action := postData.Type
		imageName := postData.EventData.Resources[0].ResourceURL
		if postData.Operator == "webhook" {
			return
		}
		err = t.authorizer.CheckPull(action, imageName, *t)
		if err != nil {
			log.Error(err)
		}
	})

	router.POST("/api", func(c *gin.Context) {

		var body map[string]interface{}
		data, _ := ioutil.ReadAll(c.Request.Body)
		if err := json.Unmarshal(data, &body); err != nil {
			fmt.Println(err)
		}
		fmt.Println("body data => ", string(data))
		for k, v := range c.Request.Header {
			fmt.Println(k, v)
		}

	})
	router.Run()
}
