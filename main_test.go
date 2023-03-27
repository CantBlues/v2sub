package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CantBlues/v2sub/types"

	// "github.com/CantBlues/v2sub/core"
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
	source := template.V2rayDefault
	data, _ := json.Marshal(source)

	var copy types.V2ray

	json.Unmarshal(data, &copy)

	copy.InboundConfigs[0].Port = 666

	fmt.Println(copy)
	fmt.Println(source)

}
