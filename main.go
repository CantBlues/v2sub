package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CantBlues/remoteWake/router"
	"github.com/CantBlues/v2sub/core"
	"github.com/CantBlues/v2sub/ping"
	"github.com/CantBlues/v2sub/types"
)

var (
	subCfg *types.Config
)

func main() {
	// heartbeat("ip:port")
	core.LoadConf()
	core.DisableIptable()

	subCfg = core.SubCfg

	http.HandleFunc("/fetch", fetch)
	http.HandleFunc("/detect", detectNode)
	http.HandleFunc("/nodes/receive", receiveNode)
	http.HandleFunc("/change", change)
	http.HandleFunc("/iptable/toggle", toggleIptable)
	http.HandleFunc("/start", startService)
	http.HandleFunc("/iptable/status", checkIptable)
	http.ListenAndServe(":89", nil)
	fmt.Println("listen")

}

func heartbeat(server string) {
	go router.HeartBeatRoute(server)
}

func fetch(w http.ResponseWriter, r *http.Request) {
	refresh := r.URL.Query().Get("refresh")

	if len(subCfg.Nodes) == 0 || refresh != "" {
		subCfg.Nodes = core.GetNodes()
		go ping.Ping(subCfg.Nodes, core.Duration)
	}
	data, _ := json.Marshal(subCfg)
	w.Write(data)
}

func detectNode(w http.ResponseWriter, r *http.Request) {
	go http.Get("http://192.168.0.174:9999/v2ray/detect")
	w.Write([]byte{'1'})
}

func receiveNode(w http.ResponseWriter, r *http.Request) {
	var nodes types.Nodes
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&nodes)
	subCfg.Nodes = nodes

	
	d, _ := json.Marshal(subCfg)
	w.Write(d)
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

func toggleIptable(w http.ResponseWriter, r *http.Request) {
	if subCfg.FwStatus {
		core.DisableIptable()
		w.Write([]byte{'o', 'k'})
	} else {
		core.EnableIptable()
		w.Write([]byte{'o', 'k'})
	}
}

func startService(w http.ResponseWriter, r *http.Request) {
	core.StartService()
	w.Write([]byte{'o', 'k'})
}

func checkIptable(w http.ResponseWriter, r *http.Request) {
	if subCfg.FwStatus {
		w.Write([]byte{'1'})
	}
}
