package huawei

import (
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/push"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	HUAWEI_HTTP_METHOD_POST = "POST"
	HUAWEI_TOKEN_TYPE       = "Bearer" //默认token类型
)

var (
	/*
		https://oauth-login.cloud.huawei.com/oauth2/v2/token
		token请求	字段			必选M/可选O	类型		描述
					grant_type			M		string		固定填”client_credentials”
					client_id			M		int			创建应用时获得的App ID。
					client_secret		M		string		Appid的密码，在开发者联盟上查


		请求示例：
		POST /oauth2/v2/token HTTP/1.1
		Content-Type: application/x-www-form-urlencoded
		grant_type=client_credentials&client_id=12345&client_secret=bKaZ0VE3EYrXaXCdCe3d2k9few

		响应示例:
		{"access_token":"CFyJ7eTl8WIPi9603E7Ro2Icy+K0JYe2qVjS8uzwCPltlO0fC7mZ0gzZX9p8CCwAaiU17nyP+N8+ORRzjjk1EA==","expires_in":3600,"token_type":"Bearer"}
	*/

	HUAWEI_AUTH_V2_URL = "https://oauth-login.cloud.huawei.com/oauth2/v2/token"

	//完整URL "https://push-api.cloud.huawei.com/v1/[appid]/messages:send"
	HUAWEI_PUSH_API_URL = "https://push-api.cloud.huawei.com/v1"
)

const (
	HUAWEI_ERROR_CODE_SUCCESS          = "80000000" //成功
	HUAWEI_ERROR_CODE_OAUTH            = "80200001" //Oauth认证错误, 请求HTTP头中Authorization参数里面的Access Token鉴权失败，请检查。
	HUAWEI_ERROR_CODE_TOKEN_EXPIRE     = "80200003" //Oauth Token过期, 请求HTTP头中Authorization参数里面的Access Token已过期，请重新申请后重试。
	HUAWEI_ERROR_CODE_MESSAGE_BODY     = "80100003" //消息结构体错误, 消息结构体参数携带错误，请按照此文档中请求参数部分进行检查。
	HUAWEI_ERROR_CODE_OVERFLOW         = "80300008" //消息体大小超过系统设置的默认值(4K) 请求消息体超默认值，请减小消息体后重试。
	HUAWEI_ERROR_CODE_PERMISSION       = "80300002" //APP被禁止发送当前应用无权限发送该消息，请检查应用权限。
	HUAWEI_ERROR_CODE_INVALID_TOKEN    = "80300007" //无效的Token消息请求携带的Push Token无效，请检查。
	HUAWEI_ERROR_CODE_INTERNAL_SERVICE = "500"      //内部服务器错误（与华为无关）
)

var errorCodeComments = map[string]string{

	HUAWEI_ERROR_CODE_SUCCESS:          "成功",
	HUAWEI_ERROR_CODE_OAUTH:            "Oauth认证错误, 请求HTTP头中Authorization参数里面的Access Token鉴权失败，请检查。",
	HUAWEI_ERROR_CODE_TOKEN_EXPIRE:     "Oauth Token过期, 请求HTTP头中Authorization参数里面的Access Token已过期，请重新申请后重试。",
	HUAWEI_ERROR_CODE_MESSAGE_BODY:     "消息结构体错误, 消息结构体参数携带错误，请按照此文档中请求参数部分进行检查。",
	HUAWEI_ERROR_CODE_OVERFLOW:         "消息体大小超过系统设置的默认值(4K) 请求消息体超默认值，请减小消息体后重试。",
	HUAWEI_ERROR_CODE_PERMISSION:       "APP被禁止发送当前应用无权限发送该消息，请检查应用权限。",
	HUAWEI_ERROR_CODE_INVALID_TOKEN:    "无效的Token消息请求携带的Push Token无效，请检查。",
	HUAWEI_ERROR_CODE_INTERNAL_SERVICE: "内部服务器错误（与华为无关）",
}

/* 典型的通知消息示例
{
    "validate_only":false, //为true时仅验证消息合法，不发送给APP
    "message":{
        "notification":{
            "title":"Big News",
            "body":"This is a Big News!"
        },
        "android":{
			"bi_tag":"JFXnBlNq53bwlxVWSHeDPxs5bWW", -- 用户定义回执ID
            "notification":{
                "title":"this is notification title",
                "body":"this is notification message body",
                "click_action":{
                    "type":1,
                    "intent":"#Intent;compo=com.rvr/.Activity;S.W=U;end"
                }
            }
        },
        "token":[
            "AAWWHI94sgUR2RU5_P1ZptUiwLq7W8XWJO2LxaAPuXw4_HOJFXnBlN-q5_3bwlxVW_SHeDPx_s5bWW-9DjtWZsvcm9CwXe1FHJg0u-D2pcQPcb3sTxDTJeiwEb9WBPl_9w"
        ]
    }
}
*/
type clickAction struct {
	Type   int    `json:"type"`
	Intent string `json:"intent"`
}

type badgeNotification struct {
	Num   int    `json:"num"`
	Class string `json:"class"`
}
type notification struct {
	Title       string      `json:"title"`        //标题
	Content     string      `json:"body"`         //内容
	ClickAction clickAction `json:"click_action"` //点击动作
}

type androidNotification struct {
	BiTag        string            `json:"bi_tag"`     //消息回执标识ID(用户自定义)
	ChannelId    string            `json:"channel_id"` //Android O版本中提供的通知消息渠道号
	Notification notification      `json:"notification"`
	Badge        badgeNotification `json:"badge"` //角标
}

type huaweiNotification struct {
	Title   string `json:"title"` //标题
	Content string `json:"body"`  //内容
}

type huaweiMessage struct {
	Token               []string            `json:"token"` //要发送的设备token
	HuaweiNotification  huaweiNotification  `json:"notification"`
	AndroidNotification androidNotification `json:"android"`
}

type pushMessage struct {
	ValidateOnly  bool          `json:"validate_only"` //为true时仅验证消息合法，不发送给APP
	HuaweiMessage huaweiMessage `json:"message"`
}

type Audience struct {
	RegId []string //接收者的设备注册id/token
}

type Message struct {
	Audience  Audience    //消息接收者
	Title     string      //标题
	Content   string      //内容
	Extra     interface{} //自定义数据结构
	ReceiptID string      //消息回执ID
	Badge     int         //角标数量
}

//获取token正确返回示例：
// HTTP/1.1 200 OK
// Content-Type: text/html;charset=UTF-8
// {"access_token":"CFyJ7eTl8WIPi9603E7Ro9Icy+K0JYe9qVjS8uzw3PltlO0fC7mZ0gzZX9p8CCwAaiU17nyP+N8+ORRzjjk1EA==","expires_in":3600,"token_type":"Bearer"}
type huaweiRespTokenOk struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

//错误返回示例：
// HTTP/1.1 400 Bad Request
// Content-Type: text/html;charset=UTF-8
// {"error":1203,"sub_error":12303,"error_description":"The request is invalid"}
type huaweiRespTokenError struct {
	MainError int    `json:"error"`
	SubError  int    `json:"sub_error"`
	Message   string `json:"error_description"`
}

//正常推送返回消息体
type huaweiResponse struct {
	Code  string `json:"code"`      //错误码
	Msg   string `json:"msg"`       //错误消息描述
	ReqId string `json:"requestId"` //推送成功后服务器返回的消息ID
}

type HuaWei struct {
	strAppId     string //APPID
	strAppSecret string //APP密钥
	strActivity  string //安卓跳转actitivy
	strToken     string //访问token
	strTokenType string //token类型: 例如Bearer(构造Authentication时使用)
	expireTime   int    //访问token过期时间
	strChannelId string //推送通道ID（跟APP设置保持一致即可）
	httpCli      *http.Client
}

func init() {
	if err := push.Register(push.AdatperType_Huawei, newHUAWEI); err != nil {
		log.Error("register %v instance error [%v]", push.AdatperType_Umeng, err.Error())
		panic("register instance failed")
	}
}

//创建推送接口对象
//args[0] => app id   		string 	[必填] 华为开发者 app id
//args[1] => app secret   	string  [必填] 华为开发者 app secret
//args[2] -> channel id 	string  [必填] 渠道ID
//args[3] -> activity 		string  [选填] 跳转activity
func newHUAWEI(args ...interface{}) push.IPush {
	var argc = len(args)
	var strAppId, strAppSecret, strChannelId, strActivity string
	if argc < 3 {
		panic(fmt.Sprintf("need 3 or 4 arguments, got %v", argc))
	}
	strAppId = args[0].(string)
	strAppSecret = args[1].(string)
	strChannelId = args[2].(string)
	if argc == 4 {
		strActivity = args[3].(string)
	}
	return &HuaWei{
		strAppId:     strAppId,
		strAppSecret: strAppSecret,
		strActivity:  strActivity,
		strTokenType: HUAWEI_TOKEN_TYPE,
		strChannelId: strChannelId,
		httpCli:      &http.Client{},
	}
}

//获取错误描述内容
func getErrorCodeComment(nErrorCode string) (strComment string) {

	var ok bool
	if strComment, ok = errorCodeComments[nErrorCode]; !ok {
		strComment = "<unknown error>"
		return
	}
	return strComment
}

//获取时间戳(秒)
func getTimestamp() int64 {

	return time.Now().UnixNano() / 1e9
}

func (h *HuaWei) getIntent(extras interface{}) (strIntentUri string) {

	var kvs map[string]string
	data, err := json.Marshal(extras)
	if err != nil {
		log.Errorf("%v", err.Error())
		return
	}

	log.Debugf("huawei getIntent -> extras json [%s]", data)

	if err = json.Unmarshal(data, &kvs); err != nil {
		log.Errorf("%v", err.Error())
		return
	}

	//intent自定义参数示例："longchat://com.legocity.longchat/notify_detail?type=1&chat_id=3"
	strIntentUri = fmt.Sprintf("longchat://%v/notify_detail?", h.strChannelId)
	strAnd := ""
	for k, v := range kvs {

		strIntentUri += fmt.Sprintf("%v%s=%d", strAnd, k, v)
		strAnd = "&"
	}
	return
}

//生成认证信息: "Bearer [access_token]"
func (h *HuaWei) getAuthentication() string {
	log.Debugf("token_type: [%v] token: [%v]", h.strTokenType, h.strToken)
	return fmt.Sprintf("%v %v", h.strTokenType, h.strToken)
}

func (h *HuaWei) getPushApiUrl() string {

	//完整URL "https://push-api.cloud.huawei.com/v1/[appid]/messages:send"
	return fmt.Sprintf("%v/%v/messages:send", HUAWEI_PUSH_API_URL, h.strAppId)
}

//向厂商获取通信token
//接口功能 开发者身份鉴权
//请求方法 POST
//请求编码 UTF-8
//Content-Type application/x-www-form-urlencoded
//参数：grant_type=client_credentials&client_id=12345&client_secret=bKaZ0VE3EYrXaXCdCe3d2k9few
//请求路径 https://oauth-login.cloud.huawei.com/oauth2/v2/token
func (h *HuaWei) getAuthToken() (err error) {

	kvs := url.Values{}
	kvs.Add("grant_type", "client_credentials")
	kvs.Add("client_id", h.strAppId)
	kvs.Add("client_secret", h.strAppSecret)

	var req *http.Request
	var httpResp *http.Response

	if req, err = http.NewRequest(HUAWEI_HTTP_METHOD_POST, HUAWEI_AUTH_V2_URL, strings.NewReader(kvs.Encode())); err != nil {
		log.Errorf("%v", err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if httpResp, err = h.httpCli.Do(req); err != nil {
		log.Errorf("[%v] to [%v] failed [%v]", HUAWEI_HTTP_METHOD_POST, HUAWEI_AUTH_V2_URL, err.Error())
		return
	}

	var respData []byte
	if respData, err = ioutil.ReadAll(httpResp.Body); err != nil {
		log.Errorf("ioutil.ReadAll from http response body failed [%v]", err.Error())
		return
	}

	log.Debugf("http response [%v]", string(respData))

	if httpResp.StatusCode != 200 { //返回错误

		var resp huaweiRespTokenError
		if err = json.Unmarshal(respData, &resp); err != nil {
			log.Errorf("json.Unmarshal from http response body failed [%v]", err.Error())
			return
		}
		err = fmt.Errorf("get huawei token error,  main code [%v] sub code [%v] message [%v]", resp.MainError, resp.SubError, resp.Message)
		log.Errorf("%v", err.Error())
	} else {
		var resp huaweiRespTokenOk
		if err = json.Unmarshal(respData, &resp); err != nil {
			log.Errorf("json.Unmarshal from http response body failed [%v]", err.Error())
			return
		}
		h.strTokenType = resp.TokenType
		h.strToken = resp.AccessToken
		h.expireTime = resp.ExpiresIn
	}

	return
}

//enable or disable debug output
func (h *HuaWei) Debug(enable bool) {
	if enable {
		log.SetLevel(0)
	} else {
		log.SetLevel(1)
	}
}

func (h *HuaWei) PushNotification(msg *push.Notification) (MsgID string, err error) {

	if len(msg.Audiences) == 0 || msg.Title == "" || msg.Alert == "" {
		err = fmt.Errorf("device token or msg title/content is nil [%+v]", msg)
		log.Errorf("%v", err.Error())
		return
	}

	var huaweiMsg = &pushMessage{
		ValidateOnly: false, //true为验证模式（不推送给手机APP端）
		HuaweiMessage: huaweiMessage{

			Token: msg.Audiences,
			HuaweiNotification: huaweiNotification{
				Title:   msg.Title,
				Content: msg.Alert,
			},
			AndroidNotification: androidNotification{
				BiTag:     "",             // 消息回执ID
				ChannelId: h.strChannelId, //推送渠道ID
				Badge: badgeNotification{
					Num:   msg.Badge,
					Class: "",
				},
				Notification: notification{
					Title:   msg.Title,
					Content: msg.Alert,
					ClickAction: clickAction{
						Type:   1, //intent activity
						Intent: h.getIntent(msg.Extra),
					},
				},
			},
		},
	}
	if h.strToken == "" {
		//获取token
		if err = h.getAuthToken(); err != nil {
			log.Errorf("%v", err.Error())
			return
		}
	}
	var Code string
	if MsgID, Code = h.sendMessage(huaweiMsg); Code != HUAWEI_ERROR_CODE_SUCCESS {

		if Code == HUAWEI_ERROR_CODE_TOKEN_EXPIRE { //token已失效，重新获取
			//获取token
			if err = h.getAuthToken(); err != nil {
				log.Errorf("%v", err.Error())
				return
			}
			//重发一次
			if MsgID, Code = h.sendMessage(huaweiMsg); Code != HUAWEI_ERROR_CODE_SUCCESS {
				err = fmt.Errorf("%v", getErrorCodeComment(Code))
				log.Errorf("%v", err.Error())
				return
			}
		} else {
			err = fmt.Errorf("%v", getErrorCodeComment(Code))
			log.Errorf("%v", err.Error())
			return
		}
	}

	return
}

func (h *HuaWei) sendMessage(msg *pushMessage) (MsgID string, Code string) {

	var err error
	var data []byte

	Code = HUAWEI_ERROR_CODE_INTERNAL_SERVICE
	data, err = json.Marshal(msg)
	if err != nil {
		log.Errorf("%v", err.Error())
		return
	}

	log.Debugf("sendMessage: message=%s", data)
	var req *http.Request
	if req, err = http.NewRequest(HUAWEI_HTTP_METHOD_POST, h.getPushApiUrl(), strings.NewReader(string(data))); err != nil {
		log.Errorf("http.NewRequest return error [%v]", err.Error())
		return
	}

	var httpResp *http.Response

	req.Header.Set("Authorization", h.getAuthentication()) //鉴权token
	req.Header.Set("Content-Type", "application/json")     //内容类型(JSON)

	if httpResp, err = h.httpCli.Do(req); err != nil {
		log.Errorf("[%v] to [%v] failed [%v]", HUAWEI_HTTP_METHOD_POST, h.getPushApiUrl(), err.Error())
		return
	}

	var respData []byte
	if respData, err = ioutil.ReadAll(httpResp.Body); err != nil {
		log.Errorf("ioutil.ReadAll from http response body failed [%v]", err.Error())
		return
	}

	defer httpResp.Body.Close()
	log.Debugf("http response [%s]", respData)

	var resp huaweiResponse
	if err = json.Unmarshal(respData, &resp); err != nil {
		log.Errorf("json.Unmarshal from http response body failed [%v]", err.Error())
		return
	}

	if resp.Code != HUAWEI_ERROR_CODE_SUCCESS {

		if resp.Code == HUAWEI_ERROR_CODE_TOKEN_EXPIRE {

			if err = h.getAuthToken(); err != nil {
				log.Errorf("%v", err.Error())
				return
			}

		} else {
			err = fmt.Errorf("%v", getErrorCodeComment(resp.Code))
			log.Warnf("[%v] to [%v] failed, reason [%v] data [%+v]", HUAWEI_HTTP_METHOD_POST, h.getPushApiUrl(), getErrorCodeComment(resp.Code), string(data))
			return "", resp.Code //华为报错
		}
	}

	return resp.ReqId, resp.Code
}
