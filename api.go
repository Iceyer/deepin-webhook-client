package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ListEventAPI     = "/events"
	GetEventAPI      = "/event/%s/%s"
	GetSubscriberAPI = "/events/%s/%s/subscribers"
)

func do(req *http.Request, token string) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", token)

	client := http.Client{}
	rsp, err := client.Do(req)
	if nil != err || nil == rsp {
		return nil, err
	}

	retdata, _ := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(retdata))
	}

	return retdata, nil
}

type HookClient struct {
	token string
	host  string
	ver   string
}

func NewHookClient(host, ver, token string) *HookClient {
	return &HookClient{
		host:  host,
		ver:   ver,
		token: token,
	}
}

func (c *HookClient) putData(url string, jsondata []byte) ([]byte, error) {
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsondata))
	return do(req, c.token)
}

func (c *HookClient) deleteData(url string) ([]byte, error) {
	req, _ := http.NewRequest("DELETE", url, nil)
	return do(req, c.token)
}

func (c *HookClient) getData(url string) ([]byte, error) {
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	return do(req, c.token)
}

func (c *HookClient) postData(url string, data interface{}) ([]byte, error) {
	jdata, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jdata))
	return do(req, c.token)
}

func (p *HookClient) api() string {
	return p.host + "/" + p.ver
}

func (c *HookClient) GetSubscriber(publisher, event string) ([]*Subscriber, error) {
	ss := []*Subscriber{}
	data, err := c.getData(c.api() + fmt.Sprintf(GetSubscriberAPI, publisher, event))
	if nil != err {
		return ss, err
	}
	err = json.Unmarshal(data, &ss)
	return ss, err
}

func (c *HookClient) ListEvent() ([]*Event, error) {
	es := []*Event{}
	data, err := c.getData(c.api() + ListEventAPI)
	if nil != err {
		return es, err
	}
	err = json.Unmarshal(data, &es)
	return es, err
}
