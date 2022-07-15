package route

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var regexImageNamePattern = `images/(.*?)/push`

func GetImageNameFromUri(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	path := u.Path
	// post url like : /v1.41/images/{name}/push
	// ref:https://docs.docker.com/engine/api/v1.41/#tag/Image/operation/ImageHistory
	reg := regexp.MustCompile(regexImageNamePattern)
	imageNames := reg.FindStringSubmatch(path)
	if len(imageNames) != 2 {
		return "", errors.New("fail to parse url path")
	}
	kv, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	// tag as a query params
	tags, ok := kv["tag"]
	if !ok {
		return "", errors.New("bad request format")
	}
	imageName := imageNames[1] + ":" + tags[0]
	return imageName, nil
}

func GetImageNameFromUrlParam(uri, key string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	// post url like : /v1.41/images/create?fromImage=huzai9527%2Fbabyrop&tag=latest
	// ref:https://docs.docker.com/engine/api/v1.41/#tag/Image/operation/ImageCreate
	kv, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	imageNames, ok := kv[key]
	if !ok {
		return "", errors.New("bad request format")
	}
	tags, ok := kv["tag"]
	if !ok {
		return "", errors.New("bad request format")
	}
	imageName := imageNames[0] + ":" + tags[0]
	return imageName, nil
}

var regexIDSPattern = `^[0-9a-f]+$`

func GetImageNameFromBodyParam(uri, contentType, key string, reqBody []byte) (string, error) {
	// parse request body from URL and get image's name
	var body map[string]interface{}
	if contentType == "application/json" && len(reqBody) > 0 {
		if err := json.Unmarshal(reqBody, &body); err != nil {
			return "", err
		}
	}
	// imageName should be tagged
	// if you don't hava Image Key, then return defaultResponse
	imageNameI, ok := body[key]
	if !ok {
		return "", errors.New("Bad request head format")
	}

	imageName, ok := imageNameI.(string)
	if !ok {
		return "", errors.New("Fail to parse image name with string")
	}
	// cause only get imagename of string type
	// it may be patial digest id, such as "92e0f4bd4b90"
	// besides, it may only a string like "tocmat"
	// in this way, docker will take "tocmat" as "tocmat:latest"
	// it also may "tocmat:1.9" with tag or "tomcat:1.9@shaxxx" with tag and id
	// specially, if the imagename is "abcd" then we just return this to vemind-sdk
	// vemind-sdk will find all relative images whether "abcd" is name or id
	// reference : https://github.com/docker/docker-ce/blob/5d94ad617b913e7eaa5adb65dd6260d0aa87f9c9/components/engine/daemon/images/image.go#L150
	var HexRegexpAnchored = regexp.MustCompile(regexIDSPattern)
	if !HexRegexpAnchored.MatchString(imageName) {
		if strings.Count("imageName", ":") == 1 {
			imageName = imageName + ":latest"
		}
	}
	return imageName, nil
}
