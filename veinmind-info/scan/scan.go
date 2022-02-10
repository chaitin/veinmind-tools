package scan

import (
	"errors"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/veinmind-tools/veinmind-info/converter"
	common "github.com/chaitin/veinmind-tools/veinmind-info/log"
	"github.com/chaitin/veinmind-tools/veinmind-info/model"
	"github.com/chaitin/veinmind-tools/veinmind-info/scan/user"
)

type EngineType int

const (
	Dockerd EngineType = iota
	Containerd
)

var EngineTypeMap = map[string]EngineType{
	"dockerd":    Dockerd,
	"containerd": Containerd,
}

type ScanOption struct {
	EngineType EngineType
	ImageName  string
}

func Scan(opt ScanOption) (results []model.ImageInfo, err error) {
	// 初始化客户端
	var client api.Runtime

	switch opt.EngineType {
	case Dockerd:
		client, err = docker.New()
		if err != nil {
			return nil, err
		}

		defer func() {
			client.Close()
		}()
	case Containerd:
		client, err = containerd.New()

		if err != nil {
			return nil, err
		}

		defer func() {
			client.Close()
		}()
	default:
		return nil, errors.New("Engine type doesn't match")
	}

	var imageIds []string
	if opt.ImageName != "" {
		imageIds, err = client.FindImageIDs(opt.ImageName)
		if err != nil {
			return
		}
	} else {
		imageIds, err = client.ListImageIDs()
		if err != nil {
			return
		}
	}

	for _, imageID := range imageIds {
		scanResult, err := ScanById(imageID, client, opt)
		if err != nil {
			common.Log.Error(err)
			continue
		}

		results = append(results, scanResult)
	}

	return results, nil
}

func ScanById(id string, client api.Runtime, opt ScanOption) (result model.ImageInfo, err error) {
	image, err := client.OpenImageByID(id)
	if err != nil {
		return model.ImageInfo{}, err
	}

	// 初始化
	result = model.ImageInfo{}

	// 填充基本信息
	result.ID = image.ID()

	// 获取 OCI 信息
	err = converter.Convert(image, &result)
	if err != nil {
		common.Log.Error(err)
	}

	// 获取 User 信息
	err = user.GetUserInfo(image, &result)
	if err != nil {
		common.Log.Error(err)
	}

	return result, nil
}
