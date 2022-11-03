package ping

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/CantBlues/v2sub/core"
	"github.com/CantBlues/v2sub/template"
	"github.com/CantBlues/v2sub/types"
	gop "github.com/sparrc/go-ping"
)

func Ping(nodes types.Nodes, duration time.Duration) {
	timer := time.After(duration)
	ch := make(chan [2]int, len(nodes))
	//defer close(ch)  后续写入会导致 panic

	for i := range nodes {
		nodes[i].Ping = -1

		go func(ch chan<- [2]int, index int) {
			pinger, err := gop.NewPinger(nodes[index].Addr)
			if err != nil {
				return // parse address error
			}

			pinger.Count = 4
			pinger.Interval = 500 * time.Millisecond
			pinger.SetPrivileged(true)
			pinger.OnFinish = func(stats *gop.Statistics) {
				ch <- [2]int{index, int(stats.AvgRtt.Nanoseconds() / 1e6)}
			}

			pinger.Run()
		}(ch, i)
	}

	for {
		select {
		case <-timer:
			return
		case res := <-ch:
			if res[1] != 0 {
				nodes[res[0]].Ping = res[1]
			}
		}
	}
}

func TestAll(nodes types.Nodes) {
	Ping(nodes, core.Duration)

	for i, node := range nodes {
		if node.Host == "127.0.0.1" || core.ParsePort(node.Port) <= 0 || node.Ping < 10 || node.Ping > 1000 {
			continue
		}
		TestNodeQuality(node)

		core.PrintAsTable(nodes)

		if i > 20 {
			break
		}

	}

}

func TestNodeQuality(node *types.Node) {

	port, err := GetFreePort()
	if err != nil {
		fmt.Println("failed to get free port")
		return
	}

	conf := template.GetTestCfg(uint32(port))
	conf.OutboundConfigs = core.SetOutbound(node)
	data, _ := json.Marshal(conf)
	fileName := fmt.Sprintf("./test%s.json", strconv.Itoa(port))
	core.WriteFile(fileName, data)
	cmd := core.StartTestProcess(fileName)
	node.Delay = httpPing(port)
	node.Speed = speedTest(port)
	fmt.Println(cmd.Process.Pid)
	cmd.Process.Kill()
	os.Remove(fileName)
}

func httpPing(port int) string {
	delay := "-"
	
	proxy, _ := url.Parse("socks5://127.0.0.1:" + strconv.Itoa(port))
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, Proxy: http.ProxyURL(proxy)}
	client := &http.Client{Transport: tr, Timeout: core.Duration}

	core.RetryDo(3, time.Second*2, func() error {
		start := time.Now()
		_, err := client.Get("https://google.com")
		if err != nil {
			fmt.Println("get http error", err)
			return err
		}
		end := time.Now()
		duration := end.Sub(start)
		delay = fmt.Sprint(duration)
		return nil
	})

	return delay
}

func speedTest(port int) string {
	speed := "-"

	proxy, _ := url.Parse("socks5://127.0.0.1:" + strconv.Itoa(port))
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, Proxy: http.ProxyURL(proxy)}
	client := &http.Client{Transport: tr, Timeout: core.Duration}

	core.RetryDo(3, time.Second*2, func() error {
		start := time.Now()
		_, err := client.Get("http://cachefly.cachefly.net/10mb.test")
		if err != nil {
			fmt.Println("get http error", err)
			return err
		}
		end := time.Now()
		duration := end.Sub(start)
		tmp := 10 / duration.Seconds()
		speed = fmt.Sprint(tmp)
		return nil
	})

	return speed
}

// 动态获取可用端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	l, err := net.Listen("tcp", addr.String())
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
