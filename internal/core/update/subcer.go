package update

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/setting"
	"github.com/bestruirui/bestsub/internal/modules/subcer"
	"github.com/bestruirui/bestsub/internal/utils/archives"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

func InitSubconverter() error {
	filePath := config.Base().SubConverter.Path + "/subconverter"
	if runtime.GOOS == "windows" {
		filePath += ".exe"
	}
	if _, err := os.Stat(filePath); err != nil {
		log.Infof("subconverter not found, downloading...")
		err = updateSubconverter()
		if err != nil {
			log.Warnf("auto update subconverter failed, please download subconverter manually from %s and move to %s: %v", op.GetSettingStr(setting.SUBCONVERTER_URL), config.Base().SubConverter.Path, err)
			os.Exit(1)
			return err
		}
		if _, err := os.Stat(filePath); err != nil {
			log.Warnf("subconverter not found, please download subconverter manually from %s and move to %s: %v", op.GetSettingStr(setting.SUBCONVERTER_URL), config.Base().SubConverter.Path, err)
			os.Exit(1)
			return err
		}
	}
	log.Infof("subconverter is already up to date")
	return nil
}

func UpdateSubconverter() error {
	log.Infof("start update subconverter")
	err := updateSubconverter()
	if err != nil {
		log.Warnf("update subconverter failed, please download subconverter manually from %s and move to %s: %v", op.GetSettingStr(setting.SUBCONVERTER_URL), config.Base().SubConverter.Path, err)
		return err
	}
	log.Infof("update subconverter success")
	return nil
}

func updateSubconverter() error {
	arch := runtime.GOARCH
	goos := runtime.GOOS

	var downloadUrl string
	baseUrl := op.GetSettingStr(setting.SUBCONVERTER_URL)

	var filename string
	switch goos {
	case "windows":
		switch arch {
		case "386":
			filename = "subconverter_win32.7z"
		case "amd64":
			filename = "subconverter_win64.7z"
		default:
			log.Errorf("unsupported windows architecture: %s", arch)
			return fmt.Errorf("unsupported windows architecture: %s", arch)
		}
	case "darwin":
		switch arch {
		case "amd64":
			filename = "subconverter_darwin64.tar.gz"
		case "arm64":
			filename = "subconverter_darwinarm.tar.gz"
		default:
			log.Errorf("unsupported darwin architecture: %s", arch)
			return fmt.Errorf("unsupported darwin architecture: %s", arch)
		}
	case "linux":
		switch arch {
		case "386":
			filename = "subconverter_linux32.tar.gz"
		case "amd64":
			filename = "subconverter_linux64.tar.gz"
		case "arm":
			filename = "subconverter_armv7.tar.gz"
		case "arm64":
			filename = "subconverter_aarch64.tar.gz"
		default:
			log.Errorf("unsupported linux architecture: %s", arch)
			return fmt.Errorf("unsupported linux architecture: %s", arch)
		}
	default:
		log.Errorf("unsupported operating system: %s", goos)
		return fmt.Errorf("unsupported operating system: %s", goos)
	}

	downloadUrl = baseUrl + "/" + filename

	bytes, err := download(downloadUrl, op.GetSettingBool(setting.SUBCONVERTER_URL_PROXY))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(config.Base().SubConverter.Path, 0755); err != nil {
		log.Errorf("failed to create directory: %v", err)
		return err
	}
	subcer.Lock()
	defer subcer.Unlock()
	subcer.Stop()
	cachePath := filepath.Join(config.Base().SubConverter.Path, "tmp")
	if err := archives.Extract(bytes, cachePath); err != nil {
		return err
	}

	srcDir := filepath.Join(cachePath, "subconverter")
	destDir := config.Base().SubConverter.Path
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		log.Debugf("read dir failed: %v", err)
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		err := os.Rename(srcPath, destPath)
		if err != nil {
			log.Debugf("move file failed: %v", err)
			return err
		}
	}
	if err := os.RemoveAll(cachePath); err != nil {
		return err
	}

	subcer.Start()
	return nil
}
