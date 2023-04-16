package service

import (
	"bufio"
	"io"
	"regexp"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
)

type vsftpdService struct {
	name     string
	filepath []string
}

func (i *vsftpdService) Name() string {
	return i.name
}

func (i *vsftpdService) FilePath() (paths []string) {
	return i.filepath
}

// GetRecords TODO: 从配置文件中解析DB文件名称，根据DB文件名称检测
func (i *vsftpdService) GetRecords(file io.Reader) (records []model.Record, err error) {
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	var data []byte
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}
		data = append(data, buffer[:n]...)
	}
	if len(data) < 128 {
		log.Error("service: vsftpd virtual_users.db format error")
		return records, nil
	}
	data = data[len(data)-128:]
	re := regexp.MustCompile("\x01(.+?)\x01(.+)")
	matches := re.FindSubmatch(data)
	if len(matches) == 3 {
		s := model.Record{}
		s.Username = string(matches[2])
		s.Password = string(matches[1])
		records = append(records, s)
	}
	return records, nil
}

func init() {
	mod := &vsftpdService{}
	ServiceMatcherMap["vsftpd"] = "plain"
	mod.name = "vsftpd"
	mod.filepath = []string{"/etc/vsftpd/virtual_users.db"}
	Register("ftp", mod)
}
