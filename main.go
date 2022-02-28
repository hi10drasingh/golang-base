package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/droomlab/drm-coupon/pkg/global"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

var fileLogger zerolog.Logger

func main() {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(global.Global.Config.Log.Dir+"/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error().Err(err).Msg("error in opening log file")
	}
	defer file.Close()

	fileLogger := zerolog.New(file).With().Logger()

	fmt.Print(global.Global)

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/hello/:name", hello)

	fileLogger.Info().Msg("Server started at port 8888 testing")
	fileLogger.Fatal().Err(http.ListenAndServe(":8888", router))
}
