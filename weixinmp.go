package weixinmp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

const (
	// request types
	MsgTypeText      = "text"
	MsgTypeImage     = "image"
	MsgTypeVoice     = "voice"
	MsgTypeVideo     = "video"
	MsgTypeLocation  = "location"
	MsgTypeLink      = "link"
	MsgTypeEvent     = "event"
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventScan        = "SCAN"
	EventLocation    = "LOCATION"
	EventClick       = "CLICK"
	EventView        = "VIEW"
	// environment variables
	retryNum        = 3
	plainPreUrl     = "https://api.weixin.qq.com/cgi-bin/"
	mediaPreUrl     = "http://file.api.weixin.qq.com/cgi-bin/media/"
	accessTokenTemp = "accesstoken.temp"
)

type qrScene struct {
	ExpireSeconds int64  `json:"expire_seconds,omitempty"`
	ActionName    string `json:"action_name"`
	ActionInfo    struct {
		Scene struct {
			SceneId int64 `json:"scene_id"`
		} `json:"scene"`
	} `json:"action_info"`
}

// message structs
type textMsg struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `json:"touser"`
	FromUserName string   `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      string   `json:"msgtype"`
	Content      string   `json:"-"`
	Text         struct {
		Content string `xml:"-" json:"content"`
	} `xml:"-" json:"text"`
}

type imageMsg struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `json:"touser"`
	FromUserName string   `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      string   `json:"msgtype"`
	Image        struct {
		MediaId string `json:"media_id"`
	} `json:"image"`
}

type voiceMsg struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `json:"touser"`
	FromUserName string   `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      string   `json:"msgtype"`
	Voice        struct {
		MediaId string `json:"media_id"`
	} `json:"voice"`
}

type videoMsg struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `json:"touser"`
	FromUserName string   `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      string   `json:"msgtype"`
	Video        *Video   `json:"video"`
}

type musicMsg struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `json:"touser"`
	FromUserName string   `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      string   `json:"msgtype"`
	Music        *Music   `json:"music"`
}

type newsMsg struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `json:"touser"`
	FromUserName string   `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      string   `json:"msgtype"`
	ArticleCount int      `json:"-"`
	Articles     struct {
		Item *[]Article `xml:"item" json:"articles"`
	} `json:"news"`
}

type Video struct {
	MediaId     string `json:"media_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Music struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	MediaId      string `json:"media_id"`
	MusicUrl     string `json:"musicurl"`
	HQMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PicUrl      string `json:"picurl"`
	Url         string `json:"url"`
}

// weixinmp goes here
type Weixinmp struct {
	accessToken *accessToken
	Request     *Request
}

func New(token, appid, secret string) *Weixinmp {
	inst := &Weixinmp{
		accessToken: &accessToken{appid: appid, secret: secret},
		Request:     &Request{token: token},
	}
	return inst
}

// reply message methods
func (this *Weixinmp) ReplyTextMsg(rw http.ResponseWriter, content string) error {
	var msg textMsg
	msg.MsgType = "text"
	msg.Content = content
	return this.replyMsg(rw, &msg)
}

func (this *Weixinmp) ReplyImageMsg(rw http.ResponseWriter, mediaId string) error {
	var msg imageMsg
	msg.MsgType = "image"
	msg.Image.MediaId = mediaId
	return this.replyMsg(rw, &msg)
}

func (this *Weixinmp) ReplyVoiceMsg(rw http.ResponseWriter, mediaId string) error {
	var msg voiceMsg
	msg.MsgType = "voice"
	msg.Voice.MediaId = mediaId
	return this.replyMsg(rw, &msg)
}

func (this *Weixinmp) ReplyVideoMsg(rw http.ResponseWriter, video *Video) error {
	var msg videoMsg
	msg.MsgType = "video"
	msg.Video = video
	return this.replyMsg(rw, &msg)
}

func (this *Weixinmp) ReplyMusicMsg(rw http.ResponseWriter, music *Music) error {
	var msg musicMsg
	msg.MsgType = "music"
	msg.Music = music
	return this.replyMsg(rw, &msg)
}

func (this *Weixinmp) ReplyNewsMsg(rw http.ResponseWriter, articles *[]Article) error {
	var msg newsMsg
	msg.MsgType = "news"
	msg.ArticleCount = len(*articles)
	msg.Articles.Item = articles
	return this.replyMsg(rw, &msg)
}

func (this *Weixinmp) replyMsg(rw http.ResponseWriter, msg interface{}) error {
	v := reflect.ValueOf(msg).Elem()
	v.FieldByName("ToUserName").SetString(this.Request.FromUserName)
	v.FieldByName("FromUserName").SetString(this.Request.ToUserName)
	v.FieldByName("CreateTime").SetInt(time.Now().Unix())
	data, err := xml.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := rw.Write(data); err != nil {
		return err
	}
	return nil
}

// send message methods
func (this *Weixinmp) SendTextMsg(touser string, content string) error {
	var msg textMsg
	msg.MsgType = "text"
	msg.Text.Content = content
	return this.sendMsg(touser, &msg)
}

func (this *Weixinmp) SendImageMsg(touser string, mediaId string) error {
	var msg imageMsg
	msg.MsgType = "image"
	msg.Image.MediaId = mediaId
	return this.sendMsg(touser, &msg)
}

func (this *Weixinmp) SendVoiceMsg(touser string, mediaId string) error {
	var msg voiceMsg
	msg.MsgType = "voice"
	msg.Voice.MediaId = mediaId
	return this.sendMsg(touser, &msg)
}

func (this *Weixinmp) SendVideoMsg(touser string, video *Video) error {
	var msg videoMsg
	msg.MsgType = "video"
	msg.Video = video
	return this.sendMsg(touser, &msg)
}

func (this *Weixinmp) SendMusicMsg(touser string, music *Music) error {
	var msg musicMsg
	msg.MsgType = "music"
	msg.Music = music
	return this.sendMsg(touser, &msg)
}

func (this *Weixinmp) SendNewsMsg(touser string, articles *[]Article) error {
	var msg newsMsg
	msg.MsgType = "news"
	msg.Articles.Item = articles
	return this.sendMsg(touser, &msg)
}

func (this *Weixinmp) sendMsg(touser string, msg interface{}) error {
	v := reflect.ValueOf(msg).Elem()
	v.FieldByName("ToUserName").SetString(touser)
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := this.post("message/custom/send", data); err != nil {
		return err
	}
	return nil
}

// get qrcode url
func (this *Weixinmp) GetQRCodeURL(ticket string) string {
	return "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + ticket
}

// create permanent qrcode
func (this *Weixinmp) CreateQRScene(sceneId int64) (string, error) {
	var inf qrScene
	inf.ActionName = "QR_SCENE"
	inf.ActionInfo.Scene.SceneId = sceneId
	return this.createQRCode(&inf)
}

// create temporary qrcode
func (this *Weixinmp) CreateQRLimitScene(expireSeconds, sceneId int64) (string, error) {
	var inf qrScene
	inf.ExpireSeconds = expireSeconds
	inf.ActionName = "QR_LIMIT_SCENE"
	inf.ActionInfo.Scene.SceneId = sceneId
	return this.createQRCode(&inf)
}

func (this *Weixinmp) createQRCode(inf *qrScene) (string, error) {
	data, err := json.Marshal(inf)
	if err != nil {
		return "", err
	}
	raw, err := this.post("qrcode/create", data)
	if err != nil {
		return "", err
	}
	var rtn struct {
		Ticket        string `json:"ticket"`
		ExpireSeconds int64  `json:"expire_seconds"`
	}
	if err := json.Unmarshal(raw, &rtn); err != nil {
		return "", err
	}
	return rtn.Ticket, nil
}

// send post request
func (this *Weixinmp) post(action string, data []byte) ([]byte, error) {
	// retry
	for i := 0; i < retryNum; i++ {
		token, err := this.accessToken.extract()
		if err != nil {
			if i < retryNum {
				continue
			}
			return nil, err
		}
		url := plainPreUrl + action + "?access_token=" + token
		resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(data))
		defer resp.Body.Close()
		if err != nil {
			if i < retryNum {
				continue
			}
			return nil, err
		}
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if i < retryNum {
				continue
			}
			return nil, err
		}
		var rtn struct {
			ErrCode int64  `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
		}
		if err := json.Unmarshal(raw, &rtn); err != nil {
			if i < retryNum {
				continue
			}
			return nil, err
		}
		// failed
		if rtn.ErrCode != 0 {
			if i < retryNum {
				continue
			}
			return nil, errors.New(rtn.ErrMsg)
		}
		return raw, nil
	}
	return nil, errors.New("send post request failed: " + action)
}
