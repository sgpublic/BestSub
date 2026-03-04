package archives

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

func UnZip(data []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Debugf("new zip reader failed: %v", err)
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !inDest(fpath, dest) {
			log.Debugf("invalid file path: %s", fpath)
			return fmt.Errorf("invalid file path: %s", fpath)
		}
		info := f.FileInfo()
		if info.IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			log.Debugf("mkdir all failed: %v", err)
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode().Perm())
		if err != nil {
			err = os.Remove(fpath)
			if err != nil {
				log.Debugf("remove file failed: %v", err)
				return err
			}
			outFile, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Debugf("open file failed: %v", err)
				return err
			}
		}
		defer outFile.Close()
		rc, err := f.Open()
		if err != nil {
			log.Debugf("open file failed: %v", err)
			return err
		}
		_, err = io.Copy(outFile, rc)
		rc.Close()
		if err != nil {
			log.Debugf("copy failed: %v", err)
			return err
		}
	}
	return nil
}
