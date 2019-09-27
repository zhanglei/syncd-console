package main

import (
	"encoding/json"
	"fmt"
	"github.com/murderxchip/cmap"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const agent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:68.0) Gecko/20100101 Firefox/68.0"

type RespData interface{}

type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    RespData `json:"data"`
}

type Request struct {
	config AccessConfig
}

var request *Request

func NewRequest(config AccessConfig) *Request {
	if request != nil {
		return request
	}
	request = &Request{config: config}
	return request
}

func (req *Request) getUrl(route string, params cmap.CMap) string {
	if err := params.Set("_t", strconv.Itoa(int(time.Now().Unix()))); err != nil {
		panic(err)
	}

	param := make([]string, 1)
	for mapItem := range params.Dump() {
		if mapItem.Key != "" {
			param = append(param, fmt.Sprintf("%s=%s", mapItem.Key, mapItem.Value))
		}
	}

	url := fmt.Sprintf("%s://%s/%s?%s", req.config.Schema, req.config.Host, route, strings.Join(param, "&")[1:])
	logger.Println("url:", url)
	return url
}

func ParseResponse(respBody string) (RespData, error) {
	response := Response{}
	err := json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		panic(err)
	}

	if response.Code == 1005 {
		TokenFail()
	}

	if response.Code != 0 {
		return nil, errors.New(response.Message)
	}

	return response.Data, nil
}

func (req *Request) Login() {
	form := &LoginForm{req.config.Username, Md5(req.config.Password)}
	params := *cmap.NewCMap()
	url := req.getUrl("api/login", params)
	_, _, errs := gorequest.New().
		Post(url).
		Type("form").
		AppendHeader("Accept", "application/json").
		AppendHeader("User-Agent", agent).
		Send(fmt.Sprintf("username=%s&password=%s", form.Username, form.Password)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != 200 {
				panic(fmt.Sprintf("%s", errs))
			}

			respData, err := ParseResponse(body)
			if err != nil {
				panic(err)
			}

			respData1 := respData.(map[string]interface{})
			//respData
			SetToken(respData1["token"].(string))
		})

	if errs != nil {
		panic("登录失败，请设置正确的用户名和密码。")
	}
}

func (req *Request) AuthCookie() *http.Cookie {
	cookie := http.Cookie{}
	cookie.Name = "_syd_identity"
	cookie.Value = GetToken()

	return &cookie
}

/**
http://deploy.tech.mofaxiao.com/api/deploy/apply/project/all?_t=1568861966520
*/
func (req *Request) Projects() (token string, err error) {
	params := *cmap.NewCMap()
	url := req.getUrl("api/deploy/apply/project/all", params)
	_, body, errs := gorequest.New().
		Get(url).
		AppendHeader("Accept", "application/json").
		AppendHeader("Host", req.config.Host).
		AppendHeader("User-Agent", agent).
		AddCookie(req.AuthCookie()).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != 200 {
				panic(fmt.Sprintf("%s", errs))
			}
		})

	if errs != nil {
		fmt.Println(errs)
	}

	fmt.Println(body)
	respData, err := ParseResponse(body)
	if err != nil {
		return "", err
	}

	//respData
	logger.Printf("%v", respData)
	return "", err
}
