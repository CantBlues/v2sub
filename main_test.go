package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CantBlues/v2sub/core"
	"github.com/CantBlues/v2sub/template"
)

func Test(t *testing.T) {

	fmt.Println("test starting ")
	// core.LoadConf()
	// nodes := core.GetNodes()
	// for _, node := range nodes {
	// 	core.NodesQueue.Enqueue(node)
	// 	fmt.Println(core.NodesQueue.Items)
	// }
	// core.SaveConf()

	rule := genRouteRule("proxy", []string{"baidu.com", "test.com"})
	v2ray := template.V2rayDefault
	v2ray.RouterConfig.Rules = append(v2ray.RouterConfig.Rules, rule)
	data, err := json.Marshal(v2ray)
	if err != nil {
		fmt.Println(err)
	}
	core.WriteFile("./v2ray.json", data)

}

func genRouteRule(tag string, domains []string) []byte {
	domain, err := json.Marshal(domains)
	if err != nil {
		print(err)
	}
	template := `{
			"type": "field",
			"outboundTag": "%s",
			"domain": %s			
		}`
	result := fmt.Sprintf(template, tag, domain)
	return []byte(result)

}
