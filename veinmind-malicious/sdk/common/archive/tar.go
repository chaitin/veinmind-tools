package archive

import (
	"archive/tar"
	"io"
	"os"
	"path"
)

func Untar(reader io.Reader, directory string) error {
	tr := tar.NewReader(reader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if !ExistDir(path.Join(directory, hdr.Name)) {
				err := os.Mkdir(path.Join(directory, hdr.Name), 0755)
				if err != nil {
					continue
				}
			}
		case tar.TypeReg:
			f, err := os.Create(path.Join(directory, hdr.Name))
			if err != nil {
				continue
			}
			_, err = io.Copy(f, tr)
			f.Close()
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}
