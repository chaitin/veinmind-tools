package pkg

import (
	"bufio"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"os"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
)

func ContainerDockerAPiCheck(fs api.FileSystem) ([]*event.EscapeDetail, error) {
	var res = make([]*event.EscapeDetail, 0)
	var file string
	if _, err := os.Open("/.dockerenv"); os.IsNotExist(err) {
		env := os.Getenv("LIBVEINMIND_HOST_ROOTFS") //读取环境变量获取宿主机根目录挂载在容器内的哪个目录下，读取该目录下的/lib/systemd/system/docker.service获取docker的配置
		file = env + "/lib/systemd/system/docker.service"
	} else {
		file = "/host/lib/systemd/system/docker.service"
	}
	content, err := os.Open(file)
	if err != nil {
		log.Error(err)
		return res, err
	}

	defer content.Close()
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else {
			if strings.Contains(scanner.Text(), "-H tcp://") {
				res = append(res, &event.EscapeDetail{
					Target: file,
					Reason: DOCKERAPIREASON,
					Detail: "Unsafe setting for Docker API :" + scanner.Text(),
				})
			}
		}
	}
	return res, nil

}

func init() {
	ContainerCheckList = append(ContainerCheckList, ContainerDockerAPiCheck)
}
