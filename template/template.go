package template

import (
	"encoding/json"
	"github.com/arkrz/v2sub/types"
)

const (
	ListenOnLocalAddr = "127.0.0.1"
	ListenOnWanAddr   = "0.0.0.0"

	ListenOnSocksProtocol = "socks"
	ListenOnSocksPort     = 1081

	ListenOnHttpProtocol = "http"
	ListenOnHttpPort     = 1082

	SocksPort = 12345
)

var ConfigTemplate = &types.Config{
	SubUrl: "",
	Nodes:  types.Nodes{},
	Current: 0,
}

// V2ray Default struct
var V2rayDefault = &types.V2ray{
	OutboundConfigs: []types.OutboundConfig{
		{
			Protocol: "blackhole",
			Tag:"blocked",
		},

	},
	InboundConfigs: []types.InboundConfig{
		{
			Protocol: "dokodemo-door",
			Port:     SocksPort,
			Sniffing: &types.Sniffing{
				Enabled:      true,
				DestOverride: []string{"http", "tls"},
			},
			Settings: &types.InBoundSetting{
				Network:        "tcp",
				FollowRedirect: true,
			},
			StreamSetting: &types.StreamSetting{
				Sockopt: &types.Sockopt{Tproxy: "redirect"},
			},
		},
	},
	RouterConfig: DefaultRouterConfigs,
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
	Strategy: "rules",
	RouteSetting: &types.RouteSetting{
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
					"outboundTag": "blocked",
					"domain": [
						"ext:site.dat:ad"
					]
				}`),
		},
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
