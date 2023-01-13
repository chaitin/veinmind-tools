package pkg

import (
	"bufio"
	"fmt"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"os"
	"strings"
)

func ContainerDockerAPiCheck(fs api.FileSystem) error {
	var file string
	if _, err := os.Open("/.dockerenv"); os.IsNotExist(err) {
		env := os.Getenv("Libveinmind-hostfs") //读取环境变量获取宿主机根目录挂载在容器内的哪个目录下，读取该目录下的/lib/systemd/system/docker.service获取docker的配置
		file = env + "/lib/systemd/system/docker.service"
	} else {
		file = "/host/lib/systemd/system/docker.service"
	}
	content, err := os.Open(file)
	if err != nil {
		log.Error(err)
		return err
	}

	defer FileClose(content, err)
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else {
			fmt.Println(scanner.Text())
			if strings.Contains(scanner.Text(), "-H tcp://") {
				AddResult(file, DOCKERAPIREASON, "Unsafe setting for Docker API :"+scanner.Text())
			}
		}
	}
	return nil

}
