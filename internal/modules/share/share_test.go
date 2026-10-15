package share

import (
	"encoding/base64"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestInlineClashProviders(t *testing.T) {
	config := []byte(`
proxies: []
proxy-providers:
  Provider_01:
    type: http
    url: https://internal.example/sub
proxy-groups:
  - name: Hong Kong
    type: url-test
    use: [Provider_01]
    filter: '^(?!.*CN)(?=.*HK).*'
    url: https://cp.cloudflare.com/generate_204
  - name: All Nodes
    type: select
    use: [Provider_01]
    proxies: [DIRECT]
rules:
  - MATCH,All Nodes
`)
	list := []byte(`
proxies:
  - name: HK Common
    type: ss
    server: hk.example
    port: 443
  - name: CN Common
    type: ss
    server: cn.example
    port: 443
`)

	result, err := inlineClashProviders(config, list)
	if err != nil {
		t.Fatalf("inlineClashProviders() error = %v", err)
	}

	var output map[string]any
	if err := yaml.Unmarshal(result, &output); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if _, ok := output["proxy-providers"]; ok {
		t.Fatal("output retained proxy-providers")
	}
	if got := proxyNames(output["proxies"]); !reflect.DeepEqual(got, []string{"HK Common", "CN Common"}) {
		t.Fatalf("proxies = %v, want all converted nodes", got)
	}

	groups := output["proxy-groups"].([]any)
	hongKong := groups[0].(map[string]any)
	if got := proxyNames(hongKong["proxies"]); !reflect.DeepEqual(got, []string{"HK Common"}) {
		t.Fatalf("Hong Kong proxies = %v, want [HK Common]", got)
	}
	if _, ok := hongKong["use"]; ok {
		t.Fatal("Hong Kong group retained use")
	}
	if _, ok := hongKong["filter"]; ok {
		t.Fatal("Hong Kong group retained filter")
	}

	allNodes := groups[1].(map[string]any)
	if got := proxyNames(allNodes["proxies"]); !reflect.DeepEqual(got, []string{"DIRECT", "HK Common", "CN Common"}) {
		t.Fatalf("All Nodes proxies = %v, want explicit nodes", got)
	}
}

func TestInlineClashProvidersWithoutProviderReturnsOriginal(t *testing.T) {
	config := []byte("proxies: []\nproxy-groups: []\n")
	result, err := inlineClashProviders(config, []byte("proxies: []\n"))
	if err != nil {
		t.Fatalf("inlineClashProviders() error = %v", err)
	}
	if !reflect.DeepEqual(result, config) {
		t.Fatalf("result = %q, want original config", result)
	}
}

func TestWithoutQueryKey(t *testing.T) {
	if got := withoutQueryKey("list=false&rename=Custom%20Name", "list"); got != "rename=Custom+Name" {
		t.Fatalf("withoutQueryKey() = %q, want rename query only", got)
	}
}

func TestDataSourceURL(t *testing.T) {
	data := []byte("proxies:\n")
	want := "data:text/yaml;base64," + base64.StdEncoding.EncodeToString(data)
	if got := dataSourceURL(data); got != want {
		t.Fatalf("dataSourceURL() = %q, want %q", got, want)
	}
}

func proxyNames(value any) []string {
	proxies, _ := value.([]any)
	names := make([]string, 0, len(proxies))
	for _, proxy := range proxies {
		if node, ok := proxy.(map[string]any); ok {
			names = append(names, stringValue(node["name"]))
			continue
		}
		names = append(names, stringValue(proxy))
	}
	return names
}
