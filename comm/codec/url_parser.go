package codec

import (
	"net/url"
	"strings"
)

type UrlInfo struct {
	Scheme     string
	Host       string // host name and port like '127.0.0.1:3306'
	User       string
	Password   string
	Path       string
	Fragment   string
	Opaque     string
	ForceQuery bool
	Queries    map[string]string
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

	// scheme://[userinfo@]host:port/path[?query][#fragment]

	strUrl = strings.TrimSpace(strUrl)
	if strings.Contains(strUrl, "@") { // if a url have user+password, there must be have '@'
		// find first '://'
		var strScheme string
		_ = strScheme

		index := strings.LastIndex(strUrl, "://")
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
				strUrl = strScheme + "://"
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
