package archives

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

func inDest(fpath, dest string) bool {
	if rel, err := filepath.Rel(dest, fpath); err == nil {
		if filepath.IsLocal(rel) {
			return true
		}
	}
	return false
}

func Extract(data []byte, dest string) error {
	if len(data) < 6 {
		return fmt.Errorf("invalid archive data")
	}

	switch {

	// 7z
	case bytes.HasPrefix(data, []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}):
		log.Debugf("archive type detected: 7z")
		return Un7z(data, dest)

	// zip
	case bytes.HasPrefix(data, []byte{0x50, 0x4B, 0x03, 0x04}):
		log.Debugf("archive type detected: zip")
		return UnZip(data, dest)

	// gzip (tar.gz / tgz)
	case bytes.HasPrefix(data, []byte{0x1F, 0x8B}):
		log.Debugf("archive type detected: tar.gz")
		return UnTarGz(data, dest)

	default:
		log.Debugf("unsupported archive type")
		return fmt.Errorf("unsupported archive format")
	}
}
