/*
 * @Author: thepoy
 * @Email: thepoy@163.com
 * @File Name: proxy.go
 * @Created: 2021-07-27 12:15:35
 * @Modified:  2021-11-05 14:12:39
 */

package predator

import (
	"net"
	"time"

	"github.com/thep0y/predator/proxy"
	"github.com/thep0y/predator/tools"
	"github.com/valyala/fasthttp"
)

// 可以从一些代理网站的 api 中请求指定数量的代理 ip
type AcquireProxies func(n int) []string

func (c *Crawler) DialWithProxy() fasthttp.DialFunc {
	return c.DialWithProxyAndTimeout(0)
}

func (c *Crawler) DialWithProxyAndTimeout(timeout time.Duration) fasthttp.DialFunc {
	return func(addr string) (net.Conn, error) {
		proxyAddr := tools.Shuffle(c.proxyURLPool)[0]
		c.log.Debug().Str("ProxyIP", proxyAddr).Msg("an proxy ip is selected from the proxy pool")
		if proxyAddr[:7] == "http://" || proxyAddr[:8] == "https://" {
			return proxy.HttpProxy(proxyAddr, addr, timeout)
		} else if proxyAddr[:9] == "socks5://" {
			return proxy.Socks5Proxy(proxyAddr, addr)
		} else {
			c.log.Fatal().Caller().
				Err(&proxy.ProxyErr{
					Code: proxy.ErrUnknownProtocolCode,
					Args: map[string]string{
						"proxy_addr": proxyAddr,
					},
					Msg: "only support http and socks5 protocol, but the incoming proxy address uses an unknown protocol",
				}).
				Str("proxy", proxyAddr).
				Send()
			return nil, nil
		}
	}
}
