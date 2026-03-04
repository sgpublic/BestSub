package archives

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

func UnTarGz(data []byte, dest string) error {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		log.Debugf("new gzip reader failed: %v", err)
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Debugf("read tar entry failed: %v", err)
			return err
		}

		fpath := filepath.Join(dest, hdr.Name)

		if !inDest(fpath, dest) {
			log.Debugf("invalid file path: %s", fpath)
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		switch hdr.Typeflag {

		case tar.TypeDir:
			if err = os.MkdirAll(fpath, os.ModePerm); err != nil {
				log.Debugf("mkdir all failed: %v", err)
				return err
			}

		case tar.TypeReg:
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				log.Debugf("mkdir parent failed: %v", err)
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				log.Debugf("open file failed: %v", err)
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tr)
			if err != nil {
				log.Debugf("copy failed: %v", err)
				return err
			}
		}
	}

	return nil
}
