package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"syncd-console/resp"
	"time"
)

var UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:68.0) Gecko/20100101 Firefox/68.0"

func main() {
	//InitSyncdConfig()

	form := &LoginForm{"qinming", Md5("123456")}
	//formParams, errs := json.Marshal(form)
	//
	//if errs != nil {
	//	panic("form param error")
	//}
	url := fmt.Sprintf("%s%d", "http://deploy.tech.mofaxiao.com/api/login?_t=", time.Now().Unix())
	_, body, errs := gorequest.New().
		Post(url).
		Type("form").
		AppendHeader("Accept", "application/json").
		AppendHeader("User-Agent", UserAgent).
		Send(fmt.Sprintf("username=%s&password=%s", form.Username, form.Password)).
		End()

	if errs != nil {
		panic("login failed")
	}

	respData := resp.DataResponse{}
	//fmt.Print(body)
	err := json.Unmarshal([]byte(body), &respData)
	if err != nil {
		panic(err)
	}

	if respData.Code != 0 {
		fmt.Println(respData.Message)
		return
	}

	fmt.Printf("%++v%v%v", respData.Code, respData.Message, respData.Data)
}
