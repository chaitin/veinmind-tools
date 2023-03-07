package veinfs

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	api "github.com/chaitin/libveinmind/go"
)

// veinfs
// walk regular files with more than file info

const (
	isLinux   = "linux"
	isWindows = "windows"
	isDarwin  = "darwin"
)

type WithFileInfo func(image api.Image, fileInfo *FileInfo) (*FileInfo, error)

type WalkFunc func(info *FileInfo, err error) error

type baseFileInfo struct {
	Path string
	Name string
	Ext  string
	Size uint64
}

var defaultMaxSizeBits = 50 * MB

func Walk(image api.Image, rootPath string, walkFunc WalkFunc) error {
	return image.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return walkFunc(nil, err)
		}

		if info.Size()*8 > int64(defaultMaxSizeBits) {
			return nil
		}

		// skip the dir file
		if info.IsDir() {
			return walkFunc(&FileInfo{
				baseFileInfo: baseFileInfo{
					Path: path,
				},
				Type: info.Mode().Type()}, nil)
		}

		// skip the non regular file
		if !info.Mode().IsRegular() {
			return nil
		}

		// new file info detail
		return walkFunc(NewFileInfo(image, path, info))
	})
}

func NewFileInfo(image api.Image, path string, info ...os.FileInfo) (*FileInfo, error) {
	switch runtime.GOOS {
	case isLinux:
	default:
		return nil, errors.New("not supported system")
	}

	if image == nil || path == "" {
		return nil, errors.New("args can't be empty or nil")
	}

	if !filepath.IsAbs(path) {
		return nil, errors.New("not supported path")
	}

	var fileInfo os.FileInfo
	if len(info) > 0 {
		fileInfo = info[0]
	}
	if fileInfo == nil {
		stat, err := image.Stat(path)
		if err != nil {
			return nil, err
		}
		fileInfo = stat
	}

	fillFileInfo := withFileInfo(fileInfo)
	return fillFileInfo(image, &FileInfo{
		baseFileInfo: baseFileInfo{
			Path: path,
			Name: fileInfo.Name(),
			Ext:  filepath.Ext(path),
			Size: uint64(fileInfo.Size()),
		}})
}

func hashMD5(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

func hashSha256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
