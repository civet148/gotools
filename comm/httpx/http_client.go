package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type UrlValues map[string]interface{} //POST提交表单参数

const (
	HTTP_METHOD_GET     = "GET"     //请求指定的页面信息，并返回实体主体。
	HTTP_METHOD_POST    = "POST"    //向指定资源提交数据进行处理请求（例如提交表单或者上传文件）。数据被包含在请求体中。POST 请求可能会导致新的资源的建立和/或已有资源的修改。
	HTTP_METHOD_PUT     = "PUT"     //从客户端向服务器传送的数据取代指定的文档的内容。
	HTTP_METHOD_CONNECT = "CONNECT" //HTTP/1.1 协议中预留给能够将连接改为管道方式的代理服务器
	HTTP_METHOD_OPTIONS = "OPTIONS" //允许客户端查看服务器的性能。
	HTTP_METHOD_DELETE  = "DELETE"  //请求服务器删除指定的页面内容。
	HTTP_METHOD_TRACE   = "TRACE"   //回显服务器收到的请求，主要用于测试或诊断。
	HTTP_METHOD_PATCH   = "PATCH"   //是对PUT方法的补充，用来对已知资源进行局部更新。
	HTTP_METHOD_HEAD    = "HEAD"    //类似于GET请求，只不过返回的响应中没有具体的内容，用于获取报头
)

type Client struct {
	cli    *http.Client
	header *Header
}

type Response struct {
	StatusCode  int
	ContentType string
	Body        string
}

//新建一个HTTP客户端对象
//timeout   连接超时秒数
func NewHttpClient(timeout int) *Client {

	return newClient(timeout)
}

//新建一个HTTPS客户端对象
//timeout   连接超时秒数
//cer       PEM证书文件路径(string)或二进制数据([]byte)
func NewHttpsClient(timeout int, cer interface{}) *Client {

	return newClient(timeout, cer)
}

func newClient(timeout int, args ...interface{}) (c *Client) {

	if timeout <= 0 {
		timeout = 3
	}
	c = &Client{
		header: &Header{
			values: map[string]string{
				HEADER_KEY_CONTENT_TYPE: CONTENT_TYPE_NAME_X_WWW_FORM_URL_ENCODED,
			},
		},
		cli: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
	return
}

func (c *Client) Header() *Header {

	return c.header
}

//发起HTTP GET请求(send a http request by GET method)
func (c *Client) Get(strUrl string, values UrlValues) (response *Response, err error) {

	return c.do(HTTP_METHOD_GET, strUrl, values)
}

//发起HTTP POST请求(send a http request by POST method with application/x-www-form-urlencoded)
func (c *Client) PostUrlEncoded(strUrl string, values UrlValues) (response *Response, err error) {

	return c.do(HTTP_METHOD_POST, strUrl, values)
}

//发起HTTP POST请求(send a http request by POST method with content-type multipart/form-data or application/x-www-form-urlencoded)
//data 类型 string、[]byte、UrlValues或结构体[对象或引用]
func (c *Client) Post(strUrl string, data interface{}) (response *Response, err error) {

	return c.do(HTTP_METHOD_POST, strUrl, data)
}

//发起HTTP PUT请求(send a http request by PUT method)
func (c *Client) Put(strUrl string) (response *Response, err error) {

	return c.do(HTTP_METHOD_PUT, strUrl, nil)
}

//发起HTTP DELETE请求(send a http request by DELETE method)
func (c *Client) Delete(strUrl string) (response *Response, err error) {

	return c.do(HTTP_METHOD_DELETE, strUrl, nil)
}

//发起HTTP TRACE请求(send a http request by TRACE method)
func (c *Client) Trace(strUrl string) (response *Response, err error) {

	return c.do(HTTP_METHOD_TRACE, strUrl, nil)
}

//发起HTTP PATCH请求(send a http request by PATCH method)
func (c *Client) Patch(strUrl string) (response *Response, err error) {

	return c.do(HTTP_METHOD_PATCH, strUrl, nil)
}

func (c *Client) do(strMethod, strUrl string, data interface{}) (response *Response, err error) {

	var body io.Reader

	if data != nil {

		switch data.(type) {
		case UrlValues:
			{
				kvs := url.Values{}
				values := data.(UrlValues)
				if values != nil {
					for k, v := range values {
						kvs.Add(k, fmt.Sprintf("%v", v))
					}
				}
				body = strings.NewReader(kvs.Encode())
				log.Debug("UrlValues -> [%+v]", data.(UrlValues))
			}
		case string:
			{
				body = strings.NewReader(data.(string))
				log.Debug("string -> [%s]", data.(string))
			}
		case []byte:
			{
				body = bytes.NewReader(data.([]byte))
				log.Debug("[]byte -> [%s]", data.([]byte))
			}
		default:
			{
				var jsonData []byte
				if jsonData, err = json.Marshal(data); err != nil {
					log.Error("can't marshal data to json, error [%v]", err.Error())
					return
				}
				body = bytes.NewReader(jsonData)
				log.Debug("object -> [%s]", jsonData)
			}
		}
	}

	if response, err = c.send(strMethod, strUrl, body); err != nil {
		return
	}
	return
}

func (c *Client) send(strMethod, strUrl string, body io.Reader) (response *Response, err error) {

	var req *http.Request
	var resp *http.Response

	if req, err = http.NewRequest(strMethod, strUrl, body); err != nil {
		return
	}

	for k, v := range c.header.values {
		req.Header.Set(k, v)
	}

	if resp, err = c.cli.Do(req); err != nil {
		return
	}

	defer resp.Body.Close()

	response = &Response{
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get(HEADER_KEY_CONTENT_TYPE),
	}

	var data []byte
	if data, err = ioutil.ReadAll(resp.Body); err != nil {

		return
	}

	response.Body = string(data)
	return
}
