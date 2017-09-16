package main

import (
	"fmt"
	"net/http"

	"flag"

	"strconv"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/api"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/logger"
)

var (
	port = flag.Uint("port", 8080, "Port for API listener")
	host = flag.String("host", "0.0.0.0", "Host for API listener")
)

func main() {
	log := logger.Log.WithField("func", "main")

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", *host, strconv.Itoa(int(*port))), api.NewRouter())
	if err != nil {
		log.Error(err)
	}
}
