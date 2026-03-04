package archives

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/bodgit/sevenzip"
)

func Un7z(data []byte, dest string) error {
	r, err := sevenzip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Debugf("new 7z reader failed: %v", err)
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

		rc, err := f.Open()
		if err != nil {
			log.Debugf("open file failed: %v", err)
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
		if err != nil {
			log.Debugf("open file failed: %v", err)
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			log.Debugf("copy failed: %v", err)
			return err
		}
	}

	return nil
}
