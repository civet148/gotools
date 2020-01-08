package push

import (
	"encoding/json"
	"fmt"
)

type AdapterType int

var (
	HTTP_METHOD_GET    = "GET"
	HTTP_METHOD_POST   = "POST"
	HTTP_METHOD_DELETE = "DELETE"
	HTTP_METHOD_UPLOAD = "UPLOAD"
)

var (
	AdapterType_JPush AdapterType = 1 //极光推送
	AdapterType_Fcm   AdapterType = 2 //谷歌推送
	AdapterType_Apns  AdapterType = 3 //苹果推送
	AdatperType_Umeng AdapterType = 4 //友盟推送
	AdapterType_XinGe AdapterType = 5 //信鸽推送
)

func (t AdapterType) String() (name string) {
	switch t {
	case AdapterType_JPush:
		name = "JPUSH"
	case AdapterType_Fcm:
		name = "FCM"
	case AdapterType_Apns:
		name = "APNs"
	case AdatperType_Umeng:
		name = "Umeng"
	case AdapterType_XinGe:
		name = "XinGe"
	default:
		name = "<unknown>"
	}
	return
}

func (t AdapterType) GoString() (name string) {
	return t.String()
}

var (
	PLATFORM_ALL      = "all"
	PLATFORM_ANDROID  = "android"
	PLATFORM_IOS      = "ios"
	PLATFORM_WINPHONE = "winphone"
)

type AudienceType int

var (
	AUDIENCE_TYPE_REGID_TOKEN AudienceType = 1 //按设备注册ID/token推送(device register id or device token)
	AUDIENCE_TYPE_TAG         AudienceType = 2 //按标签推送(message group)
)

type Message struct {
	AudienceType AudienceType //推送类型
	Audiences    []string     //设备注册ID或标签或别名
	Platforms    []string     //推送平台（可为空）
	Title        string       //标题
	Alert        string       //内容
	SoundIOS     string       //iOS声音设置 无 默认 自定义(空字符串、声音文件名、default)
	SoundAndroid int          //声音设置(0 默认 1 播放声音文件)
	Extra        interface{}  //扩展结构，用于点击推送消息自动跳转页面(选填)
}

type IPush interface {

	//push message to device
	Push(msg *Message) (MsgID string, err error)
	//enable or disable debug output
	Debug(enable bool)
}

type Instance func(args ...interface{}) IPush

var AdapterMap = make(map[AdapterType]Instance)

//register adapter instance
func Register(adapter AdapterType, ins Instance) (err error) {

	if _, ok := AdapterMap[adapter]; !ok {

		AdapterMap[adapter] = ins
		return
	}
	err = fmt.Errorf("adapter [%v] instance already exists", adapter)
	return
}

//get adapter instance with args...
func GetAdapter(adapter AdapterType, args ...interface{}) (IPush, error) {

	ins, ok := AdapterMap[adapter]
	if !ok {
		return nil, fmt.Errorf("adapter [%v] instance not exists", adapter)
	}

	return ins(args...), nil
}

func StructToMap(v interface{}) (m map[string]interface{}) {

	inrec, _ := json.Marshal(v)
	_ = json.Unmarshal(inrec, &m)

	return
}
