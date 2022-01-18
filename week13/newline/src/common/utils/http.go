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

	req.Header.Set("Authorization", "Bearer ")
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
