package share

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/bestruirui/bestsub/internal/core/mihomo"
	"github.com/bestruirui/bestsub/internal/core/node"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/setting"
	"github.com/bestruirui/bestsub/internal/models/share"
	"github.com/bestruirui/bestsub/internal/modules/subcer"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/country"
	"github.com/dlclark/regexp2"
	"github.com/google/go-querystring/query"
	"gopkg.in/yaml.v3"
)

func GenSubData(ctx context.Context, genConfigStr string, userAgent string, sourceData []byte, extraQuery string) ([]byte, error) {
	var genConfig share.GenConfig
	if err := json.Unmarshal([]byte(genConfigStr), &genConfig); err != nil {
		return nil, err
	}
	if len(sourceData) <= len(nodeData) {
		return nil, fmt.Errorf("no nodes available for this share")
	}
	subUrlParam, err := query.Values(genConfig.SubConverter)
	if err != nil {
		return nil, err
	}
	if genConfig.Proxy {
		subUrlParam.Add("config_proxy", op.GetSettingStr(setting.PROXY_URL))
	}
	subUrlParam.Add("url", dataSourceURL(sourceData))
	subUrlParam.Add("remove_emoji", "false")
	baseURL, err := subcer.GetBaseUrl()
	if err != nil {
		return nil, err
	}
	client := mihomo.Default(false)
	if client == nil {
		return nil, fmt.Errorf("failed to create HTTP client")
	}
	defer client.Release()
	body, err := getConvertedSubData(ctx, client.Client, baseURL, subUrlParam, userAgent, extraQuery)
	if err != nil {
		return nil, err
	}
	if !isClashTarget(genConfig.SubConverter.Target) {
		return body, nil
	}
	if !hasClashProviders(body) {
		return body, nil
	}

	listParams := cloneValues(subUrlParam)
	listParams.Set("list", "true")
	listBody, err := getConvertedSubData(ctx, client.Client, baseURL, listParams, userAgent, withoutQueryKey(extraQuery, "list"))
	if err != nil {
		return nil, err
	}
	return inlineClashProviders(body, listBody)
}

func dataSourceURL(data []byte) string {
	return "data:text/yaml;base64," + base64.StdEncoding.EncodeToString(data)
}

func getConvertedSubData(ctx context.Context, client *http.Client, baseURL string, params url.Values, userAgent string, extraQuery string) ([]byte, error) {
	requestURL := fmt.Sprintf("%s/sub?%s", baseURL, params.Encode())
	if extraQuery != "" {
		requestURL += "&" + extraQuery
	}
	request, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", userAgent)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("subconverter returned HTTP %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func isClashTarget(target string) bool {
	target = strings.ToLower(target)
	return target == "clash" || target == "clashr"
}

func cloneValues(values url.Values) url.Values {
	clone := make(url.Values, len(values))
	for key, value := range values {
		clone[key] = append([]string(nil), value...)
	}
	return clone
}

func withoutQueryKey(rawQuery, key string) string {
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return rawQuery
	}
	values.Del(key)
	return values.Encode()
}

func hasClashProviders(body []byte) bool {
	var config struct {
		Providers map[string]any `yaml:"proxy-providers"`
	}
	return yaml.Unmarshal(body, &config) == nil && len(config.Providers) > 0
}

func inlineClashProviders(configBody, listBody []byte) ([]byte, error) {
	var config map[string]any
	if err := yaml.Unmarshal(configBody, &config); err != nil {
		return nil, fmt.Errorf("parse Clash configuration: %w", err)
	}
	providers, ok := config["proxy-providers"].(map[string]any)
	if !ok || len(providers) == 0 {
		return configBody, nil
	}

	var nodeList struct {
		Proxies []map[string]any `yaml:"proxies"`
	}
	if err := yaml.Unmarshal(listBody, &nodeList); err != nil {
		return nil, fmt.Errorf("parse Clash node list: %w", err)
	}
	config["proxies"] = nodeList.Proxies

	providerNames := make(map[string]struct{}, len(providers))
	for name := range providers {
		providerNames[name] = struct{}{}
	}
	if groups, ok := config["proxy-groups"].([]any); ok {
		for _, groupValue := range groups {
			group, ok := groupValue.(map[string]any)
			if !ok || !usesProvider(group["use"], providerNames) {
				continue
			}
			proxyNames, err := matchingProxyNames(nodeList.Proxies, stringValue(group["filter"]))
			if err != nil {
				return nil, err
			}
			group["proxies"] = appendProxyNames(group["proxies"], proxyNames)
			delete(group, "use")
			delete(group, "filter")
		}
	}
	delete(config, "proxy-providers")

	result, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("encode inlined Clash configuration: %w", err)
	}
	return result, nil
}

func usesProvider(value any, providers map[string]struct{}) bool {
	uses, ok := value.([]any)
	if !ok {
		return false
	}
	for _, use := range uses {
		if _, ok := providers[stringValue(use)]; ok {
			return true
		}
	}
	return false
}

func matchingProxyNames(proxies []map[string]any, filter string) ([]string, error) {
	var expression *regexp2.Regexp
	var err error
	if filter != "" {
		expression, err = regexp2.Compile(filter, 0)
		if err != nil {
			return nil, fmt.Errorf("compile proxy group filter %q: %w", filter, err)
		}
	}

	names := make([]string, 0, len(proxies))
	for _, proxy := range proxies {
		name := stringValue(proxy["name"])
		if name == "" {
			continue
		}
		if expression != nil {
			matched, err := expression.MatchString(name)
			if err != nil {
				return nil, fmt.Errorf("match proxy group filter %q: %w", filter, err)
			}
			if !matched {
				continue
			}
		}
		names = append(names, name)
	}
	return names, nil
}

func appendProxyNames(value any, names []string) []any {
	proxies, _ := value.([]any)
	seen := make(map[string]struct{}, len(proxies)+len(names))
	for _, proxy := range proxies {
		seen[stringValue(proxy)] = struct{}{}
	}
	for _, name := range names {
		if _, ok := seen[name]; ok {
			continue
		}
		proxies = append(proxies, name)
		seen[name] = struct{}{}
	}
	return proxies
}

func stringValue(value any) string {
	stringValue, _ := value.(string)
	return stringValue
}

func GenNodeData(config string) []byte {
	var genConfig share.GenConfig
	if err := json.Unmarshal([]byte(config), &genConfig); err != nil {
		return nil
	}
	nodes := node.GetByFilter(genConfig.Filter)
	var result bytes.Buffer
	result.Write(nodeData)
	tmpl, err := renameTemplate.Parse(genConfig.Rename)
	if err != nil {
		return nil
	}
	var newName bytes.Buffer
	for i, node := range *nodes {
		newName.Reset()
		result.Write(dash)
		subTags := op.GetSubTagsByID(context.Background(), node.Base.SubId)
		simpleInfo := renameTmpl{
			SpeedUp:       node.Info.SpeedUp.Average(),
			SpeedDown:     node.Info.SpeedDown.Average(),
			Delay:         uint32(node.Info.Delay.Average()),
			Risk:          uint32(node.Info.Risk),
			Count:         uint32(i + 1),
			Country:       country.GetCountry(node.Info.Country),
			IP:            utils.Uint32ToIP(node.Info.IP),
			SubName:       op.GetSubNameByID(context.Background(), node.Base.SubId),
			SubTags:       fmt.Sprintf("<%s>", strings.Join(subTags, "|")),
			SubTagsOrigin: subTags,
		}
		tmpl.Execute(&newName, simpleInfo)
		result.Write(rename(node.Base.Raw, newName.Bytes()))
		result.Write(newLine)
	}
	return result.Bytes()
}

func rename(raw []byte, newName []byte) []byte {
	idx := bytes.Index(raw, serverDelim)
	if idx < 0 {
		return raw
	}
	out := make([]byte, 0, len(name)+len(newName)+len(raw)-idx)
	out = append(out, name...)
	out = append(out, newName...)
	out = append(out, raw[idx:]...)
	return out
}

var (
	name        = []byte("{name: ")
	serverDelim = []byte(", server:")

	nodeData = []byte("proxies:\n")
	newLine  = []byte("\n")
	dash     = []byte(" - ")
)

type renameTmpl struct {
	SpeedUp       uint32
	SpeedDown     uint32
	Delay         uint32
	Risk          uint32
	Country       country.Country
	Count         uint32
	IP            string
	SubName       string
	SubTags       string
	SubTagsOrigin []string
}

var renameTemplate = template.New("node").Funcs(template.FuncMap{
	"add": func(x, y uint32) uint32 {
		return x + y
	},
	"sub": func(x, y uint32) uint32 {
		return x - y
	},
	"div": func(x, y uint32) uint32 {
		if y == 0 {
			return 0
		}
		return x / y
	},
	"mod": func(x, y uint32) uint32 {
		if y == 0 {
			return 0
		}
		return x % y
	},
})
