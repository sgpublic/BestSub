package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bestruirui/bestsub/internal/core/mihomo"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/setting"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

const (
	bestsubUpdateUrl    = "https://github.com/bestruirui/bestsub/releases/latest/download"
	bestsubUpdateApiUrl = "https://api.github.com/repos/bestruirui/BestSub/releases/latest"
	subcerUpdateApiUrl  = "https://api.github.com/repos/tindy2013/subconverter/releases/latest"
)

type LatestInfo struct {
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
	Message     string `json:"message"`
}

func GetLatestSubconverterInfo() (*LatestInfo, error) {
	return getLatestInfo(subcerUpdateApiUrl, op.GetSettingBool(setting.SUBCONVERTER_URL_PROXY))
}

func GetLatestBestsubInfo() (*LatestInfo, error) {
	return getLatestInfo(bestsubUpdateApiUrl, op.GetSettingBool(setting.PROXY_ENABLE))
}

func getLatestInfo(url string, proxy bool) (*LatestInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	hc := mihomo.Default(proxy)
	if hc == nil {
		return nil, fmt.Errorf("failed to create http client")
	}
	defer hc.Release()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Debugf("new request failed: %v", err)
		return nil, err
	}
	resp, err := hc.Do(req)
	if err != nil {
		log.Debugf("request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Debugf("read body failed: %v", err)
		return nil, err
	}

	latestInfo := LatestInfo{}
	err = json.Unmarshal(body, &latestInfo)
	if err != nil {
		log.Debugf("unmarshal body failed: %v", err)
		return nil, err
	}
	if latestInfo.Message != "" {
		return nil, fmt.Errorf("failed to get latest info: %s", latestInfo.Message)
	}
	return &latestInfo, nil
}

func download(url string, proxy bool) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	hc := mihomo.Default(proxy)
	if hc == nil {
		return nil, fmt.Errorf("failed to create http client")
	}
	defer hc.Release()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Debugf("new request failed: %v", err)
		return nil, err
	}
	resp, err := hc.Do(req)
	if err != nil {
		log.Debugf("request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Debugf("read body failed: %v", err)
		return nil, err
	}
	return bytes, nil
}
