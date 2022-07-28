package action

import (
	"errors"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/pkg/auth"
	"github.com/chaitin/veinmind-common-go/runtime"
	"github.com/gin-gonic/gin"
)

// get secrect from Authorization field and check
func CheckPassword(c *gin.Context, password string) error {
	if password == "" {
		return nil
	}
	if c.Request.Header.Get("Authorization") == password {
		return nil
	}
	return errors.New("error passowrd")
}

// download relevant images
func GetImagesFromHarbor(authentity auth.Auth, imageNames []string) error {
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
