package weixinmp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type accessToken struct {
	appid  string
	secret string
}

func (this *accessToken) extract() (string, error) {
	fi, err := os.Stat(accessTokenTemp)
	if err != nil && !os.IsExist(err) {
		return this.fetchAndStore()
	}
	expires := fi.ModTime().Unix() + 7200
	if expires <= time.Now().Unix() {
		return this.fetchAndStore()
	}
	temp, err := os.OpenFile(accessTokenTemp, os.O_RDONLY, os.ModeTemporary)
	if err != nil {
		return "", err
	}
	defer temp.Close()
	raw, err := ioutil.ReadAll(temp)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (this *accessToken) fetchAndStore() (string, error) {
	token, err := this.fetch()
	if err != nil {
		return "", err
	}
	if err := this.store(token); err != nil {
		return "", err
	}
	return token, nil
}

func (this *accessToken) fetch() (string, error) {
	qs := fmt.Sprintf(
		"?grant_type=client_credential&appid=%s&secret=%s",
		this.appid,
		this.secret,
	)
	url := plainPreUrl + "token" + qs
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var rtn struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		ErrCode     int64  `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.Unmarshal(raw, &rtn); err != nil {
		return "", err
	}
	if rtn.ErrCode != 0 {
		return "", errors.New(fmt.Sprintf("%d %s", rtn.ErrCode, rtn.ErrMsg))
	}
	return rtn.AccessToken, nil
}

func (this *accessToken) store(token string) error {
	temp, err := os.OpenFile(accessTokenTemp, os.O_WRONLY|os.O_CREATE, os.ModeTemporary)
	if err != nil {
		return err
	}
	defer temp.Close()
	if _, err := temp.Write([]byte(token)); err != nil {
		return err
	}
	return nil
}
