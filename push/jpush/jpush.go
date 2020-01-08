package jpush

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/push"
	"io/ioutil"
	"net/http"
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
	appkey   string       //极光appkey
	secret   string       //极光secret
	is_prod  bool         //是否正式环境
	http_cli *http.Client //Http客户端对象
}

func init() {

	if err := push.Register(push.AdapterType_JPush, New); err != nil {
		log.Error("register %v instance error [%v]", push.AdapterType_JPush, err.Error())
		panic("register instance failed")
	}
}

//创建极光推送接口对象
//args[0] => appkey 极光访问账号(string)
//args[1] => secret 极光账号密码(string)
//args[2] => is_prod 是否正式环境(bool)
//New(appkey, secret, is_prod)
func New(args ...interface{}) push.IPush {

	if len(args) != JPUSH_PARAMS_COUNT {

		panic(fmt.Errorf("expect %v parameters, got %v", JPUSH_PARAMS_COUNT, len(args))) //参数个数错误
	}

	return &JPush{
		appkey:   args[0].(string),
		secret:   args[1].(string),
		is_prod:  args[2].(bool),
		http_cli: &http.Client{},
	}
}

//APP消息推送: 推送到极光服务器
//platforms 指定平台，空切片内部自动转为所有平台
func (j *JPush) Push(msg *push.Message) (MsgID string, err error) {

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
	content.Options.ApnsProduction = j.is_prod //判断IOS的生产还是测试环境
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

	log.Struct(&content)
	data, _ := json.Marshal(content)
	strPushContent := string(data)

	var resp []byte
	if resp, err = j.sendRequestWithAuthorization(push.HTTP_METHOD_POST, JPUSH_PUSHAPI_URL, strPushContent); err != nil {
		log.Error("%v", err.Error())
		return
	}

	log.Debug("post to [%v] with [%v] got response message = [%s]", JPUSH_PUSHAPI_URL, strPushContent, string(resp))

	delstr := string(resp)
	index := strings.Index(delstr, "sendno")

	if index < 0 { //没找到·就是失败
		ret := &IncorrectResp{}
		err = json.Unmarshal(resp, &ret)
		if err != nil {
			log.Error("parse jpush response json data [%v] to IncorrectResp error [%v]", string(resp), err.Error())
			err = errors.New("parse jpush response json data to IncorrectResp failed")
			return
		} else {
			err = fmt.Errorf("%v", ret.Error.Message)
			log.Error("jpush response [%#v]", ret)
			return
		}
	}
	log.Debug("jpush [%+v] to [%v] ok", strPushContent, JPUSH_PUSHAPI_URL)
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
func (j *JPush) sendRequestWithAuthorization(strHttpMethod, strUrl string, message interface{}) (resp_data []byte, err error) {

	var data []byte

	if message != nil {
		switch message.(type) {
		case string:
			data = []byte(message.(string))
		default:
			data, err = json.Marshal(message)
			if err != nil {
				log.Error("message json marshal error [%v]", err.Error())
				return
			}
		}
	}

	var req *http.Request
	var resp *http.Response

	authorization := j.getBase64Authorization()
	if req, err = http.NewRequest(strHttpMethod, strUrl, bytes.NewBuffer(data)); err != nil {

		log.Error("http.NewRequest error [%v]", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)
	//log.Debugf("send to http url [%v] with data [%v] ready", strUrl, string(data))
	if resp, err = j.http_cli.Do(req); err != nil {
		log.Error("send http request error [%v]", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp_data, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("ioutil.ReadAll(resp.Body) error [%v]", err.Error())
		return
	}

	log.Debug("[%v] http url [%v] with data [%s] successful, got response [%s]", strHttpMethod, strUrl, data, resp_data)
	return
}

//按照极光认证要求将appkey和secret做base64编码
func (j *JPush) getBase64Authorization() (strAuthorization string) {

	strEncode := base64.StdEncoding.EncodeToString(([]byte)(fmt.Sprintf("%v:%v", j.appkey, j.secret)))
	strAuthorization = fmt.Sprintf("Basic %v", strEncode)
	log.Debug("authorization [%v]", strAuthorization)
	return
}
