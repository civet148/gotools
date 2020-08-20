package jpush

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/civet148/gotools/comm/httpx"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/push"
	"strings"
)

/*
* 极光消息推送v3接口
 */

const (
	JPUSH_PARAMS_COUNT = 3 //两个参数（appkey+secret+is_prod）
)

var JPUSH_PUSHAPI_URL = "https://api.jpush.cn/v3/push"
var JPUSH_DEVICES_URL = "https://device.jpush.cn/v3/devices"
var JPUSH_ALIASES_URL = "https://device.jpush.cn/v3/aliases"

type JPush struct {
	appKey    string //极光appkey
	appSecret string //极光secret
	isProd    bool   //是否正式环境
}

func init() {

	if err := push.Register(push.AdapterType_JPush, newJPUSH); err != nil {
		log.Error("register %v instance error [%v]", push.AdapterType_JPush, err.Error())
		panic("register instance failed")
	}
}

//创建极光推送接口对象
//args[0] => appkey 极光访问账号(string)
//args[1] => secret 极光账号密码(string)
//args[2] => is_prod 是否正式环境(bool)
//New(appkey, secret, is_prod)
func newJPUSH(args ...interface{}) push.IPush {

	if len(args) != JPUSH_PARAMS_COUNT {

		panic(fmt.Errorf("expect %v parameters, got %v", JPUSH_PARAMS_COUNT, len(args))) //参数个数错误
	}

	return &JPush{
		appKey:    args[0].(string),
		appSecret: args[1].(string),
		isProd:    args[2].(bool),
	}
}

//APP消息推送: 推送到极光服务器
//platforms 指定平台，空切片内部自动转为所有平台
func (j *JPush) PushNotification(msg *push.Notification) (MsgID string, err error) {

	if len(msg.Platforms) == 0 { //不指定平台则为三个平台同时发
		msg.Platforms = []string{push.PLATFORM_ANDROID, push.PLATFORM_IOS, push.PLATFORM_WINPHONE}
	}
	content := Content{
		Platform: msg.Platforms,
	}

	//Notification for android
	content.Notification.Android.Title = msg.Title //标题
	content.Notification.Android.Alert = msg.Alert //内容
	content.Notification.Android.BuilderId = 0     //设置安卓声音
	if content.Notification.Android.BuilderId == 0 {
		content.Notification.Android.BuilderId = 2
	}
	content.Notification.Android.Extras = msg.Extra //自定义内容

	//Notification for iOS
	content.Notification.IOS.Alert = msg.Alert  //必须要有值才能收到
	content.Notification.IOS.Extras = msg.Extra //自定义内容
	content.Options.TimeToLive = 60
	content.Options.ApnsProduction = j.isProd  //判断IOS的生产还是测试环境
	content.Notification.IOS.Sound = "default" /*msg.SoundIOS*/ //默认有声音
	content.Notification.IOS.Badge = "+1"      //角标默+1

	switch msg.AudienceType {
	case push.AUDIENCE_TYPE_REGID_TOKEN: //device token or register id
		content.Audience.RegId = msg.Audiences
	case push.AUDIENCE_TYPE_TAG: //message group
		content.Audience.Tag = msg.Audiences
		//case push.AUDIENCE_TYPE_ALIAS: //disabled alias of jpush
		//	content.Audience.Alias = msg.Audiences
	}

	var resp *httpx.Response
	if resp, err = j.sendRequestWithAuthorization(JPUSH_PUSHAPI_URL, content); err != nil {
		log.Error("%v", err.Error())
		return
	}

	log.Debug("post to [%v] with [%+v] got response code [%v] message = [%s]", JPUSH_PUSHAPI_URL, content, resp.StatusCode, resp.Body)

	index := strings.Index(resp.Body, "sendno")

	if index < 0 { //没找到·就是失败
		ret := &IncorrectResp{}
		err = json.Unmarshal([]byte(resp.Body), &ret)
		if err != nil {
			log.Error("parse jpush response json data [%v] to IncorrectResp error [%v]", resp.Body, err.Error())
			err = errors.New("parse jpush response json data to IncorrectResp failed")
			return
		} else {
			err = fmt.Errorf("%v", ret.Error.Message)
			log.Error("jpush response [%#v]", ret)
			return
		}
	}
	log.Debug("jpush [%+v] to [%v] ok", content, JPUSH_PUSHAPI_URL)
	return
}

//开启关闭调试日志
func (j *JPush) Debug(enable bool) {
	if enable {
		log.SetLevel(0)
	} else {
		log.SetLevel(1)
	}
}

//获取HTTP客户端对象（包含认证信息）
func (j *JPush) sendRequestWithAuthorization(strUrl string, message interface{}) (response *httpx.Response, err error) {

	authorization := j.getBase64Authorization()
	c := httpx.NewHttpClient(3)
	c.Header().SetApplicationJson().SetAuthorization(authorization)
	if response, err = c.Post(strUrl, message); err != nil {
		log.Error("%v", err.Error())
		return
	}
	log.Debug("response [%+v]", response)
	return
}

//按照极光认证要求将appkey和secret做base64编码
func (j *JPush) getBase64Authorization() (strAuthorization string) {

	strEncode := base64.StdEncoding.EncodeToString(([]byte)(fmt.Sprintf("%v:%v", j.appKey, j.appSecret)))
	strAuthorization = fmt.Sprintf("Basic %v", strEncode)
	log.Debug("authorization [%v]", strAuthorization)
	return
}
