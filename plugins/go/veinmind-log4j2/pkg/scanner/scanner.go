package scanner

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
)

func ScanImage(image api.Image, result *[]*Result) error {
	err := image.Walk("/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debug(err)
			return nil
		}

		if (info.Mode() & (os.ModeDevice | os.ModeNamedPipe | os.ModeSocket | os.ModeCharDevice | os.ModeDir)) != 0 {
			log.Debug("Skip: ", path)
			return nil
		}
		return Scan(path, image.ID(), result, image.Open)
	})
	return err
}

func ScanContainer(container api.Container, result *[]*Result) error {
	err := container.Walk("/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debug(err)
			return nil
		}

		if (info.Mode() & (os.ModeDevice | os.ModeNamedPipe | os.ModeSocket | os.ModeCharDevice | os.ModeDir)) != 0 {
			log.Debug("Skip: ", path)
			return nil
		}
		return Scan(path, container.ID(), result, container.Open)
	})
	return err
}

func Scan(path string, id string, res *[]*Result, openFile func(path string) (api.File, error)) error {

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".jar" || ext == ".war" {
		log.Debugf("going to scan %s", path)
		file, err := openFile(path)
		if err != nil {
			log.Warnf("failed to open file %s, err: %v", path, err)
			return nil
		}
		result, err := ScanFile(file, path, 0)
		if err != nil {
			log.Warnf("failed to scan %s, err: %v", path, err)
			return nil
		}
		if result.Code == Vulnerable {
			*res = append(*res, result)
			log.Warnf("[Vulnerable] image: %s, file: %s", id, result.DisplayPath)
		}
		log.Debug("path %s scan result: %v", path, result)
	}
	return nil
}

func ScanFile(file api.File, path string, depth int) (*Result, error) {

	// 深度控制
	if depth >= JarDetectDepth {
		return nil, errors.New(fmt.Sprintln(file, " exceed jar detect depth"))
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	f, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, err
	}

	hasLookUp := false
	fixed := false
	refs := make([]string, 0, 3)

	for _, ff := range f.File {
		fname := strings.ToLower(ff.Name)
		ext := filepath.Ext(fname)
		if ext != ".class" && ext != ".jar" {
			continue
		}

		if strings.HasSuffix(fname, "log4j/core/lookup/jndilookup.class") {
			hasLookUp = true
		} else if strings.HasSuffix(fname, "log4j/core/net/jndimanager.class") {
			reader, err := f.Open(ff.Name)
			if err != nil {
				log.Warnf("failed to open %s in %s, err: %v", ff.Name, file, err)
				continue
			}
			buf := make([]byte, 1024000)
			_, _ = io.ReadFull(reader, buf)
			if bytes.Contains(buf, []byte("Invalid JNDI URI")) {
				fixed = true
			}
		} else if strings.HasSuffix(fname, "log4j/core/lookup/interpolator.class") {
			reader, err := f.Open(ff.Name)
			if err != nil {
				log.Warnf("failed to open %s in %s, err: %v", ff.Name, file, err)
				continue
			}
			buf := make([]byte, 1024000)
			_, _ = io.ReadFull(reader, buf)
			// > 2.0
			if bytes.Contains(buf, []byte("log4j.core.lookup.JndiLookup")) ||
				// 2.0
				bytes.Contains(buf, []byte("JNDI lookup class is not available")) {
				refs = append(refs, ff.Name)
			}
		} else if ext == ".jar" {
			reader, err := f.Open(ff.Name)
			if err != nil {
				log.Warnf("failed to open %s in %s, err: %v", ff.Name, file, err)
				continue
			}

			tmpFile, err := os.CreateTemp("", "extract_*.jar")
			if err != nil {
				log.Warnf("failed to create temp file, err: %v", err)
				continue
			}
			_, err = io.Copy(tmpFile, reader)
			if err != nil {
				log.Warnf("failed to copy file to extract nested jar, err: %v", err)
				_ = tmpFile.Close()
				_ = os.Remove(tmpFile.Name())
				continue
			}

			result, err := ScanFile(tmpFile, ff.Name, depth+1)
			if err != nil {
				log.Warnf("failed to scan nested jar %s in %s, err: %v", ff.Name, file, err)
				_ = tmpFile.Close()
				_ = os.Remove(tmpFile.Name())
				continue
			}

			//关闭文件句柄的占用
			_ = tmpFile.Close()
			err = os.Remove(tmpFile.Name())
			if err != nil {
				log.Warnf("failed to delete jar %s, err: %v", tmpFile.Name(), err)
				continue
			}
			if result.Code == Vulnerable {
				result.DisplayPath = fmt.Sprintf("%s -> %s", path, ff.Name)
				return result, nil
			}
		}
	}
	result := &Result{Code: NotDetected, File: stat.Name(), DisplayPath: path}
	if len(refs) > 0 {
		if fixed || !hasLookUp {
			result.Code = FixedVersion
			return result, nil
		}
		result.Code = Vulnerable
		return result, nil
	}
	return result, nil
}
