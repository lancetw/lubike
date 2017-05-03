package commonutil

import (
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
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
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
