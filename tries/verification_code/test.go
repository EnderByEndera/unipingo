package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	postData := url.Values{}
	postData.Add("content", "code:1234")
	postData.Add("phone_number", "18515195481")
	postData.Add("template_id", "CST_ptdie100")

	appcode := "99db8358267e4df383fdbbf1ba3bb6dc"

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://dfsns.market.alicloudapi.com/data/send_sms",
		strings.NewReader(postData.Encode()))
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	req.Header.Add("Authorization", "APPCODE "+appcode)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("fatal error", err)
		return
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	fmt.Println("content", string(content))
}
