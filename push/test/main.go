package main

import (
	"github.com/civet148/gotools/push"
	_ "github.com/civet148/gotools/push/apns"
	_ "github.com/civet148/gotools/push/fcm"
	_ "github.com/civet148/gotools/push/jpush"
)




type PushExtra struct {
	Type   int32 `json:"type"`
	ChatID int32 `json:"chat_id"`
}

func main() {

	var strRegId string = "191e35f7e02aeaa0597"
	var strAppKey string = "2978da262aa372ed199901d2"
	var strSecret string = "1033636bfff8d246294fd1c8"
	var isProdEnv = false //sandbox environment is false, production is true for jpush

	msg := push.Message{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strRegId},
		Title:        "你有一条新消息",
		Content:      "洞拐洞拐，我是电报后台，收到请回复+1",
		Extra:        &PushExtra{Type: 1, ChatID: 10000},
	}

	jpush, _ := push.GetAdapter(push.AdapterType_JPush, strAppKey, strSecret, isProdEnv)
	if jpush != nil {
		jpush.Debug(true)
		_ = jpush.Push(&msg)
	}


	//fake api key, access 'console.firebase.google.com' to get your api key
	strApiKey := "AIzaSyBtMplqJkuTIDyIx_CM74MoPHbxHCBcY-o"
	fcm, _ := push.GetAdapter(push.AdapterType_Fcm, strApiKey)

	if fcm != nil {
		fcm.Debug(true)
		_ = fcm.Push(&msg)
	}

	var strAuthKeyFile="AuthKey_U4Q9F3Y9WH.p8" //APNs JWT token auth key file (.p8)
	var strKeyID = "U4Q9F3Y9WH" //APNs key id
	var strTeamID = "2965AR985S" //APNs team id
	var strTopic = "com.chatol.thunderchat" //bundle id of your app
	apns, _ := push.GetAdapter(push.AdapterTYpe_Apns, strAuthKeyFile, strKeyID, strTeamID, strTopic)
	if apns != nil {
		apns.Debug(true)
		_ = apns.Push(&msg)
	}
}
