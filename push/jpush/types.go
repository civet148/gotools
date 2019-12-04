package jpush


type Content struct {
	Platform     []string     `json:"platform"` //需要推送的平台
	Audience     Audience     `json:"audience"` //需要推送的接受者
	Notification Notification `json:"notification"`
	Options      Options      `json:"options"`
}

type Audience struct {
	Tags  []string `json:"tag"`   //接收者的TAG标签
	Alias []string `json:"alias"` //接收者的别名
}

type Options struct {
	TimeToLive    int  `json:"time_to_live"`    //接收者的TAG标签
	ApnsProduction bool `json:"apns_production"` //接收者的别名
}

type Notification struct { //真正的平台通知接收者的内容
	Android Android `json:"android"`
	IOS     IOS     `json:"ios"`
}

type Android struct {
	Alert      string `json:"alert"`
	Title      string `json:"title"`
	BuilderId int    `json:"builder_id"`
	Extras     Extras `json:"extras"`
}

type IOS struct {
	Alert  string `json:"alert"`
	Extras Extras `json:"extras"`
	Sound  string `json:"sound"`
	//Badge  int    `json:"badge"`
}

//自定义字段，选填
type Extras struct {
	Url     string `json:"url"`
	Type    int32  `json:"type"`
	Content string `json:"content"`
	Sound   string `json:"sound"` //新增声音文件
}

type CorrectResp struct {
	SendNo string `json:"sendno"`
	MsgId int64  `json:"msg_id"`
}

type IncorrectResp struct {
	Error  Error `json:"error"`
	MsgId int64 `json:"msg_id"`
}

type Error struct {
	Message string `json:"message"`
	Code    int32  `json:"code"`
}

/*
新增注册
*/
type Register struct {
	Mobile string `json:"mobile"`
	Alias string `json:"alias"`
}