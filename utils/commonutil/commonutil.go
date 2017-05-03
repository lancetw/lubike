package commonutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
)

// Perror : error and panic
func Perror(err error) {
	if err != nil {
		panic(err)
	}
}

// DecodeJSON : request body to a json map
func DecodeJSON(body []byte) interface{} {
	data, err := simplejson.NewJson(body)
	Perror(err)
	return data
}

// MakeRequest : request resources from a specific url
func MakeRequest(url string, timeout time.Duration) []byte {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	Perror(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Perror(err)

	return body
}

// MakeRequestPost : request resources from a specific url with post method
func MakeRequestPost(url string, timeout time.Duration, values interface{}) []byte {
	client := &http.Client{
		Timeout: timeout,
	}
	request, _ := json.Marshal(values)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(request))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	Perror(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Perror(err)
	return body
}

// FetchAPIData : fetch remote API to a data map
func FetchAPIData(url string, timeout time.Duration) interface{} {
	return DecodeJSON(MakeRequest(url, timeout))
}

// PostAPIData : post json data to remote API
func PostAPIData(url string, timeout time.Duration, values interface{}) interface{} {
	return DecodeJSON(MakeRequestPost(url, timeout, values))
}
