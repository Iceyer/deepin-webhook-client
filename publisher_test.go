package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	sbody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Recive Msg:", string(sbody))
}

func init() {
	http.HandleFunc("/", viewHandler)
	go http.ListenAndServe("0.0.0.0:8890", nil)
}

func TestPublisher(t *testing.T) {
	serverHost := "http://hook.deepin.test"

	HookClient := HookClient{
		host:  serverHost,
		ver:   "v0",
		token: "OWNjN2QyMTMtMWRmNC00N",
	}
	es, err := HookClient.ListEvent()
	t.Log(err, es)
	return
	//Create Event
	p := NewPublisher(serverHost, "v0", "repo", "OWNjN2QyMTMtMWRmNC00N")
	_, err = p.CreateEvent("app_create", "123456", nil)
	t.Log(err)

	//Subscribe Event
	s := NewSubscriber(serverHost, "v0", "OWNjN2QyMTMtMWRmNC00N")
	revHost := "http://192.168.11.111:8890/"
	fmt.Println("Warning: Modify revHost to your recive hook server host")
	_, err = s.Subscribe("repo", "app_create", revHost)
	t.Log(err)

	//Publish Event
	p.PublishEvent("app_create", "testdata")

	//Wait Event
	time.Sleep(time.Second * 5)
}
