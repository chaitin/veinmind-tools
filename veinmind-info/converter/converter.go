package converter

import (
	"errors"
	api "github.com/chaitin/libveinmind/go"
	common "github.com/chaitin/veinmind-tools/veinmind-info/log"
	"github.com/chaitin/veinmind-tools/veinmind-info/model"
)

func Convert(image api.Image, info *model.ImageInfo) (err error) {
	if info == nil {
		return errors.New("Image info can't be nil")
	}

	// Convert OCI
	oci, err := image.OCISpecV1()
	if err != nil {
		common.Log.Error(err)
	} else {
		info.Created = oci.Created
		info.Env = oci.Config.Env
		info.User = oci.Config.User
		info.WorkingDir = oci.Config.WorkingDir
		info.Cmd = oci.Config.Cmd
		info.Entrypoint = oci.Config.Entrypoint
		info.Volumes = oci.Config.Volumes
		info.ExposedPorts = oci.Config.ExposedPorts
	}

	return nil
}
