package apns

import (
	"fmt"
	"nebula.chat/enterprise/pkg/log"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
)

type Apns struct {
	topic  string
	client *apns2.Client
}

type ApnsMessgae struct {
	Token   string   //token of apple device
	Message string   //message to push
	Type    int32    //chat type (自定义字段：聊天类型 1=单聊 2=群聊 3=频道)
	ChatID  int32    //chat id (单聊：发送人id 群聊：群id 频道：频道id)
}

//使用APNs token客户端推送(JWT)
func NewApns(strAuthKeyFile, strKeyID, strTeamID, strTopic string) *Apns {
	authKey, err := token.AuthKeyFromFile(strAuthKeyFile)
	if err != nil {
		log.Errorf("APNs: load auth key from file path [%v] error [%v]", strAuthKeyFile, err.Error())
		return nil
	}

	return &Apns{
		topic: strTopic,
		client: apns2.NewTokenClient(&token.Token{
			AuthKey: authKey,
			KeyID:   strKeyID,
			TeamID:  strTeamID,
		}),
	}
}

func (a *Apns) Push(msg *ApnsMessgae) (err error) {

	var notification = &apns2.Notification{
		DeviceToken: msg.Token,
		Topic:   a.topic,
	}

	notification.Payload = payload.NewPayload().
		Alert(fmt.Sprintf("%v", msg.Message)). //消息内容
		Badge(1). //角标+1
		Custom("type", msg.Type). //自定义字段
		Custom("chat_id", msg.ChatID) //自定义字段

	var resp *apns2.Response
	if resp, err = a.client.Push(notification); err != nil {
		log.Errorf("APNs: push message error [%v]", err.Error())
		return
	}

	if resp.StatusCode == 200 {
		log.Debugf("APNs: response ok [%+v]", resp)
	} else {
		log.Errorf("APNs: response error [%+v]", resp)
	}

	return
}
