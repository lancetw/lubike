package setting

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config : dotenv config
type Config struct {
	UbikeEndpoint                      string
	UbikeEndpointTimeout               int
	MapquestAPIKey                     string
	MapquestRouteMatrixEndpoint        string
	MapquestRouteMatrixEndpointTimeout int
	GoogleMapMatrixAPIKey              string
	GoogleMapMatrixEndpoint            string
	GoogleMapMatrixEndpointTimeout     int
}

var logFatalf = log.Fatalf
var gopath = os.Getenv("GOPATH")
var testingMode = os.Getenv("TESTING")
var godotenvLoad = godotenv.Load

// InitConfig : initialize dotenv config
func InitConfig() Config {
	if !(testingMode == "true") && gopath != "" {
		fullpath := gopath + "/src/github.com/lancetw/lubike/.env"
		err := godotenvLoad(fullpath)
		if err != nil {
			logFatalf("Error loading .env file %v", fullpath)
		}
	}

	var config Config
	config.UbikeEndpoint = os.Getenv("UBIKE_ENDPOINT")
	config.UbikeEndpointTimeout, _ = strconv.Atoi(os.Getenv("UBIKE_ENDPOINT_TIMEOUT"))
	config.MapquestAPIKey = os.Getenv("MAPQUEST_API_KEY")
	config.MapquestRouteMatrixEndpoint = os.Getenv("MAPQUEST_ROUTE_MATRIX_ENDPOINT")
	config.MapquestRouteMatrixEndpointTimeout, _ = strconv.Atoi(os.Getenv("MAPQUEST_ROUTE_MATRIX_ENDPOINT_TIMEOUT"))
	config.GoogleMapMatrixAPIKey = os.Getenv("GOOGLEMAP_MATRIX_API_KEY")
	config.GoogleMapMatrixEndpoint = os.Getenv("GOOGLEMAP_MATRIX_ENDPOINT")
	config.GoogleMapMatrixEndpointTimeout, _ = strconv.Atoi(os.Getenv("GOOGLEMAP_MATRIX_API_ENDPOINT_TIMEOUT"))

	return config
}
