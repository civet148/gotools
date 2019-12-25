package apns

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/push"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
)

type Apns struct {
	topic  string
	client *apns2.Client
}

type ApnsMessgae struct {
	Token   string //token of apple device
	Message string //message to push
	Type    int32  //chat type (自定义字段：聊天类型 1=单聊 2=群聊 3=频道)
	ChatID  int32  //chat id (单聊：发送人id 群聊：群id 频道：频道id)
}

func init() {

	if err := push.Register(push.AdapterType_Fcm, New); err != nil {
		log.Error("register jpush instance error [%v]", err.Error())
	}
}

var APNS_PARAMS_COUNT = 4

//创建APNs推送接口对象(create APNs push object)
//args[0] => strAuthKeyFile APNs auth key file AuthKey_XXXXX.p8 (string)
//args[1] => strKeyID APNs key id (string)
//args[2] => strTeamID APNs team id (string)
//args[3] => strTopic APNs bundle id of your app (string)
func New(args ...interface{}) push.IPush {

	if len(args) != APNS_PARAMS_COUNT {

		panic(fmt.Errorf("expect %v parameters, got %v", APNS_PARAMS_COUNT, len(args))) //参数个数错误(wrong parameters count)
	}
	strAuthKeyFile := args[0].(string)
	authKey, err := token.AuthKeyFromFile(strAuthKeyFile)
	if err != nil {
		log.Error("APNs: load auth key from file path [%v] error [%v]", strAuthKeyFile, err.Error())
		panic(err.Error())
	}

	return &Apns{

		client: apns2.NewTokenClient(&token.Token{
			AuthKey: authKey,
			KeyID:   args[1].(string),
			TeamID:  args[2].(string),
		}),
		topic: args[3].(string),
	}
}

//push message to device
func (a *Apns) Push(msg *push.Message) (err error) {

	if msg.AudienceType != push.AUDIENCE_TYPE_REGID_TOKEN {
		log.Error("APNs just can use AUDIENCE_TYPE_REGID_TOKEN to push message")
		return fmt.Errorf("APNs just can use AUDIENCE_TYPE_REGID_TOKEN to push message")
	}

	var notification = &apns2.Notification{
		DeviceToken: msg.Audiences[0],
		Topic:       a.topic,
	}

	Payload := payload.NewPayload().
                       Alert(fmt.Sprintf("%v", msg.Alert)). //消息内容(alert of push)
                       Badge(1) //角标+1

	m := push.StructToMap(msg.Extra)
	for k, v := range m {
		Payload.Custom(k, v) //自定义字段(custom key-value)
	}
	notification.Payload = Payload

	var resp *apns2.Response
	if resp, err = a.client.Push(notification); err != nil {
		log.Error("APNs: push message error [%v]", err.Error())
		return
	}

	if resp.StatusCode == 200 {

		log.Debug("APNs: response ok [%+v]", resp)

	} else {
		log.Error("APNs: response error [%+v]", resp)
	}
	return
}

//enable or disable debug output
func (n *Apns) Debug(enable bool) {
	if enable {
		log.SetLevel(0)
	} else {
		log.SetLevel(1)
	}
}
