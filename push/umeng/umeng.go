package umeng

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/push"
	"io/ioutil"
	"github.com/civet148/gotools/log"
	"net/http"
	"strings"
	"time"
)

/*
* 友盟推送(支持华为/小米/vivo/oppo/魅族厂商推送通道)
 */

var UMENG_PARAMS_COUNT = 2
var UMENG_METHOD_POST = "POST"

//http://msg.umeng.com/api/send?sign=mysign
//https://msgapi.umeng.com/api/send?sign=mysign
var UMENG_PUSH_API_URL = "https://msgapi.umeng.com/api/send"

var (
	UMENG_RET_SUCCESS = "SUCCESS" //成功返回
	UMENG_RET_FAIL    = "FAIL"    //失败返回
)

var (
	DISPLAY_TYPE_MESSAGE      = "message"
	DISPLAY_TYPE_NOTIFICATION = "notification"
)

var (
	UMENG_PUSH_TYPE_UNICAST        = "unicast"        //单播, 针对某个设备推送
	UMENG_PUSH_TYPE_LISTCAST       = "listcast"       //列播，要求不超过500个device_token
	UMENG_PUSH_TYPE_FILECAST       = "filecast"       //文件播，多个device_token可通过文件形式批量发送
	UMENG_PUSH_TYPE_BROADCAST      = "broadcast"      //广播,针对所有设备推送
	UMENG_PUSH_TYPE_GROUPCAST      = "groupcast"      //组播，按照filter筛选用户群, 请参照filter参数
	UMENG_PUSH_TYPE_CUSTOMIZEDCAST = "customizedcast" //通过alias进行推送，包括以下两种case: - alias: 对单个或者多个alias进行推送  - file_id: 将alias存放到文件后，根据file_id来推送
)

type Audience struct {
	RegId []string //接收者的设备注册id/token(单播/列播 MAX=500)
	Tag   []string //按照filter筛选用户群(组播)
	Alias []string //通过alias进行推送(需客户端SDK设置别名)
}

type Message struct {
	Audience Audience    //消息接收者
	Title    string      //标题
	Content  string      //内容
	Extra    interface{} //自定义数据结构
}

type Umeng struct {
	appKey    string       //AppKey
	appSecret string       //AppSecret
	httpCli  *http.Client //http client
}

type umengBody struct {
	Ticker      string      `json:"ticker"`                 //必填，通知栏提示文字
	Title       string      `json:"title"`                  //必填，通知标题
	Text        string      `json:"text"`                   //必填，通知文字描述
	Icon        string      `json:"icon,omitempty"`         //可选，状态栏图标ID，R.drawable.[smallIcon]，如果没有，默认使用应用图标
	LargeIcon   string      `json:"largeIcon,omitempty"`    //可选，通知栏拉开后左侧图标ID
	Image       string      `json:"img,omitempty"`          //可选，通知栏大图标的URL链接
	Sound       string      `json:"sound,omitempty"`        //可选，通知声音文件，文件不存在，则使用系统默认Notification提示音
	BuilderID   int         `json:"builder_id"`             //可选，默认为0，用于标识该通知采用的样式(使用该参数时，开发者必须在SDK里面实现自定义通知栏样式)
	PlayVibrate bool        `json:"play_vibrate,omitempty"` //可选，收到通知是否震动，默认为"true"
	PlayLights  bool        `json:"play_lights,omitempty"`  //可选，收到通知是否闪灯，默认为"true"
	PlaySound   bool        `json:"play_sound,omitempty"`   //可选，收到通知是否发出声音，默认为"true"
	AfterOpen   string      `json:"after_open,omitempty"`   //可选，点击"通知"的后续行为, 默认为"go_app":打开应用; "go_url":跳转到URL; "go_activity":打开特定的activity; "go_custom": 用户自定义内容
	Url         string      `json:"url,omitempty"`          //当after_open="go_url"时，必填
	Activity    string      `json:"activity,omitempty"`     //当after_open=go_activity时，必填。
	Custom      interface{} `json:"custom,omitempty"`       //当display_type=message时, 必填(用户自定义内容，可以为字符串或者JSON格式)
}

type umengPayload struct { //only for android device

	//消息类型[必填]: notification(通知)、message(消息)
	DisplayType string `json:"display_type"`

	//数据body
	//当display_type=message时，body的内容只需填写custom字段。
	Body umengBody `json:"body"`

	Extra interface{} `json:"extra,omitempty"` //用户自定义内容(JSON格式)
}

type umengNotification struct {
	//应用唯一标识[必填]
	AppKey string `json:"appkey"`

	//时间戳[必填]，10位或者13位均可，时间戳有效期为10分钟(时区+8)
	Timestamp string `json:"timestamp"`

	//消息发送类型[必填], 其值可以为:
	//unicast-单播, 针对某个设备推送
	//listcast-列播，要求不超过500个device_token
	//filecast-文件播，多个device_token可通过文件形式批量发送
	//broadcast-广播,针对所有设备推送
	//groupcast-组播，按照filter筛选用户群, 请参照filter参数
	//customizedcast，通过alias进行推送，包括以下两种case:
	//     - alias: 对单个或者多个alias进行推送
	//     - file_id: 将alias存放到文件后，根据file_id来推送
	Type string `json:"type"`

	//设备tokens
	//当type=unicast时, 必填, 表示指定的单个设备
	//当type=listcast时, 必填, 要求不超过500个, 以英文逗号分隔
	DeviceTokens string `json:"device_tokens,omitempty"`

	//别名类型(某些情况下必填)
	//当type=customizedcast时, 必填
	//alias的类型, alias_type可由开发者自定义, 开发者在SDK中
	//调用setAlias(alias, alias_type)时所设置的alias_type
	AliasType string `json:"alias_type,omitempty"`

	//别名(某些情况下必填)
	//当type=customizedcast时, 选填(此参数和file_id二选一)
	//开发者填写自己的alias, 要求不超过500个alias, 多个alias以英文逗号间隔
	//在SDK中调用setAlias(alias, alias_type)时所设置的alias
	AliasName string `json:"alias,omitempty"`

	//文件id(某些情况下必填)
	//当type=filecast时，必填，file内容为多条device_token，以回车符分隔
	//当type=customizedcast时，选填(此参数和alias二选一)
	//file内容为多条alias，以回车符分隔。注意同一个文件内的alias所对应
	//的alias_type必须和接口参数alias_type一致。使用文件播需要先调用文件上传接口获取file_id，参照"文件上传"
	FileID string `json:"file_id,omitempty"`

	//过滤(某些情况下必填)
	//当type=groupcast时，必填，用户筛选条件，如用户标签、渠道等，参考附录G。
	//filter的内容长度最大为3000B）
	Filter string `json:"filter,omitempty"`

	Payload umengPayload `json:"payload"`
}

type Response struct {
	Ret  string `json:"ret"` //成功返回"SUCCESS" 失败返回"FAIL"
	Data struct {
		MsgID     string `json:"msg_id"`     //单播类消息(type为unicast、listcast、customizedcast且不带file_id)返回：
		TaskID    string `json:"task_id"`    //任务类消息(type为broadcast、groupcast、filecast、customizedcast且file_id不为空)返回
		ErrorCode string `json:"error_code"` //当ret返回值为FAIL时指示错误码
		ErrorMsg  string `json:"error_msg"`  //当ret返回值为FAIL时指示错误细信息
	} `json:"data"`
}


func init() {

	if err := push.Register(push.AdatperType_Umeng, New); err != nil {
		log.Error("register %v instance error [%v]", push.AdatperType_Umeng, err.Error())
		panic("register instance failed")
	}
}

//创建友盟推送接口对象
//args[0] => appkey string 友盟App key
//args[1] => secret string 友盟App master secret
func New(args ...interface{}) push.IPush {

	if len(args) != UMENG_PARAMS_COUNT {
		panic(fmt.Errorf("expect %v parameters, got %v", UMENG_PARAMS_COUNT, len(args))) //参数个数错误
	}

	return &Umeng{
		appKey:    args[0].(string),
		appSecret: args[1].(string),
		httpCli:  &http.Client{},
	}
}

//获取签名
//签名
//为了确保用户发送的请求不被更改，我们设计了签名算法。该算法基本可以保证请求是合法者发送且参数没有被修改，但无法保证不被偷窥。 签名生成规则：
// - 提取请求方法method（POST，全大写）
// - 提取请求url信息，包括Host字段的域名(或ip:端口)和URI的path部分。注意不包括path的querystring。比如http://msg.umeng.com/api/send 或者 http://msg.umeng.com/api/status
// - 提取请求的post-body json字符串
// - 拼接请求方法、url、post-body及应用的app_master_secret
// - 将D形成字符串计算MD5值，形成一个32位的十六进制（字母小写）字符串，即为本次请求sign（签名）的值；Sign=MD5($http_method$url$post-body$app_master_secret)
//python代码参考 https://developer.umeng.com/docs/66632/detail/68343#h2--k-17
func (u *Umeng) getSignature(strMethod, strUrl, body string) (strSign string) {

	strMethod = strings.ToUpper(strMethod) //http method转为大写
	strText := fmt.Sprintf("%v%v%v%v", strMethod, strUrl, body, u.appSecret)
	log.Debug("getSignature() text-> [%v]", strText)
	m := md5.New()
	m.Write([]byte(strText))
	strSign = hex.EncodeToString(m.Sum(nil))
	log.Debug("getSignature() sign -> [%v]", strSign)
	return
}

func (u *Umeng) sendPushRequest(strMethod, strUrl string, reqBody interface{}) (MsgID string, err error) {

	var body []byte

	if body, err = json.Marshal(reqBody); err != nil {
		log.Error("request body marshal error [%v]", err.Error())
		return
	}

	strSign := u.getSignature(strMethod, strUrl, string(body))

	var req *http.Request

	strUrl = fmt.Sprintf("%v?sign=%v", strUrl, strSign)
	if req, err = http.NewRequest(strMethod, strUrl, bytes.NewBuffer(body)); err != nil {
		log.Error("sendHttpRequest http.NewRequest return error [%v]", err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")

	var respHttp *http.Response
	log.Debug("sendHttpRequest send push to http url [%v] body [%v]", strUrl, string(body))
	if respHttp, err = u.httpCli.Do(req); err != nil {
		log.Error("sendHttpRequest send push notification return error [%v], http url [%v] body [%v]", err.Error(), strUrl, string(body))
		return
	}

	defer respHttp.Body.Close()

	var respData []byte
	if respData, err = ioutil.ReadAll(respHttp.Body); err != nil {
		log.Error("sendHttpRequest read response data error [%v]", err.Error())
		return
	}

	log.Info("sendHttpRequest response raw data [%v]", string(respData))

	var respUmeng Response
	if err = json.Unmarshal(respData, &respUmeng); err != nil {
		log.Error("sendHttpRequest unmarshal response to struct error [%v] resp data [%v]", err.Error(), string(respData))
		return
	}

	if respUmeng.Ret == UMENG_RET_SUCCESS {
		log.Debug("sendHttpRequest response ok [%+v]", respUmeng)
		if respUmeng.Data.MsgID != "" {
			MsgID = respUmeng.Data.MsgID
		} else if respUmeng.Data.TaskID != "" {
			MsgID = respUmeng.Data.TaskID
		}
	} else {

		err = fmt.Errorf("%v", respUmeng.Data.ErrorMsg)
		log.Warn("sendHttpRequest response resp [%+v]", respUmeng)
	}

	return
}

//推送消息(当前只支持Android单播和列播)
func (u *Umeng) Push(msg *push.Message) (MsgID string, err error) {

	if msg.AudienceType != push.AUDIENCE_TYPE_REGID_TOKEN {

		err = fmt.Errorf("Umeng just can use AUDIENCE_TYPE_REGID_TOKEN to push message")
		log.Error("%v", err.Error())
		return
	}

	var targets = len(msg.Audiences)
	var strType = UMENG_PUSH_TYPE_UNICAST
	if targets == 0 {
		log.Error("no any target found at msg [%+v]", msg)
		return
	}
	if targets > 1 {
		strType = UMENG_PUSH_TYPE_LISTCAST
	}

	ts := time.Now().Unix()

	tokens := strings.Join(msg.Audiences, ",")

	notification := &umengNotification{
		AppKey:       u.appKey,
		Timestamp:    fmt.Sprintf("%v", ts),
		Type:         strType,
		DeviceTokens: tokens,
		Payload: umengPayload{
			DisplayType: DISPLAY_TYPE_NOTIFICATION,
			Body: umengBody{
				Ticker:      msg.Title,
				Title:       msg.Title,
				Text:        msg.Alert,
				//Icon:        "",
				//LargeIcon:   "",
				//Image:       "",
				//Sound:       "",
				BuilderID:   0,
				//PlayVibrate: false,
				//PlayLights:  false,
				//PlaySound:   false,
				//AfterOpen:   "",
				//Url:         "",
				//Activity:    "",
				//Custom:      nil,
			},
			Extra: msg.Extra,
		},
	}

	return u.sendPushRequest(UMENG_METHOD_POST, UMENG_PUSH_API_URL, notification)
}

//enable or disable debug output
func (u *Umeng) Debug(enable bool) {
	if enable {
		log.SetLevel(0)
	} else {
		log.SetLevel(1)
	}
}