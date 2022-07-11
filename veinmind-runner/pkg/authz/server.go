package authz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/plugins"
	"github.com/gin-gonic/gin"
)

const (
	pluginSockPath = "/run/docker/plugins/veinmind-broker.sock"
)

type AuthZServer struct {
	listener   net.Listener
	authorizer Authorizer
	logPath    string
}

func NewAuthZServer(config PolicyConfig) (*AuthZServer, error) {
	s := new(AuthZServer)
	s.Init(config)
	return s, nil
}
func (a *AuthZServer) Init(cfg PolicyConfig) error {
	listener, err := net.ListenUnix("unix", &net.UnixAddr{Net: "unix", Name: pluginSockPath})
	if err != nil {
		return err
	}
	a.listener = listener
	a.logPath = cfg.Log.AuthZLogPath
	ba := new(BasicAuthorizer)
	err = ba.Init(&cfg)
	if err != nil {
		return err
	}
	a.authorizer = ba
	return nil
}
func (a *AuthZServer) Run() {
	router := gin.Default()
	// serverLogFile, err := GetLogFile(a.logPath)
	// if err != nil {
	// 	log.Error(err)
	// } else {
	// 	gin.DefaultWriter = io.MultiWriter(serverLogFile, os.Stdout)
	// 	defer serverLogFile.Close()
	// }

	router.Any("/Plugin.Activate", func(context *gin.Context) {
		context.JSON(200, plugins.Manifest{
			Implements: []string{authorization.AuthZApiImplements},
		})
	})

	router.Any(fmt.Sprintf("/%s", authorization.AuthZApiRequest), func(context *gin.Context) {
		b, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			log.Error(err)
			context.JSON(200, authorization.Response{
				Allow: true,
			})
			return
		}
		req := authorization.Request{}
		err = json.Unmarshal(b, &req)
		if err != nil {
			log.Error(err)
			context.JSON(200, authorization.Response{
				Allow: true,
			})
			return
		}

		resp := a.authorizer.AuthZReq(&req)
		context.JSON(200, resp)
	})

	router.Any(fmt.Sprintf("/%s", authorization.AuthZApiResponse), func(context *gin.Context) {
		b, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			log.Error(err)
			context.JSON(200, authorization.Response{
				Allow: true,
			})
			return
		}
		req := authorization.Request{}
		err = json.Unmarshal(b, &req)
		if err != nil {
			log.Error(err)
			context.JSON(200, authorization.Response{
				Allow: true,
			})
			return
		}

		resp := a.authorizer.AuthZRes(&req)
		context.JSON(200, resp)
	})

	err := router.RunListener(a.listener)
	if err != nil {
		panic(err)
	}
}
