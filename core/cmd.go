package core

import (
	"bufio"
	"os"
	"os/exec"
)

const (
	V2rayService = "/etc/init.d/v2ray"
)

func EnableIptable() {
	SubCfg.FwStatus = true
	exec.Command("cp", SubCfg.IptableSource, SubCfg.IptablePath).Run()
	file, err := os.OpenFile(SubCfg.IptablePath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString("iptables -t nat -A PREROUTING -p tcp -j V2RAY")
	writer.Flush()
	exec.Command("fw3", "restart").Run()
}

func DisableIptable() {
	SubCfg.FwStatus = false
	exec.Command("cp", SubCfg.IptableSource, SubCfg.IptablePath).Run()
	exec.Command("fw3", "restart").Run()
}

func StartService() {
	exec.Command(V2rayService, "start").Run()
}

func RestartService() {
	if !SubCfg.FwStatus {
		EnableIptable()
	}
	exec.Command(V2rayService, "restart").Run()
}

func StartTestProcess(configPath string) (*exec.Cmd, error) {
	cmd := exec.Command("./v2ray.exe", "run", "-c", configPath)
	err := cmd.Start()
	return cmd, err
}

func CopyVlessJson(){
	cmd := exec.Command("cp", "/etc/config/v2ray/vless.json", "/etc/config/v2ray/config.json")
	cmd.Start()
}