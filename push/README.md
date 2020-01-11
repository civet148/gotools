# [APP消息推送]

## 极光(JPush)

* API集成文档官方URL https://docs.jiguang.cn/jpush/server/push/rest_api_v3_push/

## 谷歌FCM(FireBase)

* 开发者官方URL https://firebase.google.com/

## 苹果APNs

* 开发者官方URL https://developer.apple.com/

## 信鸽(XinGe)

支持 FCM/HUAWEI/XIAOMI/MEIZU 厂商通道(需在信鸽开发者平台开启厂商通道并配置相关厂商的AppKey+AppSecret并且APP端需集成相关厂商的推送SDK)

* API集成文档官方URL https://xg.qq.com/docs/server_api/v3/rest_api_summary.html

## 友盟推送

支持 HUAWEI/XIAOMI/MEIZU/VIVO/OPPO(不支持Oppo v8.0以上系统) 厂商通道(服务端需在发送推送通知报文中指定mipush=true并在友盟开发者平台配置相关厂商的AppKey+AppSecret并且APP端需集成相关厂商的推送SDK)


* API集成文档官方URL https://developer.umeng.com/docs/66632/detail/68343

## 代码示例

```golang

package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/push"
	_ "github.com/civet148/gotools/push/apns"
	_ "github.com/civet148/gotools/push/fcm"
	_ "github.com/civet148/gotools/push/jpush"
	_ "github.com/civet148/gotools/push/umeng"
	_ "github.com/civet148/gotools/push/xinge"
)

type PushExtra struct {
	Type   int32 `json:"type"`
	ChatID int32 `json:"chat_id"`
}

func main() {

	//JPushMessage()  //JPUSH(极光)
	//FcmMessage()    //FCM(Google FireBase)
	//UmengMessage()  //Umeng(友盟)
	XingeMessage()  //信鸽(tencent)
	//ApnsMessage()   //APNs(Apple)
}

func JPushMessage() {


	var strAppKey string = "2978da262aa372ed199901d2"
	var strSecret string = "1033636bfff8d246294fd1c8"
	var strToken string = "191e35f7e02aeaa0597" //极光token/registered_id 19字节
	var isProdEnv = false //sandbox environment is false, production is true for jpush

	jMsg := push.Message{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strToken},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra:        &PushExtra{Type: 1, ChatID: 10086},
	}

	JPUSH, err := push.GetAdapter(push.AdapterType_JPush, strAppKey, strSecret, isProdEnv)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if JPUSH != nil {
		JPUSH.Debug(true)
		MsgID, err := JPUSH.Push(&jMsg)
		log.Debug("JPUSH msg id [%v] error [%v]", MsgID, err)
	}

}

func FcmMessage() {
	//fake api key, access 'console.firebase.google.com' to get your api key
	strApiKey := "AIzaSyBtMplqJkuTIDyIx_CM74MoPHbxHCBcY-o"
	strToken := "fpICefK-jfE:APA11bHjZTxe503tpFoFCmXXX9LAiMmg7OwgTPYmTb8Ox-yF88umTQnmTQUGbALplxqre7R6v3d0-vSK5MyT4jFtSqklbY1GIaM4d8uZ0wJlwWrRWdBDeOJ4rlpvamd3aGyBlHKAH18N"
	// test for FCM
	fcmMsg := push.Message{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strToken},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra:        &PushExtra{Type: 1, ChatID: 10086},
	}
	FCM, err := push.GetAdapter(push.AdapterType_Fcm, strApiKey)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if FCM != nil {
		FCM.Debug(true)
		MsgID, err := FCM.Push(&fcmMsg)
		log.Debug("FCM msg id [%v] error [%v]", MsgID, err)
	}
}

func ApnsMessage() {
	var strAuthKeyFile = "AuthKey_U4Q9F3Y9WH.p8" //APNs JWT token auth key file (.p8)
	var strKeyID = "U4Q9F3Y9WH"                  //APNs key id
	var strTeamID = "2965AR985S"                 //APNs team id
	var strTopic = "com.chatol.thunderchat"      //bundle id of your app
	var strToken = "47b9ec1cb4ea5d9e6a62edb183e425a3f2b44715aa26d9848089b1d548e0a0b5" //iOS 64字节
	apnsMsg := push.Message{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strToken},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra:        &PushExtra{Type: 1, ChatID: 10086},
	}
	APNs, err := push.GetAdapter(push.AdapterType_Apns, strAuthKeyFile, strKeyID, strTeamID, strTopic)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if APNs != nil {
		APNs.Debug(true)
		MsgID, err := APNs.Push(&apnsMsg)
		log.Debug("APNs msg id [%v] error [%v]", MsgID, err)
	}
}

func UmengMessage() {
	
	var strAppKey = "5e1qfcp30cxfb29u570k01f3"
	var strAppSecret = "jynhjnhhlcgt298ke9q6abpwz3n1j1pv"
	var strToken = "6a17coue30csp91mxy8zafb29f570001f3yaxhw913lx" //iOS 64字节 安卓 44字节
	Umeng, err := push.GetAdapter(push.AdatperType_Umeng, strAppKey, strAppSecret)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if Umeng != nil {
		Umeng.Push(&push.Message{
			AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
			Audiences:    []string{strToken},
			Platforms:    nil,
			Title:        "this is message title",
			Alert:        "you have a new message",
			SoundIOS:     "",
			SoundAndroid: 0,
			Extra:         &PushExtra{Type: 1, ChatID: 10086},
		})
	}
}

func XingeMessage() {


	var strAppKey string = "b1b86bd2a16b8"
	var strSecret string = "zbzzc05cfc8317ea224ze4b5c367z737"
	var strToken string = "6a17coue30csp91mxy8zafb29f570001f3yaxhwz" //iOS 64字节 安卓40字节
	var isProdEnv = false //sandbox environment is false, production is true for iOS

	msg := push.Message{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strToken},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra:        &PushExtra{Type: 1, ChatID: 10086},
	}

	XinGe, err := push.GetAdapter(push.AdapterType_XinGe, strAppKey, strSecret, isProdEnv)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if XinGe != nil {
		XinGe.Debug(true)
		MsgID, err := XinGe.Push(&msg)
		log.Debug("XinGe msg id [%v] error [%v]", MsgID, err)
	}
}
```