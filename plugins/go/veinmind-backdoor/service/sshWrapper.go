package service

import (
	"io"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

const ELFMagicByte = "\x7FELF"

func sshWrapperBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	sshdPath := "/usr/sbin/sshd"
	check := false
	var res []*event.BackdoorDetail

	fileInfo, err := apiFileSystem.Stat(sshdPath)
	if err != nil {
		return false, nil
	}
	file, err := apiFileSystem.Open(sshdPath)
	if err != nil {
		return false, nil
	}
	header := make([]byte, 4)
	_, err = io.ReadFull(file, header)
	if err != nil {
		return false, nil
	}

	// 判断文件头是否被修改
	if string(header) != ELFMagicByte {
		check = true
		fileDetail, err := file2FileDetail(fileInfo, sshdPath)
		if err != nil {
			return false, nil
		}
		res = append(res, &event.BackdoorDetail{
			FileDetail:  fileDetail,
			Content:     "sshd file is not ELF format",
			Description: "sshwrapper backdoor",
		})
	}

	return check, res
}

func init() {
	ImageCheckFuncMap["sshWrapper"] = sshWrapperBackdoorCheck
	ContainerCheckFuncMap["sshWrapper"] = sshWrapperBackdoorCheck
}
