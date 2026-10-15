package subcer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bestruirui/bestsub/internal/core/mihomo"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/setting"
)

// GetBaseUrl returns the address of the independently managed Subconverter service.
func GetBaseUrl() (string, error) {
	value := strings.TrimRight(strings.TrimSpace(op.GetSettingStr(setting.SUBCONVERTER_SERVICE_URL)), "/")
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("invalid subconverter service URL: configure an http(s) base URL in settings")
	}
	return value, nil
}

func GetVersion() string {
	baseURL, err := GetBaseUrl()
	if err != nil {
		return ""
	}
	client := mihomo.Default(false)
	if client == nil {
		return ""
	}
	defer client.Release()
	client.Timeout = 5 * time.Second
	resp, err := client.Get(baseURL + "/version")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 256))
	if err != nil {
		return ""
	}
	parts := strings.Fields(string(body))
	if len(parts) > 1 {
		return parts[1]
	}
	return strings.TrimSpace(string(body))
}
