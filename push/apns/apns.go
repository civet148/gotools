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
	Topic  string
	Client *apns2.Client
}

func init() {

	if err := push.Register(push.AdapterType_Apns, newAPNS); err != nil {
		log.Error("register %v instance error [%v]", push.AdapterType_Apns, err.Error())
		panic("register instance failed")
	}
}

var APNS_PARAMS_COUNT = 5

//创建APNs推送接口对象(create APNs push object)
//args[0] => strAuthKeyFile (string) APNs auth key file AuthKey_XXXXX.p8
//args[1] => strKeyID (string) APNs key id
//args[2] => strTeamID (string) APNs team id
//args[3] => strTopic (string) APNs bundle id of your app
//args[4] => bProd (bool)   APNs production environment
func newAPNS(args ...interface{}) push.IPush {

	if len(args) != APNS_PARAMS_COUNT {
		panic(fmt.Errorf("expect %v parameters, got %v", APNS_PARAMS_COUNT, len(args))) //参数个数错误(wrong parameters count)
	}
	strAuthKeyFile := args[0].(string)
	authKey, err := token.AuthKeyFromFile(strAuthKeyFile)
	if err != nil {
		log.Error("APNs: load auth key from file path [%v] error [%v]", strAuthKeyFile, err.Error())
		panic(err.Error())
	}

	ac := &Apns{

		Client: apns2.NewTokenClient(&token.Token{
			AuthKey: authKey,
			KeyID:   args[1].(string),
			TeamID:  args[2].(string),
		}),
		Topic: args[3].(string),
	}

	if args[4].(bool) {
		ac.Client.Production()
	}
	return ac
}

//push message to device
func (a *Apns) PushNotification(msg *push.Notification) (MsgID string, err error) {

	if msg.AudienceType != push.AUDIENCE_TYPE_REGID_TOKEN {

		err = fmt.Errorf("APNs just can use AUDIENCE_TYPE_REGID_TOKEN to push message")
		log.Error("%v", err.Error())
		return
	}

	var notification = &apns2.Notification{
		DeviceToken: msg.Audiences[0],
		Topic:       a.Topic,
	}

	Payload := payload.NewPayload().
		Alert(msg.Alert). //消息内容(alert of push)
		Badge(msg.Badge)  //角标+1

	for k, v := range msg.Extra {
		Payload.Custom(k, v) //自定义字段(custom key-value)
	}
	notification.Payload = Payload

	var resp *apns2.Response
	if resp, err = a.Client.Push(notification); err != nil {
		log.Error("APNs: push message error [%v]", err.Error())
		return
	}

	if resp.StatusCode == 200 {

		MsgID = resp.ApnsID
		log.Debug("APNs: response ok [%+v]", resp)
	} else {
		log.Error("APNs: response error [%+v]", resp)
	}
	return
}

//enable or disable debug output
func (a *Apns) Debug(enable bool) {
	if enable {
		log.SetLevel(0)
	} else {
		log.SetLevel(1)
	}
}
