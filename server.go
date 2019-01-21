package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const URL_POST = "https://chat.googleapis.com/v1/spaces/XXXXXXXXXXX/messages?key=XXXX-XXXX&token=XXXXXX"
const TRHEADID = `{"name": "spaces/XXXXXX/threads/XXXXX"}`

type Message struct {
	Text string  `json:"text"`
	Thread string `json:"thread"`
}

type Grafana struct {
	Title string `json:title`
	RuleID int `json:ruleId`
	RuleUrl string `json:ruleUrl`
	State string `json:state`
	Message string `json:message`
}


func listenService(w http.ResponseWriter, r *http.Request) {

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var grafana Grafana
	err = json.Unmarshal(b, &grafana)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println(grafana)

	in:=[]byte(`{"text" : "*`+grafana.Title+`:* `+grafana.Message+` `+grafana.State+` `+grafana.RuleUrl+`","thread": `+TRHEADID+` }`)
	var raw map[string]interface{}
	json.Unmarshal(in, &raw)
	output, _ := json.Marshal(raw)
	//w.Header().Set("content-type", "application/json")
	//w.Write(output)
	sendAlertToGoogleChat(output,w)
}

func sendAlertToGoogleChat(in []byte,w http.ResponseWriter){
	log.Println(in)
	resp, err := http.Post(URL_POST, "application/json", bytes.NewBuffer(in))
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)
	log.Println(result["data"])
	
	w.Header().Set("Server", "A Go Web Server")
	w.WriteHeader(200)
}




func main() {
	http.HandleFunc("/sendToGchat", listenService)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
