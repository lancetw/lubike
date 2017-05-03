package commonutil

import (
	"errors"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lancetw/lubike/setting"
)

func TestPerror(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Perror did not panic")
		}
	}()

	Perror(errors.New("error"))
}

func TestFetchAPIData(t *testing.T) {
	config := setting.InitConfig()
	timeout := time.Duration(time.Duration(config.UbikeEndpointTimeout) * time.Second)
	json := FetchAPIData(config.UbikeEndpoint, timeout).(*simplejson.Json)
	collect, err := json.Get("retVal").Map()
	if err != nil {
		t.Error(err)
	}
	if len(collect) == 0 {
		t.Fail()
	}
}

func TestPostAPIData(t *testing.T) {
	config := setting.InitConfig()
	timeout := time.Duration(time.Duration(config.GoogleMapMatrixEndpointTimeout) * time.Second)
	values := map[string]interface{}{}
	json := PostAPIData(config.GoogleMapMatrixEndpoint, timeout, values).(*simplejson.Json)
	status := json.Get("status").MustString()
	if status != "INVALID_REQUEST" {
		t.Error(status)
	}
}
