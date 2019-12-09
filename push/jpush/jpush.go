package jpush

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/civet148/gotools/comm/httpc"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/push"
	"strings"
)

/*
* 极光消息推送v3接口
 */

const (
	JPUSH_PARAMS_COUNT = 2 //两个参数（appkey+secret+is_prod）
)

var JPUSH_PUSHAPI_URL = "https://api.jpush.cn/v3/push"
var JPUSH_DEVICES_URL = "https://device.jpush.cn/v3/devices"
var JPUSH_ALIASES_URL = "https://device.jpush.cn/v3/aliases"

type JPush struct {
	appkey string //极光appkey
	secret string //极光secret
}

func init() {

	if err := push.Register(push.AdapterType_JPush, New); err != nil {
		log.Error("register jpush instance error [%v]", err.Error())
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
		appkey: args[0].(string),
		secret: args[1].(string),
	}
}

//注册别名（后续大部分消息推送主要通过别名完成推送）
//AliasName 	别名，作为消息推送的唯一标识
//RegisterId 	App集成极光推送SDK后获得并由App客户端负责传递到服务器端
//Mobile 		用户手机号
//备注：			别名规则需根据实际情况定义好，确保一台设备一个账户唯一避免出现一台设备多个账户登录或一个账户多个设备登录造成消息推送失败
func (j *JPush) Register(strAliasName, strRegisterId, strMobile string) (err error) {

	reg := &Register{}
	reg.Mobile = strMobile
	reg.Alias = strAliasName
	httpCli := j.getHttpClientWithAuthorization()
	strRegUrl := fmt.Sprintf("%v/%v", JPUSH_DEVICES_URL, strRegisterId)
	log.Debug("register post to url [%v]\n", strRegUrl)
	data, _ := json.Marshal(reg)
	log.Debug("register data [%v]", string(data))

	resp, err := httpCli.SendUpstream(string(data), "POST", strRegUrl)
	if err != nil {
		log.Error("post to [%v] with register id [%v] mobile [%v] alias [%v] error [%v]\n", strRegUrl, strRegisterId, strMobile, strAliasName, err.Error())
		return err
	}

	log.Debug("jpush response [%s]\n", string(resp))
	if string(resp) != "" {
		ret := IncorrectResp{}
		err = json.Unmarshal(resp, &ret)
		if err != nil {
			return err
		}
		log.Error("jpush response incorrect code [%v] message [%v]", ret.Error.Code, ret.Error.Message)
		return fmt.Errorf("%d %s", ret.Error.Code, ret.Error.Message)
	}

	return
}

//从极光服务器删除设备别名(用户退出登录时)
//strAliasName  别名，作为消息推送的唯一标识
func (j *JPush) Unregister(strAliasName string) (err error) {

	httpCli := j.getHttpClientWithAuthorization()
	strDelUrl := fmt.Sprintf("%v/%v", JPUSH_ALIASES_URL, strAliasName)
	log.Debug("post to [%v] delete alias [%v]", strDelUrl, strAliasName)

	resp, err := httpCli.SendUpstream("", "DELETE", strDelUrl) //删除别名

	if err != nil {
		log.Error("delete alias [%v] error [%v]", strAliasName, err)
		return err
	}

	log.Debug("delete alias response message [%v]", string(resp))
	return
}

//APP消息推送: 推送到极光服务器
//platforms 指定平台，空切片内部自动转为所有平台
func (j *JPush) Push(platforms []string, msg *push.Message) (err error) {

	if len(platforms) == 0 { //不指定平台则为三个平台同时发
		platforms = []string{push.PLATFORM_ANDROID, push.PLATFORM_IOS, push.PLATFORM_WINPHONE}
	}
	content := Content{
		Platform: platforms,
	}

	//Notification for android
	content.Notification.Android.Title = msg.Title            //标题
	content.Notification.Android.Alert = msg.Content          //内容
	content.Notification.Android.BuilderId = msg.SoundAndroid //设置安卓声音
	if content.Notification.Android.BuilderId == 0 {
		content.Notification.Android.BuilderId = 2
	}
	content.Notification.Android.Extras.Content = msg.Extra.Content //自定义内容
	content.Notification.Android.Extras.Type = msg.Extra.Type       //自定义类型
	content.Notification.Android.Extras.Url = msg.Extra.Url         //自定义跳转URL

	//Notification for iOS
	content.Notification.IOS.Alert = msg.Title    //必须要有值才能收到
	content.Notification.IOS.Sound = msg.SoundIOS //默认有声音
	content.Notification.IOS.Extras.Content = msg.Extra.Content
	content.Notification.IOS.Extras.Type = msg.Extra.Type //1|2的异或运算
	content.Notification.IOS.Extras.Url = msg.Extra.Url
	content.Options.TimeToLive = 60
	content.Options.ApnsProduction = true //j.is_prod //判断IOS的生产还是测试环境
	//content.Notification.IOS.Badge = 1         //角标默认为空时+1

	content.Audience.Tags = msg.Tags
	content.Audience.Alias = msg.Alias

	log.Struct(&content)
	data, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	strPushContent := string(data)
	if len(msg.Tags) > 0 { //按标签推送(json数据中要删除alias key)
		strPushContent = strings.Replace(strPushContent, `,"alias":null`, "", -1)
	} else { //按别名推送(json数据中要删除tag key)
		strPushContent = strings.Replace(strPushContent, `"tag":null,`, "", -1)
	}

	httpCli := j.getHttpClientWithAuthorization()
	resp, err := httpCli.SendUpstream(strPushContent, "POST", JPUSH_PUSHAPI_URL)
	if err != nil {

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
			return err
		} else {
			log.Error("jpush response [%#v]", ret)
			return errors.New(ret.Error.Message)
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
func (j *JPush) getHttpClientWithAuthorization() (httpCli *httpc.HttpClient) {

	authorization := j.getBase64Authorization()
	httpCli = httpc.NewHttpsC(1)
	httpCli.SetHeader("Content-Type", "application/json")
	httpCli.SetHeader("Authorization", authorization)
	return
}

//按照极光认证要求将appkey和secret做base64编码
func (j *JPush) getBase64Authorization() (strAuthorization string) {

	strEncode := base64.StdEncoding.EncodeToString(([]byte)(fmt.Sprintf("%v:%v", j.appkey, j.secret)))
	strAuthorization = fmt.Sprintf("Basic %v", strEncode)
	log.Debug("authorization [%v]", strAuthorization)
	return
}
