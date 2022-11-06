package main

import (
	// "encoding/json"
	"fmt"
	"testing"

	"github.com/CantBlues/v2sub/core"
)

func Test(t *testing.T) {

	fmt.Println("test starting ")
	core.LoadConf()
	nodes := core.GetNodes()
	for _, node := range nodes {
		core.NodesQueue.Enqueue(node)
		fmt.Println(core.NodesQueue.Items)
	}
	core.SaveConf()

}
