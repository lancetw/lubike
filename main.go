package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/lancetw/lubike/utils/ubikeutil"
	e "github.com/lonnblad/negroni-etag/etag"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var r = render.New()

func lubikeCommonEndpoint(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		lat := req.FormValue("lat")
		lng := req.FormValue("lng")

		num := 2
		result, errno := ubikeutil.LoadNearbyUbikes(lat, lng, num)

		r.JSON(w, http.StatusOK, map[string]interface{}{
			"code":   errno,
			"result": result,
		})
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//var err error
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.Use(gzip.Gzip(gzip.DefaultCompression))

	v1outes := mux.NewRouter().PathPrefix("/v1").Subrouter().StrictSlash(true)
	v1outes.HandleFunc("/ubike-station/taipei", lubikeCommonEndpoint)

	n.Use(c)
	n.UseHandler(v1outes)
	n.Use(e.Etag())
	n.Run(":" + port)
}
