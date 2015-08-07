package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Publisher struct {
	publisher string
	HookClient
}

func NewPublisher(host, apiVer, publisher, token string) *Publisher {
	return &Publisher{
		publisher: publisher,
		HookClient: HookClient{
			host:  host,
			ver:   apiVer,
			token: token,
		},
	}

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

func (p *Publisher) DeleteEvent(event string) (*Event, error) {
	url := p.api() + fmt.Sprintf("/events/%s/%s", p.publisher, event)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Access-Token", p.HookClient.token)
	client := http.Client{}
	rsp, err := client.Do(req)
	if nil != err || nil == rsp {
		return nil, err
	}

	retdata, _ := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(retdata))
	}

	e := Event{}
	err = json.Unmarshal(retdata, &e)
	return &e, err
}

func (p *Publisher) PublishEvent(event string, data interface{}) error {
	_, err := p.postData(p.api()+"/events/"+p.publisher+"/"+event, data)
	if nil != err {
		return err
	}
	return nil
}
