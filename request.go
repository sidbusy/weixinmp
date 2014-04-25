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

// request from weixinmp
type Request struct {
	Token string
	// request common fields
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	// message request fields
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
}

// validate request
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
	if err := this.parseRequest(req); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return false
	}
	return true
}

func (this *Request) parseRequest(req *http.Request) error {
	raw, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if err := xml.Unmarshal(raw, this); err != nil {
		return err
	}
	return nil
}

func (this *Request) checkSignature(req *http.Request) bool {
	ss := sort.StringSlice{
		this.Token,
		req.FormValue("timestamp"),
		req.FormValue("nonce"),
	}
	sort.Strings(ss)          // sort strings by dictionary
	s := strings.Join(ss, "") // concatenate strings
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil)) == req.FormValue("signature")
}
