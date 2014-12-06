package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	nlogrus "github.com/meatballhat/negroni-logrus"
)

var Log *logrus.Logger

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	logger := nlogrus.NewMiddleware()
	Log = logger.Logger

	n := negroni.Classic()
	n.Use(logger)
	n.UseHandler(mux)
	n.Run(":8000")
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	invoice := Invoice{}
	invoice.OrderNumber = 1920

	Log.Printf("%#v\n", invoice.GetLineItems())
}
