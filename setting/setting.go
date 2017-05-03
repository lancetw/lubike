package setting

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config : dotenv config
type Config struct {
	UbikeEndpoint        string
	UbikeEndpointTimeout int
	MapquestAPIKey       string
}

var logFatalf = log.Fatalf
var gopath = os.Getenv("GOPATH")

// InitConfig : initialize dotenv config
// *TODO* remove gopath workaround
func InitConfig() Config {
	if gopath != "" {
		fullpath := gopath + "/src/github.com/lancetw/lubike/.env"
		err := godotenv.Load(fullpath)
		if err != nil {
			logFatalf("Error loading .env file %v", fullpath)
		}
	}

	var config Config
	config.UbikeEndpoint = os.Getenv("UBIKE_ENDPOINT")
	config.UbikeEndpointTimeout, _ = strconv.Atoi(os.Getenv("UBIKE_ENDPOINT_TIMEOUT"))
	config.MapquestAPIKey = os.Getenv("MAPQUEST_API_KEY")

	return config
}
