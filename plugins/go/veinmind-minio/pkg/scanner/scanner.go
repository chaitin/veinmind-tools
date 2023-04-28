package scanner

import (
	"io"
	"io/fs"
	"os"
	"regexp"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
)

// MinIO RELEASE.2019-12-17T23-16-33Z <= MinIO Version < MinIO RELEASE.2023-03-20T20-16-18Z
const REGEX = `RELEASE\.(\d{4}-\d{2}-\d{2})T(\d{2}-\d{2}-\d{2})Z`

var startVersion, _ = time.Parse("2006-01-02", "2019-12-17")
var endVersion, _ = time.Parse("2006-01-02", "2023-03-20")

func ScanImage(image api.Image) Result {
	res := Result{}
	image.Walk("/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debug(err)
			return nil
		}
		if (info.Mode() & (os.ModeDevice | os.ModeNamedPipe | os.ModeSocket | os.ModeCharDevice | os.ModeDir)) != 0 {
			log.Debug("Skip: ", path)
			return nil
		}

		if info.Name() == "minio" {
			r := regexp.MustCompile(REGEX)
			file, err := image.Open(path)
			if err != nil {
				log.Debug(err)
				return nil
			}
			data, err := io.ReadAll(file)
			if err != nil {
				log.Debug(err)
				return nil
			}
			matches := r.FindAllStringSubmatch(string(data), -1)
			if len(matches) > 0 && len(matches[len(matches)-1]) > 1 {
				version := matches[len(matches)-1][1]
				v, err := time.Parse("2006-01-02", version)
				if err != nil {
					log.Debug(err)
					return nil
				}
				if v.Equal(startVersion) || (v.Before(endVersion) && v.After(startVersion)) {
					res.File = path
					res.Version = matches[len(matches)-1][0]
				}
			}
			return nil
		}
		return nil
	})

	return res
}

func ScanContainer(container api.Container) Result {
	res := Result{}
	container.Walk("/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debug(err)
			return nil
		}

		if (info.Mode() & (os.ModeDevice | os.ModeNamedPipe | os.ModeSocket | os.ModeCharDevice | os.ModeDir)) != 0 {
			log.Debug("Skip: ", path)
			return nil
		}

		if info.Name() == "minio" {
			r := regexp.MustCompile(REGEX)
			file, err := container.Open(path)
			if err != nil {
				log.Debug(err)
				return nil
			}
			data, err := io.ReadAll(file)
			if err != nil {
				log.Debug(err)
				return nil
			}
			matches := r.FindAllStringSubmatch(string(data), -1)
			if len(matches) > 0 && len(matches[len(matches)-1]) > 1 {
				version := matches[len(matches)-1][1]
				v, err := time.Parse("2006-01-02", version)
				if err != nil {
					log.Debug(err)
					return nil
				}
				if v.Equal(startVersion) || (v.Before(endVersion) && v.After(startVersion)) {
					res.File = path
					res.Version = matches[len(matches)-1][0]
				}
			}
			return nil
		}
		return nil
	})
	return res
}
