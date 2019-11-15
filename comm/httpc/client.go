package httpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"bytes"
	"mime/multipart"
	"io"
	"os"
)

type HttpClient struct {
	Cli        http.Client
	header     http.Header
	streamFlag bool
}

// 只有一个参数
func NewHttpC(maxconns ...int) *HttpClient {
	h := &HttpClient{}
	h.SetStreamFlag(false)
	if len(maxconns) == 1 && maxconns[0] > 0 {
		mineTransport := &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			MaxIdleConnsPerHost: maxconns[0], //连接池个数
		}

		h.Cli = http.Client{
			Transport: mineTransport,
			Timeout:   600 * time.Second,
		}
		h.SetStreamFlag(true)
	}
	h.header = make(http.Header)

	return h
}

func NewHttpsC(maxconns int, crtpath ...string) *HttpClient {
	h := &HttpClient{}
	pool := x509.NewCertPool()
	h.SetStreamFlag(true)

	// 可添加多个证书
	for _, v := range crtpath {
		fmt.Println(fmt.Sprintf("crtpath [%v]", v))
		caCrt, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Println(fmt.Sprintf("append crt ReadFile[%s] err:[%s]", v, err.Error()))
			break
		}

		if ok := pool.AppendCertsFromPEM(caCrt); !ok {
			fmt.Println(fmt.Sprintf("AppendCertsFromPEM error"))
		}
	}

	mineTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost: maxconns, //连接池个数
		TLSClientConfig:     &tls.Config{RootCAs: pool, InsecureSkipVerify: true},
	}

	h.Cli = http.Client{
		Transport: mineTransport,
		Timeout:   50 * time.Second,
	}
	h.header = make(http.Header)

	return h
}

func (h *HttpClient) SetStreamFlag(v bool) *HttpClient {
	h.streamFlag = v
	return h
}

func (h *HttpClient) SetHeader(k, v string) *HttpClient {
	h.header.Set(k, v)
	return h
}


// http 客户端POST文件数据
// file 要发送的文件数据([]byte)
// url 地址
func (h *HttpClient) PostFileBytes(data []byte, url string) (rsp [] byte, err error) {
/*	bodyBuf := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	request.Header.Set("Connection", "Keep-Alive")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	defer resp.Body.Close()
	rsp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("http.Do failed,[err=%s][url=%s]", err, url)
	}
	return rsp, err*/

	var resp *http.Response
	bodyBuf := bytes.NewReader(data)

	contentType := "multipart/form-data"//内容类型
	resp, err = h.Cli.Post(url, contentType, bodyBuf)
	if err != nil {
		fmt.Println(fmt.Sprintf("Post file error [%v]", err))
		return nil, err
	}

	defer resp.Body.Close()
	rsp, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("Read response body error [%v]", err))
		return nil, err
	}
	return rsp, nil
}

// http 客户端POST文件数据
// file 要发送的文件数据([]byte)或文件路径名(string)
// url 地址
func (h *HttpClient) PostFile(file string, url string) (body [] byte, err error) {

	//var FileData []byte
	var resp *http.Response
	var fileWriter io.Writer
	var bodyBuf = &bytes.Buffer{}

	//switch file.(type){
	//case []byte:
	//	{
	//		FileData = file.([]byte)
	//	}
	//case string://参数是文件路径
	//	{
	//		strFileName = file.(string)
	//		//FileData, err = ioutil.ReadFile(strFileName)
	//		//if err != nil{
	//		//	log.Error("Read file [%v] error [%v]", strFileName, err)
	//		//	return nil, err
	//		//}
	//	}
	//default:
	//	{
	//		log.Error("Unsupport parameter type [%v]", file)
	//		return nil, fmt.Errorf("Unsupport parameter type [%v]", file)
	//	}
	//}

	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err = bodyWriter.CreateFormFile("upload", file)
	if err != nil {
		fmt.Println(fmt.Sprintf("error writing to buffer"))
		return nil, err
	}

	//打开文件句柄操作
	fh, err := os.Open(file)
	if err != nil {
		fmt.Println(fmt.Sprintf("error opening file '%v' err (%v)", file, err.Error()))
		return nil, err
	}

	//拷贝到文件缓存
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fmt.Println(fmt.Sprintf("Copy to writer error: %v", err))
		return nil, err
	}

	contentType := bodyWriter.FormDataContentType()//内容类型
	bodyWriter.Close()

	resp, err = h.Cli.Post(url, contentType, bodyBuf)
	if err != nil {
		fmt.Println(fmt.Sprintf("Post file error [%v]", err))
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("Read response body error [%v]", err))
		return nil, err
	}
	return body, nil
}

// http 客户端发送数据
// packet 要发送报文
// method 选择发送方式 post get
// url 地址
func (h *HttpClient) SendUpstream(packet string, method string, url string) ([]byte, error) {
	if !h.streamFlag {
		return nil, errors.New("init http client error")
	}
	// log.Debug("http client request : ", url, packet)
	request, err := http.NewRequest(method, url, strings.NewReader(packet))
	if err != nil {
		return nil, errors.New("NewRequest: " + err.Error())
	}
	request.Header = h.header

	response, err := h.Cli.Do(request)
	if err != nil {
		// log.Error("http client request error : ", url, err)
		return nil, errors.New("Do: " + err.Error())
	}
	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("ReadAll: " + err.Error())
	}

	// log.Debug("http client answer : ", string(buf))
	return buf, nil
}

func (h *HttpClient) HttpsPostForm(url string, data url.Values) ([]byte, error) {
	var resp *http.Response
	var err error

	if h.streamFlag {
		resp, err = h.Cli.PostForm(url, data)
	} else {
		resp, err = http.PostForm(url, data)
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (h *HttpClient) HttpGet(url string) ([]byte, error) {
	var resp *http.Response
	var err error

	if h.streamFlag {
		resp, err = h.Cli.Get(url)
	} else {
		resp, err = http.Get(url)
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

//Add by dragon 20170609 返回Cookie

func (h *HttpClient) GetCookie(packet string, method string, url string) (string, error) {
	if !h.streamFlag {
		return "", errors.New("init http client error")
	}
	// log.Debug("http client request : ", url, packet)
	request, err := http.NewRequest(method, url, strings.NewReader(packet))
	if err != nil {
		return "", errors.New("NewRequest: " + err.Error())
	}
	request.Header = h.header

	response, err := h.Cli.Do(request)
	if err != nil {
		// log.Error("http client request error : ", url, err)
		return "", errors.New("Do: " + err.Error())
	}
	defer response.Body.Close()
	Cookie := ""
	if response.Header != nil {
		if _, ok := response.Header["Set-Cookie"]; ok {
			Cookie = response.Header["Set-Cookie"][0]
		}
	}
	// log.Debug("http client answer : ", string(buf))
	return Cookie, nil
}

func NewHttpsK(maxconns int, svrCrt, cliCrt, cliKey string) *HttpClient {

	h := &HttpClient{}
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(svrCrt)
	if err != nil {
		fmt.Println(fmt.Sprintf("append crt ReadFile[%s] err:[%s]", svrCrt, err.Error()))
		return h
	}
	//ok := pool.AppendCertsFromPEM(caCrt)
	//log.Debug("pool AppendCertsFromPEM: %v", ok)
	cert, err1 := x509.ParseCertificate(caCrt)
	if err1 != nil {
		fmt.Println(fmt.Sprintf("parses certificate err:[%s]", err1.Error()))
		return h
	}
	pool.AddCert(cert)

	clientCrt, err2 := tls.LoadX509KeyPair(cliCrt, cliKey)
	if err2 != nil {
		fmt.Println(fmt.Sprintf("Loadx509keypair err:[%s]", err2.Error()))
		return h
	}
//
	mineTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			//RootCAs: pool,
			Certificates: []tls.Certificate{clientCrt},
			InsecureSkipVerify: true,
		},
	}

	h.Cli = http.Client{Transport: mineTransport}
	h.header = make(http.Header)

	return h
}

