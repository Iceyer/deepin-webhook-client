package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type hookClient struct {
	token string
	host  string
	ver   string
}

type Publisher struct {
	publisher string
	hookClient
}

func NewPublisher(host, apiVer, publisher, token string) *Publisher {
	return &Publisher{
		publisher: publisher,
		hookClient: hookClient{
			host:  host,
			ver:   apiVer,
			token: token,
		},
	}

}

func (c *hookClient) postData(url string, data interface{}) ([]byte, error) {
	jdata, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jdata))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", c.token)

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

func (p *hookClient) api() string {
	return p.host + "/" + p.ver
}

func (p *Publisher) CreateEvent(event, secret string, schema map[string]interface{}) (*Event, error) {
	e := Event{
		Name:      event,
		Secret:    secret,
		Schema:    schema,
		Publisher: p.publisher,
	}
	retData, err := p.postData(p.api()+"/events", e)
	if nil != err {
		return nil, err
	}
	err = json.Unmarshal(retData, &e)
	return &e, err
}

func (p *Publisher) PublishEvent(event string, data interface{}) error {
	_, err := p.postData(p.api()+"/events/"+p.publisher+"/"+event, data)
	if nil != err {
		return err
	}
	return nil
}
