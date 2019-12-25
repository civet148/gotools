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

	var strJpushRegId string = "191e35f7e02aeaa0597"
	var strAppKey string = "2978da262aa372ed199901d2"
	var strSecret string = "1033636bfff8d246294fd1c8"
	var isProdEnv = false //sandbox environment is false, production is true for jpush

	jMsg := push.Message{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strJpushRegId},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra:        &PushExtra{Type: 1, ChatID: 10000},
	}

	JPUSH, _ := push.GetAdapter(push.AdapterType_JPush, strAppKey, strSecret, isProdEnv)
	if JPUSH != nil {
		JPUSH.Debug(true)
		_ = JPUSH.Push(&jMsg)
	}

	//fake api key, access 'console.firebase.google.com' to get your api key
	strApiKey := "AIzaSyBtMplqJkuTIDyIx_CM74MoPHbxHCBcY-o"
	strFcmRegId := "fpICefK-jfE:APA11bHjZTxe503tpFoFCmXXX9LAiMmg7OwgTPYmTb8Ox-yF88umTQnmTQUGbALplxqre7R6v3d0-vSK5MyT4jFtSqklbY1GIaM4d8uZ0wJlwWrRWdBDeOJ4rlpvamd3aGyBlHKAH18N"
	// test for FCM
	fcmMsg := push.Message{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strFcmRegId},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra:        &PushExtra{Type: 1, ChatID: 10000},
	}
	FCM, _ := push.GetAdapter(push.AdapterType_Fcm, strApiKey)
	if FCM != nil {
		FCM.Debug(true)
		_ = FCM.Push(&fcmMsg)
	}

	var strAuthKeyFile = "AuthKey_U4Q9F3Y9WH.p8" //APNs JWT token auth key file (.p8)
	var strKeyID = "U4Q9F3Y9WH"                  //APNs key id
	var strTeamID = "2965AR985S"                 //APNs team id
	var strTopic = "com.chatol.thunderchat"      //bundle id of your app
	var strDeviceToken = "47b9ec1cb4ea5d9e6a62edb183e425a3f2b44715aa26d9848089b1d548e0a0b5"
	apnsMsg := push.Message{
		AudienceType: push.AUDIENCE_TYPE_REGID_TOKEN,
		Audiences:    []string{strDeviceToken},
		Title:        "this is message title",
		Alert:        "you have a new message",
		Extra:        &PushExtra{Type: 1, ChatID: 10000},
	}
	APNs, _ := push.GetAdapter(push.AdapterTYpe_Apns, strAuthKeyFile, strKeyID, strTeamID, strTopic)
	if APNs != nil {
		APNs.Debug(true)
		_ = APNs.Push(&apnsMsg)
	}
}
