package push

import "fmt"

type AdapterType int

var (
	AdapterType_JPush  AdapterType = 1 //极光推送
	AdapterType_Google AdapterType = 2 //谷歌推送
)

var (
	PLATFORM_ALL      = "all"
	PLATFORM_ANDROID  = "android"
	PLATFORM_IOS      = "ios"
	PLATFORM_WINPHONE = "winphone"
)

type Message struct {
	Tags         []string //接收者的TAG标签（极光推送：tags和alias同时存在的情况下只发送给指定tag的接收者）
	Alias        []string //接收者的别名
	RegId        string   //极光注册ID
	Title        string   //标题
	Content      string   //内容
	SoundIOS     string   //iOS声音设置 无 默认 自定义(空字符串、声音文件名、default)
	SoundAndroid int      //声音设置(0 默认 1 播放声音文件)
	Extra        Extras   //扩展结构，用于点击推送消息自动跳转页面(选填)
}

//扩展结构，用于点击推送消息自动跳转页面
type Extras struct {
	Url     string `json:"url"`
	Type    int32  `json:"type"`
	Content string `json:"content"`
}

type IPush interface {
	//app event: user signed in from a mobile device
	Register(strAliasName, strRegisterId, strMobile string) (err error)
	//app event: user signed out from a mobile device
	Unregister(strAliasName string) (err error)
	//push message to device
	Push(platforms []string, msg *Message) (err error)
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
