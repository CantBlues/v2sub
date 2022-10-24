package main

import (
	"encoding/json"
	"fmt"
	"github.com/arkrz/v2sub/core"
	"github.com/arkrz/v2sub/types"
	"net/http"
	"strconv"
)

var (
	subCfg *types.Config
)

func main() {

	core.LoadConf()
	core.DisableIptable()

	subCfg = core.SubCfg

	http.HandleFunc("/fetch", fetch)
	http.HandleFunc("/change", change)
	http.HandleFunc("/close", closeIptable)
	http.HandleFunc("/start", startService)
	http.ListenAndServe(":89", nil)
	fmt.Println("listen")

}

func fetch(w http.ResponseWriter, r *http.Request) {
	refresh := r.URL.Query().Get("refresh")

	if len(subCfg.Nodes) == 0 || refresh != "" {
		subCfg.Nodes = core.GetNodes()
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
		subCfg.Current = subCfg.Nodes[nums]
		err := core.SwitchNode(subCfg.Nodes[nums])
		if err != nil {
			return
		}
		w.Write([]byte{'o', 'k'})
	}
}

func closeIptable(w http.ResponseWriter, r *http.Request) {
	core.DisableIptable()
	w.Write([]byte{'o', 'k'})
}

func startService(w http.ResponseWriter, r *http.Request) {
	core.StartService()
	w.Write([]byte{'o', 'k'})
}
