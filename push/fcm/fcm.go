// Copyright (c) 2019-present,  NebulaChat Studio (https://nebula.chat).
//  All rights reserved.
//
// Author: Benqi (wubenqi@gmail.com)
//

package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/push"
	"io/ioutil"
	"net/http"
	"time"
)

var FMC_PARAMS_COUNT = 1

// Message represents fcm request message
type (
	Fcm struct {
		fcm_cli *Client
	}

	Message struct {
		// Data parameter specifies the custom key-value pairs of the message's payload.
		//
		// For example, with data:{"score":"3x1"}:
		//
		// On iOS, if the message is sent via APNS, it represents the custom data fields.
		// If it is sent via FCM connection server, it would be represented as key value dictionary
		// in AppDelegate application:didReceiveRemoteNotification:.
		// On Android, this would result in an intent extra named score with the string value 3x1.
		// The key should not be a reserved word ("from" or any word starting with "google" or "gcm").
		// Do not use any of the words defined in this table (such as collapse_key).
		// Values in string types are recommended. You have to convert values in objects
		// or other non-string data types (e.g., integers or booleans) to string.
		//
		Data interface{} `json:"data,omitempty"`

		// To this parameter specifies the recipient of a message.
		//
		// The value must be a registration token, notification key, or topic.
		// Do not set this field when sending to multiple topics. See Condition.
		To string `json:"to,omitempty"`

		// RegistrationIDs for all registration ids
		// This parameter specifies a list of devices
		// (registration tokens, or IDs) receiving a multicast message.
		// It must contain at least 1 and at most 1000 registration tokens.
		// Use this parameter only for multicast messaging, not for single recipients.
		// Multicast messages (sending to more than 1 registration tokens)
		// are allowed using HTTP JSON format only.
		RegistrationIDs []string `json:"registration_ids,omitempty"`

		// CollapseKey This parameter identifies a group of messages
		// (e.g., with collapse_key: "Updates Available") that can be collapsed,
		// so that only the last message gets sent when delivery can be resumed.
		// This is intended to avoid sending too many of the same messages when the
		// device comes back online or becomes active (see delay_while_idle).
		CollapseKey string `json:"collapse_key,omitempty"`

		// Priority Sets the priority of the message. Valid values are "normal" and "high."
		// On iOS, these correspond to APNs priorities 5 and 10.
		// By default, notification messages are sent with high priority, and data messages
		// are sent with normal priority. Normal priority optimizes the client app's battery
		// consumption and should be used unless immediate delivery is required. For messages
		// with normal priority, the app may receive the message with unspecified delay.
		// When a message is sent with high priority, it is sent immediately, and the app
		// can wake a sleeping device and open a network connection to your server.
		// For more information, see Setting the priority of a message.
		Priority string `json:"priority,omitempty"`

		// Notification parameter specifies the predefined, user-visible key-value pairs of
		// the notification payload. See Notification payload support for detail.
		// For more information about notification message and data message options, see
		// Notification
		Notification Notification `json:"notification,omitempty"`

		// ContentAvailable On iOS, use this field to represent content-available
		// in the APNS payload. When a notification or message is sent and this is set
		// to true, an inactive client app is awoken. On Android, data messages wake
		// the app by default. On Chrome, currently not supported.
		ContentAvailable bool `json:"content_available,omitempty"`

		// DelayWhenIdle When this parameter is set to true, it indicates that
		// the message should not be sent until the device becomes active.
		// The default value is false.
		DelayWhileIdle bool `json:"delay_while_idle,omitempty"`

		// TimeToLive This parameter specifies how long (in seconds) the message
		// should be kept in FCM storage if the device is offline. The maximum time
		// to live supported is 4 weeks, and the default value is 4 weeks.
		// For more information, see
		// https://firebase.google.com/docs/cloud-messaging/concept-options#ttl
		TimeToLive int `json:"time_to_live,omitempty"`

		// RestrictedPackageName This parameter specifies the package name of the
		// application where the registration tokens must match in order to
		// receive the message.
		RestrictedPackageName string `json:"restricted_package_name,omitempty"`

		// DryRun This parameter, when set to true, allows developers to test
		// a request without actually sending a message.
		// The default value is false
		DryRun bool `json:"dry_run,omitempty"`

		// Condition to set a logical expression of conditions that determine the message target
		// This parameter specifies a logical expression of conditions that determine the message target.
		// Supported condition: Topic, formatted as "'yourTopic' in topics". This value is case-insensitive.
		// Supported operators: &&, ||. Maximum two operators per topic message supported.
		Condition string `json:"condition,omitempty"`

		// Currently for iOS 10+ devices only. On iOS, use this field to represent mutable-content in the APNS payload.
		// When a notification is sent and this is set to true, the content of the notification can be modified before
		// it is displayed, using a Notification Service app extension. This parameter will be ignored for Android and web.
		MutableContent bool `json:"mutable_content,omitempty"`

		Android Android `json:"android,omitempty"`
	}

	Android struct {
		Priority string `json:"priority,omitempty"`
	}

	// Result Downstream result from FCM, sent in the "results" field of the Response packet
	Result struct {
		// String specifying a unique ID for each successfully processed message.
		MessageID string `json:"message_id"`

		// Optional string specifying the canonical registration token for the
		// client app that the message was processed and sent to. Sender should
		// use this value as the registration token for future requests.
		// Otherwise, the messages might be rejected.
		RegistrationID string `json:"registration_id"`

		// String specifying the error that occurred when processing the message
		// for the recipient. The possible values can be found in table 9 here:
		// https://firebase.google.com/docs/cloud-messaging/http-server-ref#table9
		Error string `json:"error"`
	}

	// Response represents fcm response message - (tokens and topics)
	Response struct {
		Ok         bool
		StatusCode int

		// MulticastID a unique ID (number) identifying the multicast message.
		MulticastID int `json:"multicast_id"`

		// Success number of messages that were processed without an error.
		Success int `json:"success"`

		// Fail number of messages that could not be processed.
		Fail int `json:"failure"`

		// CanonicalIDs number of results that contain a canonical registration token.
		// A canonical registration ID is the registration token of the last registration
		// requested by the client app. This is the ID that the server should use
		// when sending messages to the device.
		CanonicalIDs int `json:"canonical_ids"`

		// Results Array of objects representing the status of the messages processed. The objects are listed in the same order as the request (i.e., for each registration ID in the request, its result is listed in the same index in the response).
		// message_id: String specifying a unique ID for each successfully processed message.
		// registration_id: Optional string specifying the canonical registration token for the client app that the message was processed and sent to. Sender should use this value as the registration token for future requests. Otherwise, the messages might be rejected.
		// error: String specifying the error that occurred when processing the message for the recipient. The possible values can be found in table 9.
		Results []Result `json:"results,omitempty"`

		// The topic message ID when FCM has successfully received the request and will attempt to deliver to all subscribed devices.
		MsgID int `json:"message_id,omitempty"`

		// Error that occurred when processing the message. The possible values can be found in table 9.
		Err string `json:"error,omitempty"`

		// RetryAfter
		RetryAfter string
	}

	// Notification notification message payload
	Notification struct {
		// Title indicates notification title. This field is not visible on iOS phones and tablets.
		Title string `json:"title,omitempty"`

		// Body indicates notification body text.
		Body string `json:"body,omitempty"`

		// Sound indicates a sound to play when the device receives a notification.
		// Sound files can be in the main bundle of the client app or in the
		// Library/Sounds folder of the app's data container.
		// See the iOS Developer Library for more information.
		// http://apple.co/2jaGqiE
		Sound string `json:"sound,omitempty"`

		// Badge indicates the badge on the client app home icon.
		Badge string `json:"badge,omitempty"`

		// Icon indicates notification icon. Sets value to myicon for drawable resource myicon.
		// If you don't send this key in the request, FCM displays the launcher icon specified
		// in your app manifest.
		Icon string `json:"icon,omitempty"`

		// Tag indicates whether each notification results in a new entry in the notification
		// drawer on Android. If not set, each request creates a new notification.
		// If set, and a notification with the same tag is already being shown,
		// the new notification replaces the existing one in the notification drawer.
		Tag string `json:"tag,omitempty"`

		// Color indicates color of the icon, expressed in #rrggbb format
		Color string `json:"color,omitempty"`

		// ClickAction indicates the action associated with a user click on the notification.
		// When this is set, an activity with a matching intent filter is launched when user
		// clicks the notification.
		ClickAction string `json:"click_action,omitempty"`

		// BodyLockKey indicates the key to the body string for localization. Use the key in
		// the app's string resources when populating this value.
		BodyLocKey string `json:"body_loc_key,omitempty"`

		// BodyLocArgs indicates the string value to replace format specifiers in the body
		// string for localization. For more information, see Formatting and Styling.
		BodyLocArgs string `json:"body_loc_args,omitempty"`

		// TitleLocKey indicates the key to the title string for localization.
		// Use the key in the app's string resources when populating this value.
		TitleLocKey string `json:"title_loc_key,omitempty"`

		// TitleLocArgs indicates the string value to replace format specifiers in the title string for
		// localization. For more information, see
		// https://developer.android.com/guide/topics/resources/string-resource.html#FormattingAndStyling
		TitleLocArgs string `json:"title_loc_args,omitempty"`
	}
)

const (
	// PriorityHigh used for high notification priority
	PriorityHigh = "high"

	// PriorityNormal used for normal notification priority
	PriorityNormal = "normal"

	// HeaderRetryAfter HTTP header constant
	HeaderRetryAfter = "Retry-After"

	// ErrorKey readable error caching
	ErrorKey = "error"

	// MethodPOST indicates http post method
	MethodPOST = "POST"

	// ServerURL push server url
	ServerURL = "https://fcm.googleapis.com/fcm/send"
)

// retryableErrors whether the error is a retryable
var retryableErrors = map[string]bool{
	"Unavailable":         true,
	"InternalServerError": true,
}

// Client stores client with api key to firebase
type Client struct {
	apiKey  string
	httpCli *http.Client
}

func init() {

	if err := push.Register(push.AdapterType_Fcm, newFCM); err != nil {
		log.Error("register %v instance error [%v]", push.AdapterType_Fcm, err.Error())
		panic("register instance failed")
	}
}

// NewClient creates a new client
func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		apiKey:  apiKey,
		httpCli: &http.Client{Timeout: timeout},
	}
}

func (f *Client) authorization() string {
	return fmt.Sprintf("key=%v", f.apiKey)
}

// Send sends message to FCM
func (f *Client) Send(message *Message) (*Response, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return &Response{}, err
	}
	req, err := http.NewRequest(MethodPOST, ServerURL, bytes.NewBuffer(data))
	if err != nil {
		return &Response{}, err
	}
	req.Header.Set("Authorization", f.authorization())
	req.Header.Set("Content-Type", "application/json")
	resp, err := f.httpCli.Do(req)
	if err != nil {
		// fmt.Println(err)
		return &Response{}, err
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	response := &Response{StatusCode: resp.StatusCode}
	if resp.StatusCode >= 500 {
		response.RetryAfter = resp.Header.Get(HeaderRetryAfter)
	}
	if resp.StatusCode != 200 {
		return response, fmt.Errorf("fcm status code(%d)", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	if err := f.Failed(response); err != nil {
		return response, err
	}
	response.Ok = true
	return response, nil
}

// Failed method indicates if the server couldn't process
// the request in time.
func (f *Client) Failed(response *Response) error {
	for _, response := range response.Results {
		if retryableErrors[response.Error] {
			return fmt.Errorf("fcm push error(%s)", response.Error)
		}
	}
	return nil
}

//创建极光推送接口对象
//args[0] => apikey 开发者API密钥(string)
func newFCM(args ...interface{}) push.IPush {

	if len(args) != FMC_PARAMS_COUNT {

		panic(fmt.Errorf("expect %v parameters, got %v", FMC_PARAMS_COUNT, len(args))) //参数个数错误
	}

	return &Fcm{

		fcm_cli: NewClient(args[0].(string), 5*time.Second),
	}
}

//push message to device (by device token or register id)
func (f *Fcm) PushNotification(msg *push.Notification) (MsgID string, err error) {

	if msg.AudienceType != push.AUDIENCE_TYPE_REGID_TOKEN {
		err = fmt.Errorf("FCM just can use AUDIENCE_TYPE_REGID_TOKEN to push message")
		log.Error("FCM just can use AUDIENCE_TYPE_REGID_TOKEN to push message")
		return
	}

	fcmMsg := &Message{
		// DryRun:          true, // 如果是 true，消息不会下发给用户，用于测试
		Data:            msg.Extra,
		RegistrationIDs: msg.Audiences,
		Priority:        PriorityHigh,
		DelayWhileIdle:  true,
		Notification: Notification{
			Title: msg.Title,
			Body:  msg.Alert,
			//ClickAction: "com.bilibili.app.in.com.bilibili.push.FCM_MESSAGE", //点击触发事件，暂时不支持
		},
		//CollapseKey: strings.TrimFunc("t123456", func(r rune) bool {
		//	return !unicode.IsNumber(r)
		//}), // 消息分组, 值转成 int 传到客户端(暂时不支持分组)
		TimeToLive: int(time.Hour.Seconds()),
		Android:    Android{Priority: PriorityHigh},
	}

	response, err := f.fcm_cli.Send(fcmMsg)
	if err != nil {
		log.Error("pushToFcm error [%v]", err.Error())
		return
	}
	if response.Ok {
		MsgID = fmt.Sprintf("%v", response.MsgID)
		log.Debug("pushToFcm response ok [%+v]", response)
	} else {
		log.Error("pushToFcm response error [%+v]", response)
	}
	return
}

//enable or disable debug output
func (f *Fcm) Debug(enable bool) {
	if enable {
		log.SetLevel(0)
	} else {
		log.SetLevel(1)
	}
}

// GetRetryAfterTime converts the retry after response header to a time.Duration
func (r *Response) GetRetryAfterTime() (time.Duration, error) {
	return time.ParseDuration(r.RetryAfter)
}
