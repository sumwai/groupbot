package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	. "net/url"
)

const API_HOST = "10.0.0.130"
const API_PORT = 5700
const ChangeLinkApi = "http://test.sixgm.cn/index/api/up_qun?url="

// Call 调用bot接口
//
// api: 调用的API端点
//
// params: 调用API需要的参数
//
func Call(api string, params map[string]string) string {
	url := fmt.Sprintf("http://%s:%d/%s?", API_HOST, API_PORT, api)

	for k, v := range params {
		url += k + "=" + QueryEscape(v)
		url += "&"
	}

	log.Println()
	log.Printf("API: {%s}\n", api)
	log.Println("PARAM: ")

	for k, v := range params {
		log.Printf("\t%s:%s\n", k, v)
	}
	log.Println()

	response, err := http.Get(url)
	if err != nil {
		log.Println("Http请求失败", err)
		return ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("close http body error", err)
		}
	}(response.Body)

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("读取响应结果失败")
		return ""
	}

	return string(bytes)
}

type ChangeResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

// ChangeGroupLink 修改群链接
//
// link: 新的群链接
//
func ChangeGroupLink(link string) (msg string) {
	log.Println("change join group link")
	log.Printf("[Link]: %s\n", link)
	link = QueryEscape(link)
	url := fmt.Sprintf("%s%s", ChangeLinkApi, link)
	response, err := http.Get(url)
	if err != nil {
		log.Println("change group link: http get error", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("change group link: close response body error", err)
		}
	}(response.Body)

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("change group link: read response body error", err)
		return
	}

	var result *ChangeResponse
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Println("[Error]: change group link, response error, unable to unmarshal json string")
		return
	}

	log.Printf("[Change]: %s", result.Msg)

	return result.Msg
}
