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

func main() {

	//JPushMessage()  //JPUSH(极光)
	//FcmMessage()    //FCM(Google FireBase)
	//UmengMessage()  //Umeng(友盟)
	XingeMessage() //信鸽(tencent)
	//ApnsMessage()   //APNs(Apple)
}

func JPushMessage() {

	var strAppKey string = "2978da262aa372ed199901d2"
	var strSecret string = "1033636bfff8d246294fd1c8"
	var strToken string = "191e35f7e02aeaa0597" //极光token/registered_id 19字节
	var isProdEnv = false                       //sandbox environment is false, production is true for jpush

	jMsg := push.Notification{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strToken},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra: push.PushExtra{
			"type": "1",
			"id":   "10086",
		},
	}

	JPUSH, err := push.GetAdapter(push.AdapterType_JPush, strAppKey, strSecret, isProdEnv)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if JPUSH != nil {
		JPUSH.Debug(true)
		MsgID, err := JPUSH.PushNotification(&jMsg)
		log.Debug("JPUSH msg id [%v] error [%v]", MsgID, err)
	}

}

func FcmMessage() {
	//fake api key, access 'console.firebase.google.com' to get your api key
	strApiKey := "AIzaSyBtMplqJkuTIDyIx_CM74MoPHbxHCBcY-o"
	strToken := "fpICefK-jfE:APA11bHjZTxe503tpFoFCmXXX9LAiMmg7OwgTPYmTb8Ox-yF88umTQnmTQUGbALplxqre7R6v3d0-vSK5MyT4jFtSqklbY1GIaM4d8uZ0wJlwWrRWdBDeOJ4rlpvamd3aGyBlHKAH18N"
	// test for FCM
	fcmMsg := push.Notification{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strToken},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra: push.PushExtra{
			"type": "1",
			"id":   "10086",
		},
	}
	FCM, err := push.GetAdapter(push.AdapterType_Fcm, strApiKey)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if FCM != nil {
		FCM.Debug(true)
		MsgID, err := FCM.PushNotification(&fcmMsg)
		log.Debug("FCM msg id [%v] error [%v]", MsgID, err)
	}
}

func ApnsMessage() {
	var strAuthKeyFile = "AuthKey_U4Q9F3Y9WH.p8"                                      //APNs JWT token auth key file (.p8)
	var strKeyID = "U4Q9F3Y9WH"                                                       //APNs key id
	var strTeamID = "2965AR985S"                                                      //APNs team id
	var strTopic = "com.chatol.thunderchat"                                           //bundle id of your app
	var strToken = "47b9ec1cb4ea5d9e6a62edb183e425a3f2b44715aa26d9848089b1d548e0a0b5" //iOS 64字节
	apnsMsg := push.Notification{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strToken},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra: push.PushExtra{
			"type": "1",
			"id":   "10086",
		},
	}
	APNs, err := push.GetAdapter(push.AdapterType_Apns, strAuthKeyFile, strKeyID, strTeamID, strTopic)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if APNs != nil {
		APNs.Debug(true)
		MsgID, err := APNs.PushNotification(&apnsMsg)
		log.Debug("APNs msg id [%v] error [%v]", MsgID, err)
	}
}

func UmengMessage() {

	var strAppKey = "5e1qfcp30cxfb29u570k01f3"
	var strAppSecret = "jynhjnhhlcgt298ke9q6abpwz3n1j1pv"
	var strToken = "6a17coue30csp91mxy8zafb29f570001f3yaxhw913lx" //iOS 64字节 安卓 44字节
	var strActivity = "org.mychat.ui.LaunchActivity"
	Umeng, err := push.GetAdapter(push.AdatperType_Umeng, strAppKey, strAppSecret, strActivity)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if Umeng != nil {
		strMsgID, err := Umeng.PushNotification(&push.Notification{
			AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
			Audiences:    []string{strToken},
			Platforms:    nil,
			Title:        "this is message title",
			Alert:        "you have a new message",
			SoundIOS:     "",
			SoundAndroid: 0,
			Extra: push.PushExtra{
				"type": "1",
				"id":   "10086",
			},
		})

		log.Debug("MsgID [%v] error [%v]", strMsgID, err)
	}
}

func XingeMessage() {

	var strAppKey string = "b1b86bd2a16b8"
	var strSecret string = "9b99c05cfc8317ea2249e4b5c367z737"
	var strToken string = "6a17coue30csp91mxy8zafb29f570001f3yaxhwz" //iOS 64字节 安卓40字节
	var isProdEnv = false                                            //sandbox environment is false, production is true for iOS

	msg := push.Notification{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strToken},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra: push.PushExtra{
			"type": "1",
			"id":   "10086",
		},
	}

	XinGe, err := push.GetAdapter(push.AdapterType_XinGe, strAppKey, strSecret, isProdEnv)
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	if XinGe != nil {
		XinGe.Debug(true)
		MsgID, err := XinGe.PushNotification(&msg)
		log.Debug("XinGe msg id [%v] error [%v]", MsgID, err)
	}
}
