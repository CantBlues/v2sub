package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/arkrz/v2sub/ping"
	"github.com/arkrz/v2sub/template"
	"github.com/arkrz/v2sub/types"
)

const (
	V2subConfig = "/etc/config/v2ray/v2sub.json"
)

var (
	SubCfg *types.Config
)

func LoadConf() {
	SubCfg, _ = ReadConfig(V2subConfig)

	DisableIptable()
}

func saveConf() {
	bytes, _ := json.Marshal(SubCfg)
	WriteFile(V2subConfig, bytes)
}

func GetNodes() types.Nodes {
	urls := SubCfg.SubUrl

	fmt.Println("开始解析订阅信息...")

	var nodes types.Nodes
	subCh := make(chan []string, 1)
	defer close(subCh)
	for i := 0; i < len(urls); i++ {
		go GetSub(urls[i], subCh)

		select {
		case <-time.After(duration):
			ExitWithMsg(fmt.Sprintf("%s 后仍未获取到订阅信息, 请检查订阅地址和网络状况", duration.String()), 0)

		case data := <-subCh:
			if data == nil {
				ExitWithMsg("base64 解码错误, 请核实订阅编码", 0)
			}
			node, data := ParseNodes(data)
			nodes = append(nodes, node...)
			if len(data) != 0 {
				fmt.Println("无法解析下列节点:")
				for i := range data {
					fmt.Println(data[i])
				}
			}
		}
	}

	ping.Ping(nodes, duration)
	SubCfg.Nodes = nodes
	saveConf()
	// sort.Sort(nodes)
	return nodes
}

func setOutbound(node *types.Node) []types.OutboundConfig {
	config := template.DefaultOutboundConfigs
	config = append(config, resolve(node))
	return config
}

func resolve(node *types.Node) types.OutboundConfig {
	var outbound types.OutboundConfig
	var v2rayOutboundProtocol string
	var outboundSetting interface{}
	var streamSetting types.StreamSetting // v2ray.streamSettings
	switch node.Protocol {
	case vmessProtocol:
		v2rayOutboundProtocol = vmessProtocol
		outboundSetting = &types.VnextOutboundSetting{VNext: []types.VNextConfig{
			{
				Address: node.Addr,
				Port:    parsePort(node.Port),
				Users: []types.VNextUser{{
					ID:       node.UID,
					Security: node.Type,
					AlterId:  parsePort(node.AID),
				}},
			},
		}}
		streamSetting.Network = node.Net
		streamSetting.Security = node.TLS
		streamSetting.WsStream = &types.WsStream{Path: node.Path, Header: &types.Header{Host: node.Host}}
		streamSetting.TlsStream = &types.TlsStream{ServerName: node.Host, AllowInsecure: true}

	case ssProtocol:
		v2rayOutboundProtocol = ssProtocol
		outboundSetting = &types.SSOutboundSetting{Servers: []types.SSServerConfig{
			{
				Address:  node.Addr,
				Port:     parsePort(node.Port),
				Method:   node.Type,
				Password: node.UID,
			},
		}}
		streamSetting.Network = "tcp"
		streamSetting.Security = "none"

	case trojanProtocol:
		v2rayOutboundProtocol = "trojan"
		outboundSetting = struct {
			Address  string `json:"address"`
			Port     int    `json:"port"`
			Password string `json:"password"`
		}{
			Address:  node.Addr,
			Port:     parsePort(node.Port),
			Password: node.UID,
		}
		streamSetting.Network = "tcp"
		streamSetting.Security = "tls"
		streamSetting.TlsStream = &types.TlsStream{AllowInsecure: true, ServerName: node.Host}

	default:
		ExitWithMsg("unexpected protocol: "+node.Protocol, 1)
	}

	setting, _ := json.Marshal(outboundSetting)
	var settingRaw json.RawMessage = setting

	outbound = types.OutboundConfig{
		Protocol:       v2rayOutboundProtocol,
		Settings:       &settingRaw,
		Tag:            "proxy",
		StreamSettings: &streamSetting,
	}

	return outbound
}

func SwitchNode(node *types.Node) error {
	v2ray := template.V2rayDefault
	v2ray.OutboundConfigs = setOutbound(node)
	data, _ := json.Marshal(v2ray)
	err := WriteFile(SubCfg.V2rayCfg, data)
	if err != nil {
		return errors.New("write file error")
	}
	SubCfg.Current = node
	saveConf()
	RestartService()
	return nil
}
