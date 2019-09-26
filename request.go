package main

import (
	"encoding/json"
	"fmt"
	"github.com/murderxchip/cmap"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const agent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:68.0) Gecko/20100101 Firefox/68.0"

type RespData map[string]interface{}

type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    RespData `json:"data"`
}

type Request struct {
	config AccessConfig
}

var request Request

func NewRequest(config AccessConfig) *Request {
	return &Request{config: config}
}

func (req *Request) getUrl(route string, params cmap.CMap) string {
	if err := params.Set("_t", strconv.Itoa(int(time.Now().Unix()))); err != nil {
		panic(err)
	}

	param := make([]string, 5)
	for mapItem := range params.Dump() {
		param = append(param, fmt.Sprintf("%s=%s", mapItem.Key, mapItem.Value))
	}

	return fmt.Sprintf("%s/%s?%s", req.config.Host, route, strings.Join(param, "&"))
}

func ParseResponse(respBody string) (RespData, error) {
	response := Response{}
	err := json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		panic(err)
	}

	if response.Code != 0 {
		return nil, errors.New(response.Message)
	}

	return response.Data, nil
}

func (req *Request) Login() (token string, err error) {
	form := &LoginForm{req.config.Username, Md5(req.config.Password)}
	params := *cmap.NewCMap()
	url := req.getUrl("api/login", params)
	_, body, errs := gorequest.New().
		Post(url).
		Type("form").
		AppendHeader("Accept", "application/json").
		AppendHeader("User-Agent", agent).
		Send(fmt.Sprintf("username=%s&password=%s", form.Username, form.Password)).
		End()

	if errs != nil {
		panic("login failed")
	}

	//fmt.Println(body)
	respData, err := ParseResponse(body)
	if err != nil {
		return "", err
	}

	//respData
	logger.Println(respData)
	token = respData["token"].(string)
	return token, err
}
