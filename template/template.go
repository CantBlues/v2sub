package template

import (
	"encoding/json"
	"github.com/CantBlues/v2sub/types"
)

const (
	ListenOnLocalAddr = "127.0.0.1"
	ListenOnWanAddr   = "0.0.0.0"

	SocksPort = 12345
)

var ConfigTemplate = &types.Config{
	FwStatus:      false,
	SubUrl:        []string{"https://bulink.me/sub/3dpe3/vm", "https://g.luxury/link/oSDnm6qP5MdkyDvc?sub=4"},
	Nodes:         types.Nodes{},
	IptablePath:   "/etc/firewall.user",
	IptableSource: "/etc/config/v2ray/firewall.user",
	ExecPath:      "/etc/config/v2ray/v2ray",
	V2rayCfg:      "/etc/config/v2ray/config.json",
}

// V2ray Default struct
var V2rayDefault = &types.V2ray{
	InboundConfigs: []types.InboundConfig{
		{
			Protocol: "dokodemo-door",
			Port:     SocksPort,
			Listen:   ListenOnWanAddr,
			Sniffing: &types.Sniffing{
				Enabled:      true,
				DestOverride: []string{"http", "tls"},
			},
			Settings: &types.InBoundSetting{
				Network:        "tcp",
				FollowRedirect: true,
			},
			StreamSetting: &types.StreamSetting{
				Network:  "tcp",
				Security: "none",
				Sockopt:  &types.Sockopt{Tproxy: "redirect"},
			},
		},
	},
	RouterConfig: DefaultRouterConfigs,
}

// GetTestCfg get v2ray config for speed test or latency time
func GetTestCfg(port uint32) *types.V2ray {
	var config = &types.V2ray{
		InboundConfigs: []types.InboundConfig{
			{
				Protocol: "socks",
				Listen:   ListenOnLocalAddr,
				Port:     port,
				Settings: &types.InBoundSetting{},
				Sniffing: &types.Sniffing{
					Enabled:      true,
					DestOverride: []string{"http", "tls"},
				},
			},
		},
	}
	return config
}

// DefaultDNSConfigs 默认路由规则
// 参考 https://toutyrater.github.io/routing/configurate_rules.html
var DefaultDNSConfigs = &types.DNSConfig{Servers: []json.RawMessage{
	[]byte(`"114.114.114.114"`),
	[]byte(
		`{
			"address": "1.1.1.1",
			"port": 53,
			"domains": [
				"geosite:geolocation-!cn"
			]
		}`),
}}

var DefaultRouterConfigs = &types.RouterConfig{
	DomainStrategy: "AsIs",

	Rules: []json.RawMessage{
		[]byte(
			`{
				"type": "field",
				"outboundTag": "proxy",
				"domain": [
					"ext:site.dat:gw"
				]
			}`),
		[]byte(
			`{
				"type": "field",
				"outboundTag": "block",
				"domain": [
					"ext:site.dat:ad"
				]
			}`),
	},
}

var DefaultOutboundConfigs = []types.OutboundConfig{
	{
		Protocol: "freedom",
		Tag:      "direct",
	},
	{
		Protocol: "blackhole",
		Tag:      "block",
	},
}
