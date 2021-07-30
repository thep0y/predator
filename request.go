/*
 * @Author: thepoy
 * @Email: thepoy@163.com
 * @File Name: request.go (c) 2021
 * @Created: 2021-07-24 13:29:11
 * @Modified: 2021-07-30 17:55:02
 */

package predator

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"sync/atomic"

	pctx "github.com/thep0y/predator/context"
	"github.com/thep0y/predator/json"
	"github.com/valyala/fasthttp"
)

type Request struct {
	// 访问的链接
	URL string
	// 请求方法
	Method string
	// 请求头
	Headers *fasthttp.RequestHeader
	// 请求和响应之间共享的上下文
	Ctx pctx.Context
	// 请求体
	Body []byte
	// 唯一标识
	ID uint32
	// 中断本次请求
	abort bool
	// 基于原 crawler 重试或发出新请求
	crawler *Crawler
	// 重试计数器
	retryCounter uint32
}

// New 使用原始请求的上下文创建一个新的请求
func (r *Request) New(method, URL string, body []byte) *Request {
	return &Request{
		Method:  method,
		URL:     URL,
		Body:    body,
		Ctx:     r.Ctx,
		Headers: &fasthttp.RequestHeader{},
		ID:      atomic.AddUint32(&r.crawler.requestCount, 1),
		crawler: r.crawler,
	}
}

func (r *Request) Abort() {
	r.abort = true
}

func (r *Request) SetContentType(contentType string) {
	r.Headers.Set("Content-Type", contentType)
}

func (r *Request) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		r.Headers.Set(k, v)
	}
}

func (r Request) NumberOfRetries() uint32 {
	return r.retryCounter
}

func (r Request) Get(u string) error {
	return r.crawler.Get(u)
}

func (r Request) Post(URL string, requestData map[string]string, ctx pctx.Context) error {
	return r.crawler.Post(URL, requestData, ctx)
}

type cacheRequest struct {
	// 访问的链接
	URL string
	// 请求方法
	Method string
	// 请求体
	Body []byte
}

func (r Request) Marshal() ([]byte, error) {

	cr := &cacheRequest{
		URL:    r.URL,
		Method: r.Method,
		Body:   r.Body,
	}

	return json.Marshal(cr)
}

func (r Request) Hash() (string, error) {
	var s strings.Builder

	s.WriteString(r.URL)
	s.WriteString(r.Method)
	for _, b := range r.Body {
		s.WriteByte(b)
	}

	return fmt.Sprintf("%x", sha1.Sum([]byte(s.String()))), nil
}
