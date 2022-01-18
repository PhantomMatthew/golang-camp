package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Alert(content string){

	client := &http.Client{}
	robotSdkUrl:="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=64b69d58-896a-420a-a647-78d8467f9930"

	if(len(content)>2048){
		content = content[0:2000];
	}
	message:= QWRobotMessage{
		Msgtype: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: content,
		},
	}
	data, err := json.Marshal(message)
	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", robotSdkUrl, body)
	if err != nil {
		// handle error
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	fmt.Println(string(rspBody))
	jsonStr := string(rspBody)
	dynamic := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &dynamic)

}
type QWRobotMessage struct{
	Msgtype string  `json:"msgtype"`
	Text struct{
		Content string  `json:"content"`
	} `json:"text"`
}