package types

import (
	"encoding/json"
)

type Config struct {
	SubUrl        []string `json:"subUrl"`
	Nodes         Nodes    `json:"nodes"`
	Current       *Node    `json:"current"`
	IptablePath   string   `json:"iptable"`
	IptableSource string   `json:"ipSource"`
	FwStatus      bool     `json:"status"`
	V2rayCfg      string   `json:"v2rayConfig"`
	ExecPath      string   `json:"v2ray"`
}

type V2ray struct {
	DNSConfigs      *DNSConfig       `json:"dns"`
	RouterConfig    *RouterConfig    `json:"routing"`
	OutboundConfigs []OutboundConfig `json:"outbounds"`
	InboundConfigs  []InboundConfig  `json:"inbounds"`
}

type DNSConfig struct {
	Servers []json.RawMessage `json:"servers"`
}

type RouterConfig struct {
	Strategy     string        `json:"strategy"`
	RouteSetting *RouteSetting `json:"settings"`
}

type RouteSetting struct {
	Rules []json.RawMessage `json:"rules"`
}

type OutboundConfig struct {
	Protocol       string           `json:"protocol"`
	Settings       *json.RawMessage `json:"settings"`
	Tag            string           `json:"tag"`
	StreamSettings *StreamSetting   `json:"streamSettings"`
}

type InboundConfig struct {
	Protocol      string          `json:"protocol"`
	Port          uint32          `json:"port"`
	Sniffing      *Sniffing       `json:"sniffing"`
	Settings      *InBoundSetting `json:"settings"`
	StreamSetting *StreamSetting  `json:"streamSettings"`
}

type Sniffing struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}

type InBoundSetting struct {
	Network        string `json:"network"`
	FollowRedirect bool   `json:"followRedirect"`
}
type VnextOutboundSetting struct {
	VNext []VNextConfig `json:"vnext"`
}

type VNextConfig struct {
	Address string      `json:"address"`
	Port    int         `json:"port"`
	Tag     string      `json:"tag"`
	Users   []VNextUser `json:"users"`
}

type VNextUser struct {
	ID       string `json:"id"`
	Security string `json:"security"`
	AlterId  int    `json:"alterId"`
}

type SocksOutboundSetting struct {
	Servers []SocksServerConfig `json:"servers"`
}

type SocksServerConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type SSOutboundSetting struct {
	Servers []SSServerConfig `json:"servers"`
}

type SSServerConfig struct {
	Email    string `json:"email"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Method   string `json:"method"`
	Password string `json:"password"`
	OTA      bool   `json:"ota"`
	Level    int    `json:"level"`
}

type StreamSetting struct {
	Network  string `json:"network"`
	Security string `json:"security"`

	TcpStream  *TcpStream  `json:"tcpSettings"`
	TlsStream  *TlsStream  `json:"tlsSettings"`
	KcpStream  *KcpStream  `json:"kcpSettings"`
	WsStream   *WsStream   `json:"wsSettings"`
	QuicStream *QuicStream `json:"quicSettings"`

	Sockopt *Sockopt `json:"sockopt"`
}

type TlsStream struct {
	ServerName    string `json:"serverName"`
	AllowInsecure bool   `json:"allowInsecure"`
}

type WsStream struct {
	Path   string  `json:"path"`
	Header *Header `json:"headers"`
}
type Header struct {
	Host string `json:"Host"`
}

type QuicStream struct {
}
type TcpStream struct {
}
type KcpStream struct {
}

type Sockopt struct {
	Tproxy string `json:"tproxy"`
	Mark   int    `json:"mark"`
}

type Trojan struct {
	RunType    string   `json:"run_type"`
	LocalAddr  string   `json:"local_addr"`
	LocalPort  int      `json:"local_port"`
	RemoteAddr string   `json:"remote_addr"`
	RemotePort int      `json:"remote_port"`
	Password   []string `json:"password"`
}

type Node struct {
	Name     string      `json:"ps"`
	Addr     string      `json:"add"`
	Port     interface{} `json:"port"`
	UID      string      `json:"id"`
	Net      string      `json:"net"`
	Type     string      `json:"type"`
	Host     string      `json:"host"`
	TLS      string      `json:"tls"`
	Protocol string      `json:"protocol"`
	AID      interface{} `json:"aid"`
	Path     string      `json:"path"`

	Ping int `json:"-"`
}

type Nodes []*Node

func (ns Nodes) Len() int { return len(ns) }
func (ns Nodes) Less(i, j int) bool {
	switch {
	case ns[i].Ping == -1:
		return false
	case ns[j].Ping == -1:
		return true
	default:
		return ns[i].Ping < ns[j].Ping
	}
}
func (ns Nodes) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }

type TableRow struct {
	Index int
	Name  string
	Addr  string
	Port  int
	Ping  int
}
