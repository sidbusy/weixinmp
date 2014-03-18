package weixinmp

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// weixinmp request
type Request struct {
	// request header fields
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	// general request fields
	Content      string
	MsgId        int64
	PicUrl       string
	MediaId      string
	Format       string
	ThumbMediaId string
	LocationX    float64 `xml:"Location_X"`
	LocationY    float64 `xml:"Location_Y"`
	Scale        float64
	Label        string
	Title        string
	Description  string
	Url          string
	Recognition  string
	// event request fields
	Event     string
	EventKey  string
	Ticket    string
	Latitude  float64
	Longitude float64
	Precision float64

	token string
}

func (this *Request) IsValid(rw http.ResponseWriter, req *http.Request) bool {
	if !this.checkSignature(req) {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return false
	}
	if req.Method != "POST" {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(req.FormValue("echostr")))
		return false
	}
	if err := this.parseXML(req); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return false
	}
	return true
}

func (this *Request) checkSignature(req *http.Request) bool {
	ss := sort.StringSlice{
		this.token,
		req.FormValue("timestamp"),
		req.FormValue("nonce"),
	}
	sort.Strings(ss)          // sort strings by dictionary
	s := strings.Join(ss, "") // concatenate strings
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil)) == req.FormValue("signature")
}

func (this *Request) parseXML(req *http.Request) error {
	defer req.Body.Close()
	if raw, err := ioutil.ReadAll(req.Body); err != nil {
		return err
	} else if err := xml.Unmarshal(raw, this); err != nil {
		return err
	}
	return nil
}
