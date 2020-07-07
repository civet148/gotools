package parser

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	URL_SCHEME_HTTPS     = "https"
	URL_SCHEME_HTTP      = "http"
	URL_SCHEME_MYSQL     = "mysql"
	URL_SCHEME_POSTGRES  = "postgres"
	URL_SCHEME_REDIS     = "redis"
	URL_SCHEME_ETCD      = "etcd"
	URL_SCHEME_KAFKA     = "kafka"
	URL_SCHEME_MSSQL     = "mssql"
	URL_SCHEME_ZOOKEEPER = "zookeeper"
	URL_SCHEME_SEP       = "://"
)

type UrlInfo struct {
	Scheme     string            // http...
	Host       string            // host name and port eg. '127.0.0.1:3306'
	User       string            // user
	Password   string            // password
	Path       string            // path
	Fragment   string            // fragment
	Opaque     string            // opaque
	ForceQuery bool              // force query
	Queries    map[string]string // queries
}

//URL have some special characters in password(支持URL中密码包含特殊字符)
func ParseUrl(strUrl string) (ui *UrlInfo) {

	ui = &UrlInfo{Queries: make(map[string]string, 1)}

	var encodes = map[string]string{
		"`":  "%60",
		"#":  "%23",
		"?":  "%3f",
		"<":  "%3c",
		">":  "%3e",
		"[":  "%5b",
		"]":  "%5d",
		"{":  "%7b",
		"}":  "%7d",
		"/":  "%2f",
		"|":  "%7c",
		"\\": "%5c",
		"%":  "%25",
		"^":  "%5e",
	}

	var decodes = map[string]string{
		"%60": "`",
		"%23": "#",
		"%3f": "?",
		"%3c": "<",
		"%3e": ">",
		"%5b": "[",
		"%5d": "]",
		"%7b": "{",
		"%7d": "}",
		"%2f": "/",
		"%7c": "|",
		"%5c": "\\",
		"%25": "%",
		"%5e": "^",
	}
	_ = decodes

	strUrl = strings.TrimSpace(strUrl)
	if strings.Contains(strUrl, "@") { // if a url have user+password, there must be have '@'
		// find first '://'
		var strScheme string
		_ = strScheme

		index := strings.LastIndex(strUrl, URL_SCHEME_SEP)
		if index > 0 {
			strScheme = strUrl[:index]
			strUrl = strUrl[index+3:]
		}

		// find last '@'
		index = strings.LastIndex(strUrl, "@")
		if index > 0 {
			strPrefix := strUrl[:index]
			strSuffix := strUrl[index:]
			for k, v := range encodes {
				//encode user and password special character(s) to url encode
				strPrefix = strings.ReplaceAll(strPrefix, k, v)
			}

			if strScheme != "" {
				strUrl = strScheme + URL_SCHEME_SEP
			}
			strUrl += strPrefix + strSuffix
		}
	}

	u, err := url.Parse(strUrl)
	if err != nil {
		return
	}
	ui.Path = u.Path
	ui.Host = u.Host
	ui.Scheme = u.Scheme
	ui.Fragment = u.Fragment
	ui.Opaque = u.Opaque
	ui.ForceQuery = u.ForceQuery

	if u.User != nil {
		ui.User = u.User.Username()
		ui.Password, _ = u.User.Password()
		for k, v := range decodes {
			//decode password from url encode to special character(s)
			ui.Password = strings.ReplaceAll(ui.Password, k, v)
		}
	}
	vs, _ := url.ParseQuery(u.RawQuery)
	for k, v := range vs {
		ui.Queries[k] = v[0]
	}
	return
}

func (u *UrlInfo) GetScheme() string {
	return u.Scheme
}

func (u *UrlInfo) GetHost() string {
	return u.Host
}

func (u *UrlInfo) GetIPAndPort() (ip string, port int) {

	ss := strings.Split(u.Host, ":")
	if len(ss) == 0 {
		return
	}

	if len(ss) == 1 {
		port = u.getSchemePort()
	} else {
		port, _ = strconv.Atoi(ss[1])
	}
	ip = ss[0]
	return
}

func (u *UrlInfo) GetPath() string {
	return u.Path
}

func (u *UrlInfo) GetUser() string {
	return u.User
}

func (u *UrlInfo) GetPassword() string {
	return u.Password
}

func (u *UrlInfo) GetOpaque() string {
	return u.Opaque
}

func (u *UrlInfo) GetFragment() string {
	return u.Fragment
}

func (u *UrlInfo) GetQueries() map[string]string {
	return u.Queries
}

func (u *UrlInfo) GetForceQuery() bool {
	return u.ForceQuery
}

func (u *UrlInfo) getSchemePort() (port int) {
	switch u.Scheme {
	case URL_SCHEME_HTTPS:
		port = 443
	case URL_SCHEME_HTTP:
		port = 80
	case URL_SCHEME_MYSQL:
		port = 3306
	case URL_SCHEME_POSTGRES:
		port = 5432
	case URL_SCHEME_REDIS:
		port = 6379
	case URL_SCHEME_ETCD:
		port = 2379
	case URL_SCHEME_KAFKA:
		port = 9092
	case URL_SCHEME_MSSQL:
		port = 1433
	case URL_SCHEME_ZOOKEEPER:
		port = 2181
	default:
		port = 80
	}
	return
}
