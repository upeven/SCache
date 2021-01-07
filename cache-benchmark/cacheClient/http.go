package cacheClient

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//http客户端实现

type httpClient struct {
	*http.Client
	server string
}

func (c *httpClient) get(key string) string {
	resp, e := c.Get(c.server + key)
	if e != nil {
		log.Println(key + "unknown error")
		//panic(e)
	}
	if resp.StatusCode == http.StatusNotFound {
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println("未读取到数据")
	}
	return string(b)
}

func (c *httpClient) set(key, value string) {
	//发送请求
	req, e := http.NewRequest(http.MethodPut,
		c.server+key, strings.NewReader(value))
	if e != nil {
		log.Println(key + "设置缓存失败")
		//panic(e)
	}
	resp, e := c.Do(req)
	if e != nil {
		log.Println(key + "设置缓存失败")
		//panic(e)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
	}
}

//实现client接口
func (c *httpClient) Run(cmd *Cmd) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if cmd.Name == "get" {
		cmd.Value = c.get(cmd.Key)
		return
	}
	if cmd.Name == "set" {
		c.set(cmd.Key, cmd.Value)
		return
	}
	panic("unknown cmd name " + cmd.Name)
}

func newHTTPClient(server string) *httpClient {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, "http://" + server + ":12345/cache/"}
}

//http不支持PiplineRun
func (c *httpClient) PipelinedRun([]*Cmd) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	panic("httpClient pipelined run not implement")
}
