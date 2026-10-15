package share

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGetConvertedSubDataIncludesErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("no nodes found"))
	}))
	defer server.Close()

	_, err := getConvertedSubData(context.Background(), server.Client(), server.URL, url.Values{"target": {"clash"}}, "", "")
	if err == nil {
		t.Fatal("getConvertedSubData() error = nil, want HTTP error")
	}
	if got, want := err.Error(), "subconverter returned HTTP 400: no nodes found"; got != want {
		t.Errorf("getConvertedSubData() error = %q, want %q", got, want)
	}
}

func TestHasNodeData(t *testing.T) {
	if HasNodeData([]byte("proxies:\n")) {
		t.Fatal("HasNodeData() = true for an empty node document")
	}
	if !HasNodeData([]byte("proxies:\n - {name: node}\n")) {
		t.Fatal("HasNodeData() = false for a document with a node")
	}
}

func TestListConversionParamsUsesDataSource(t *testing.T) {
	nodeData := []byte("proxies:\n - {name: node}\n")
	params := listConversionParams(url.Values{"target": {"clash"}, "url": {"https://bestsub.example/node"}}, nodeData)

	if got, want := params.Get("list"), "true"; got != want {
		t.Errorf("list = %q, want %q", got, want)
	}
	wantURL := "data:text/yaml;base64," + base64.StdEncoding.EncodeToString(nodeData)
	if got := params.Get("url"); got != wantURL {
		t.Errorf("url = %q, want %q", got, wantURL)
	}
}

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
