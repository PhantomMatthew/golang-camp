package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func ReqWithAuth(method, url string, params string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(params))

	//if len(params) > 0 {
	//	param := params[0]
	//	req, err := http.NewRequest(method, url, strings.NewReader(param))
	//} else {
	//
	//}

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibWljX3dlY2hhdCIsInBvcnQiOjMwMDQsImlhdCI6MTU2NjQ0NjQ4OSwiZXhwIjo0NzIyMjA2NDg5fQ.MqAq3BZu0n2hCRMs1hw0Zd6DPSRd0yHzoXJTFhJuoM0")
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
