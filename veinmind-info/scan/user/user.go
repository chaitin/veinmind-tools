package user

import (
	"bufio"
	api "github.com/chaitin/libveinmind/go"
	common "github.com/chaitin/veinmind-tools/veinmind-info/log"
	"github.com/chaitin/veinmind-tools/veinmind-info/model"
	"strings"
)

type Passwd struct {
	Name   string
	Passwd string
	Uid    string
	Gid    string
	Gecos  string
	Dir    string
	Shell  string
}

func GetUserInfo(image api.Image, info *model.ImageInfo) (err error) {
	f, err := image.Open("/etc/passwd")
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		userinfo := strings.Split(scanner.Text(), ":")
		if len(userinfo) != 7 {
			common.Log.Error("Passwd format error")
			continue
		}

		passwd := Passwd{}
		passwd.Name = userinfo[0]
		passwd.Passwd = userinfo[1]
		passwd.Uid = userinfo[2]
		passwd.Gid = userinfo[3]
		passwd.Gecos = userinfo[4]
		passwd.Dir = userinfo[5]
		passwd.Shell = userinfo[6]

		info.Users = append(info.Users, model.ImageUserInfo{
			Username:    passwd.Name,
			Uid:         passwd.Uid,
			Gid:         passwd.Gid,
			Shell:       passwd.Shell,
			Description: passwd.Gecos,
		})
	}

	return nil
}
