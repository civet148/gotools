package main

import (
	"github.com/civet148/gotools/push"
	_ "github.com/civet148/gotools/push/jpush"
)

var strAlias string = "10000007"
var strRegId string = "191e35f7e02aeaa0597"
var strAppKey string = "2978da262aa372ed199901d2"
var strSecret string = "1033636bfff8d246294fd1c8"

func main() {

	var isProdEnv = false //测试环境为false，正式环境为true
	jpush, _ := push.GetAdapter(push.AdapterType_JPush, strAppKey, strSecret, isProdEnv)
	err := jpush.Register(strAlias, strRegId, "17788880000")

	if err == nil {
		msg := push.Message{
			//Tags:         nil,
			Alias: []string{strAlias},
			//RegId:		  strRegId,
			Title:        "你有一条新消息",
			Content:      "洞拐洞拐，我是电报后台，收到请回复+1",
			SoundIOS:     "",
			SoundAndroid: 0,
			Extra:        push.Extras{Content: "洞拐洞拐，我是电报后台，收到请回复+2"},
		}

		_ = jpush.Push([]string{push.PLATFORM_ANDROID, push.PLATFORM_IOS}, &msg)
		//_ = jpush.Unregister(strAlias)
	}
}
