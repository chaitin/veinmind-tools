//go:build linux

package veinfs

import (
	"bytes"
	"debug/elf"
	"errors"
	"io"
	"os"
	"strings"
	"syscall"
	"time"

	api "github.com/chaitin/libveinmind/go"
)

var withFileInfo = withLinuxFileInfo

type FileInfo struct {
	baseFileInfo
	Md5        string      // md5
	Sha256     string      // sha256
	Uid        uint32      // uid bits
	Gid        uint32      // gid bits
	ELF        bool        // elf file
	UserFile   bool        // user file
	Hidden     bool        // hidden file
	Temporary  bool        // temporary file
	Perm       os.FileMode // permission bits
	Type       os.FileMode // type bits: dir/symlink/namedPipe/socket/device/charDevice
	CreateTime time.Time   // create time
	ModifyTime time.Time   // modify time
	AccessTime time.Time   // access time
}

func withLinuxFileInfo(info os.FileInfo) WithFileInfo {
	return func(image api.Image, fileInfo *FileInfo) (*FileInfo, error) {
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return nil, errors.New("not supported linux attr")
		}

		imageFp, err := image.Open(fileInfo.Path)
		if err != nil {
			return nil, err
		}

		content, err := io.ReadAll(imageFp)
		if err != nil {
			return nil, err
		}
		imageFp.Close()

		// file hash
		fileInfo.Md5 = hashMD5(content)
		fileInfo.Sha256 = hashSha256(content)

		fp, err := elf.NewFile(bytes.NewReader(content))
		if err == nil {
			fp.Close()
			fileInfo.ELF = true
		}

		// user file
		if strings.HasPrefix(fileInfo.Path, "/home/") ||
			strings.HasPrefix(fileInfo.Path, "/root/") {
			fileInfo.UserFile = true
		}

		// hidden file
		if strings.HasPrefix(fileInfo.Name, ".") {
			fileInfo.Hidden = true
		}

		// temporary file
		if strings.HasPrefix(fileInfo.Path, "/tmp/") {
			fileInfo.Temporary = true
		}

		// other file info
		fileInfo.Uid = stat.Uid
		fileInfo.Gid = stat.Gid
		fileInfo.Perm = info.Mode().Perm()
		fileInfo.Type = info.Mode().Type()
		fileInfo.CreateTime = time.Unix(int64(stat.Ctim.Sec), 0)
		fileInfo.ModifyTime = time.Unix(int64(stat.Mtim.Sec), 0)
		fileInfo.AccessTime = time.Unix(int64(stat.Atim.Sec), 0)

		return fileInfo, nil
	}
}
