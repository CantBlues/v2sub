package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/arkrz/v2sub/core"
	"github.com/arkrz/v2sub/template"
	"github.com/arkrz/v2sub/types"
)

const (
	v2subConfig = "./v2sub.json"
	v2rayConfig = "./v2ray.json"
	duration    = 5 * time.Second // 建议至少 5s
)

var (
	urls   = []string{"https://bulink.me/sub/3dpe3/vm", "https://g.luxury/link/oSDnm6qP5MdkyDvc?sub=4"}
	subCfg *types.Config
	v2ray  *types.V2ray
)

func main() {
	subCfg, _ = core.ReadConfig(v2subConfig)
	v2ray = template.V2rayDefault
	// http.HandleFunc("/fetch", fetch)
	// http.HandleFunc("/change", change)
	// http.ListenAndServe(":89", nil)
	fmt.Println("listen")

	nodes := core.GetNodes(urls)
	v2ray.OutboundConfigs = core.SetOutbound(nodes[34])

	data, _ := json.Marshal(v2ray)
	core.WriteFile(v2rayConfig, data)
	// core.PrintAsTable(nodes)
}

func fetch(w http.ResponseWriter, r *http.Request) {

	refresh := r.URL.Query().Get("refresh")

	if len(subCfg.Nodes) == 0 || refresh != "" {
		subCfg.Nodes = core.GetNodes(urls)
	}
	data, _ := json.Marshal(subCfg.Nodes)

	w.Write(data)
}

func change(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	nums, err := strconv.Atoi(target)
	if err != nil {
		return
	}

	if nums < subCfg.Nodes.Len() {
		subCfg.Current = nums
		core.SetOutbound(subCfg.Nodes[nums])
	}
}
