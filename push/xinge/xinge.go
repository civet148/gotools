package xinge

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/comm/httpx"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/push"
)

/*
* 信鸽消息推送(支持Google FCM/小米/华为/魅族厂商推送通道)
* 官方API文档链接：https://xg.qq.com/docs/server_api/v3/push_api_v3.html
 */

const (
	XINGE_PARAMS_COUNT = 3 //必填参数个数
)

var (
	XINGE_AUDIENCE_ALL          = "all"
	XINGE_AUDIENCE_TAG          = "tag"
	XINGE_AUDIENCE_TOKEN        = "token"
	XINGE_AUDIENCE_TOKEN_LIST   = "token_list"
	XINGE_AUDIENCE_ACCOUNT      = "account"
	XINGE_AUDIENCE_ACCOUNT_LIST = "account_list"
)

var (
	XINGE_ENV_DEV     = "dev"     //开发环境
	XINGE_ENV_PRODUCT = "product" //生产环境
)

var (
	XINGE_PUSH_TYPE_NOTIFY  = "notify"  //通知
	XINGE_PUSH_TYPE_MESSAGE = "message" //透传消息
)

var (
	XINGE_PLATFORM_ANDROID = "android"
	XINGE_PLATFORM_IOS     = "ios"
)

var XINGE_PUSHAPI_URL = "https://openapi.xg.qq.com/v3/push/app"

type XinGe struct {
	appKey      string //信鸽appkey
	appSecret   string //信鸽secret
	isProd      bool   //是否正式环境[仅适用于苹果iOS设备]（true=正式环境 false=测试环境）
	strActivity string //安卓通知点击跳转activity
}

type xingeTime struct {
	Hour string `json:"hour"`
	Min  string `json:"min"`
}

type xingeAccetpTime struct {
	Start xingeTime `json:"start"`
	End   xingeTime `json:"end"`
}

type xingeAtyAttr struct {
	If int `json:"if"` // Intent的Flag属性
	Pf int `json:"pf"` // PendingIntent的Flag属性
}

type xingeBrowser struct {
	Url     string `json:"url"`     //仅支持http、https
	Confirm int    `json:"confirm"` // 是否需要用户确认(0=不需要 1=需要)
}

type xingeAction struct {
	ActionType int           `json:"action_type"` // 动作类型，1，打开activity或app本身；2，打开浏览器；3，打开Intent
	Activity   string        `json:"activity"`
	AtyAttr    *xingeAtyAttr `json:"aty_attr,omitempty"` //activity属性，只针对action_type=1的情况
	Browser    *xingeBrowser `json:"browser,omitempty"`  //URL跳转
	Intent     string        `json:"intent"`             //SDK版本需要大于等于3.2.3，然后在客户端的intent配置data标签，并设置scheme属性
}

type xingeTagList struct {
	Tags []string `json:"tags"`
	Op   string   `json:"op"`
}

type xingeAndroid struct {

	//通知消息对象的唯一标识（只对信鸽通道生效） 默认值 0 是否必填[否]
	//1）大于0：会覆盖先前相同id的消息
	//2）等于0：展示本条通知且不影响其他消息
	//3）等于-1：将清除先前所有消息，仅展示本条消息
	NotificationID int `json:"n_id"`
	//本地通知样式标识(默认0)
	BuilderID int `json:"builder_id"`
	//是否有铃声(0=无铃声 1=有铃声 默认1)
	Ring int `json:"ring"`
	//指定Android工程里raw目录中的铃声文件名，不需要后缀名
	RingRaw string `json:"ring_raw"`
	//是否使用震动(0=没有震动 1=有震动 默认1)
	Vibrate int `json:"vibrate"`
	//是否使用呼吸灯(0=不使用呼吸灯 1=使用呼吸灯 默认1)
	Lights int `json:"lights"`
	//通知栏是否可清除(0=不可清除 1=可清除 默认1)
	Clearable int `json:"clearable"`
	//通知栏图标是应用内图标还是上传图标(0=应用内图标 1=上传图标 默认0)
	IconType int `json:"icon_type"`
	//应用内图标文件名或者下载图标的url地址
	IconRes string `json:"icon_res"`
	//设置是否覆盖指定编号的通知样式(默认1)
	StyleID int `json:"style_id"`
	//消息在状态栏显示的图标，若不设置，则显示应用图标
	SmallIcon string `json:"small_icon"`
	//设置点击通知栏之后的行为，默认为打开app
	Action *xingeAction `json:"action,omitempty"`
	//用户自定义的键值对
	CustomContent interface{} `json:"custom_content,omitempty"`
}

type xingeMessage struct {
	//消息标题[必填]
	Title string `json:"title"`
	//消息内容[必填]
	Content string `json:"content"`
	//消息将在哪些时间段允许推送给用户，建议小于10个，不能为空
	AcceptTime []xingeAccetpTime `json:"accept_time"`
	//富媒体元素地址，建议小于5个 （仅限SDK4.2.0及以上版本使用）
	XgMediaResources string `json:"xg_media_resources"`
	//信鸽消息体[必填]
	Android *xingeAndroid `json:"android"` //Android消息体
}
type xingeNotification struct {

	//推送目标[必填]
	//1）all：全量推送
	//2）tag：标签推送
	//3）token：单设备推送
	//4）token_list：设备列表推送
	//5）account：单账号推送
	//6）account_list：账号列表推送
	AudienceType string `json:"audience_type"`

	//客户端平台类型[必填]
	//1）android：安卓
	//2）ios：苹果
	Platform string `json:"platform"`

	//设备token列表
	TokenList []string `json:"token_list,omitempty"`

	//标签列表
	TagList *xingeTagList `json:"tag_list,omitempty"`

	//消息类型[必填]
	//1）notify：通知
	//2）message：透传消息/静默消息
	MessageType string `json:"message_type"`

	//消息体
	Message xingeMessage `json:"message"`

	//用户指定推送环境[必填]
	//1）product： 推送生产环境
	//2）dev： 推送开发环境
	//仅限iOS平台推送使用(默认为product)
	Environment string `json:"environment"`
}

/*
{
    "seq": 0,
    "environment": "product",
    "ret_code": 0,
    "push_id": "3895624686"
}
*/
type xingeResponse struct {
	SeqNo       int    `json:"seq"`
	RetCode     int    `json:"ret_code"` //0表示成功，其他表示失败
	PushID      string `json:"push_id"`  //推送成功后返回的消息ID
	Environment string `json:"environment"`
	ErrorMsg    string `json:"err_msg"` //ret_code != 0时指示错误信息
}

func init() {

	if err := push.Register(push.AdapterType_XinGe, newXINGE); err != nil {
		log.Error("register %v instance error [%v]", push.AdapterType_XinGe, err.Error())
		panic("register instance failed")
	}
}

//创建信鸽推送接口对象
//args[0] => appKey 	string 信鸽app key
//args[1] => appSecret 	string 信鸽app secret
//args[2] => isProd   	bool   是否iOS正式环境（true=正式环境 false=测试环境）
//args[3] => activity   string [可选]推送通知点击跳转activity
func newXINGE(args ...interface{}) push.IPush {

	var nArgs = len(args)
	if nArgs < XINGE_PARAMS_COUNT {

		panic(fmt.Errorf("expect %v parameters, got %v", XINGE_PARAMS_COUNT, len(args))) //参数个数错误
	}

	var strActivity string
	if nArgs > XINGE_PARAMS_COUNT {
		strActivity = args[3].(string)
	}
	return &XinGe{
		appKey:      args[0].(string),
		appSecret:   args[1].(string),
		isProd:      args[2].(bool),
		strActivity: strActivity,
	}
}

//APP消息推送: 推送到信鸽服务器
//platforms 指定平台(目前友盟仅支持Android)
func (x *XinGe) PushNotification(msg *push.Notification) (MsgID string, err error) {

	var strEnv = XINGE_ENV_DEV
	if x.isProd {
		strEnv = XINGE_ENV_PRODUCT
	}

	targets := len(msg.Audiences)
	if targets == 0 {
		err = fmt.Errorf("no any target found")
		log.Error("%v", err.Error())
		return
	}

	var notification = xingeNotification{

		Platform:    XINGE_PLATFORM_ANDROID,
		MessageType: XINGE_PUSH_TYPE_NOTIFY,
		Message: xingeMessage{

			Title:   msg.Title,
			Content: msg.Alert,
			AcceptTime: []xingeAccetpTime{{
				Start: xingeTime{
					Hour: "00",
					Min:  "01",
				},
				End: xingeTime{
					Hour: "23",
					Min:  "59",
				},
			}},
			XgMediaResources: "",
			Android: &xingeAndroid{
				NotificationID: 0,
				BuilderID:      0,
				Ring:           0,
				RingRaw:        "",
				Vibrate:        1,
				Lights:         1,
				Clearable:      1,
				IconType:       0,
				IconRes:        "",
				StyleID:        0,
				SmallIcon:      "",
				Action: &xingeAction{
					ActionType: 1,
					Activity:   x.strActivity,
					AtyAttr:    &xingeAtyAttr{},
					Browser:    &xingeBrowser{},
					Intent:     "",
				},
				CustomContent: msg.Extra,
			},
		},
		Environment: strEnv,
	}

	switch msg.AudienceType {
	case push.AUDIENCE_TYPE_REGID_TOKEN:
		{
			notification.AudienceType = XINGE_AUDIENCE_TOKEN_LIST
			notification.TokenList = msg.Audiences
		}
	case push.AUDIENCE_TYPE_TAG:
		{
			notification.AudienceType = XINGE_AUDIENCE_TAG
			notification.TagList = &xingeTagList{Tags: msg.Audiences, Op: "AND"}
		}
	default:
		err = fmt.Errorf("can't handle push type [%v]", msg.AudienceType)
		log.Error("%v", err.Error())
		return
	}

	var response *httpx.Response
	if response, err = x.sendRequestWithAuthorization(XINGE_PUSHAPI_URL, notification); err != nil {
		log.Error("%v", err.Error())
		return
	}
	log.Debug("post to [%v] with [%+v] got response data = [%+v]", XINGE_PUSHAPI_URL, notification, response)

	if response.StatusCode == 200 {
		var resp xingeResponse

		if err = json.Unmarshal([]byte(response.Body), &resp); err != nil {
			log.Error("unmarshal http response data [%+v] to xingeResponse object error [%v]", response.Body, err.Error())
			err = fmt.Errorf("%s", response.Body)
			return
		}
		if resp.RetCode == 0 {
			MsgID = resp.PushID
			log.Debug("XINGE push to url [%+v] content [%+v] ok, MsgID [%v]", XINGE_PUSHAPI_URL, notification, resp.PushID)
		} else {

			err = fmt.Errorf("%+v", resp)
			log.Error("XINGE push to url [%+v] content [%+v] failed, response [%+v]", XINGE_PUSHAPI_URL, notification, resp)
			return
		}
	} else {
		err = fmt.Errorf("%+v", response.Body)
	}

	return
}

//开启关闭调试日志
func (x *XinGe) Debug(enable bool) {
	if enable {
		log.SetLevel(0)
	} else {
		log.SetLevel(1)
	}
}

//获取HTTP客户端对象（包含认证信息）
func (x *XinGe) sendRequestWithAuthorization(strUrl string, message interface{}) (response *httpx.Response, err error) {

	c := httpx.NewHttpClient(3)
	c.Header().SetApplicationJson().SetAuthorization(x.getBase64Authorization())

	if response, err = c.Post(strUrl, message); err != nil {
		log.Error("http post error [%v]", err.Error())
		return
	}

	log.Debug("http post url [%v] with message [%+v] successful, got response [%+v]", strUrl, message, response)

	return
}

//按照极光认证要求将appkey和secret做base64编码
func (x *XinGe) getBase64Authorization() (strAuthorization string) {

	strEncode := base64.StdEncoding.EncodeToString(([]byte)(fmt.Sprintf("%v:%v", x.appKey, x.appSecret)))
	strAuthorization = fmt.Sprintf("Basic %v", strEncode)
	log.Debug("authorization [%v]", strAuthorization)
	return
}
