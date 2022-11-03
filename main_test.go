package main

import (
	// "encoding/json"
	"fmt"
	"os/exec"
	"testing"


	"github.com/CantBlues/v2sub/core"
	"github.com/CantBlues/v2sub/ping"
)

func Test(t *testing.T) {

	fmt.Println("test starting ")
	core.SubCfg, _ = core.ReadConfig("./v2sub.json")

	nodes := core.GetNodes()
	// nodes := core.SubCfg.Nodes
	// core.TestDelay(nodes[2])

	core.PrintAsTable(nodes)
	// node := nodes[10]
	// conf.OutboundConfigs = core.SetOutbound(node)
	// bytes, _ := json.Marshal(core.SubCfg)
	// core.WriteFile("./v2sub.json", bytes)

}

func TestGetPort(t *testing.T) {
	port, err := ping.GetFreePort()

	fmt.Println(port, err)
	port, err = ping.GetFreePort()

	fmt.Println(port, err)

}

func TestSpeed(t *testing.T){
	nodes := core.GetNodes()
	ping.TestAll(nodes)
}

func _TestGetPid(t *testing.T) {
	t.Log("test get command pid")
	cmd := exec.Command("./v2ray.exe", "run", "-c", "./v2ray.json")
	cmd.Start()
	pid := cmd.Process.Pid
	fmt.Println(pid)
	// cmd.Wait()
	fmt.Println("get pid then")

}
