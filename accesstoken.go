package weixinmp

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type AccessToken struct {
	AppId     string
	AppSecret string
	TmpName   string
}

// get fresh access_token string
func (this *AccessToken) Fresh() (string, error) {
	if this.TmpName == "" {
		this.TmpName = "accesstoken.tmp"
	}
	fi, err := os.Stat(this.TmpName)
	if err != nil && !os.IsExist(err) {
		return this.fetchAndStore()
	}
	expires := fi.ModTime().Unix() + 7200
	if expires <= time.Now().Unix() {
		return this.fetchAndStore()
	}
	tmp, err := os.Open(this.TmpName)
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	data, err := ioutil.ReadAll(tmp)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (this *AccessToken) fetchAndStore() (string, error) {
	token, err := this.fetch()
	if err != nil {
		return "", err
	}
	if err := this.store(token); err != nil {
		return "", err
	}
	return token, nil
}

func (this *AccessToken) store(token string) error {
	tmp, err := os.OpenFile(this.TmpName, os.O_WRONLY|os.O_CREATE, os.ModeTemporary)
	if err != nil {
		return err
	}
	defer tmp.Close()
	if _, err := tmp.Write([]byte(token)); err != nil {
		return err
	}
	return nil
}

func (this *AccessToken) fetch() (string, error) {
	rtn, err := get(fmt.Sprintf(
		"%stoken?grant_type=client_credential&appid=%s&secret=%s",
		UrlPrefix,
		this.AppId,
		this.AppSecret,
	))
	if err != nil {
		return "", err
	}
	return rtn.AccessToken, nil
}
