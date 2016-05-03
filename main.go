package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	lightstepot "github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"

	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
)

var (
	port           = flag.Int("port", 8080, "Example app port.")
	appdashPort    = flag.Int("appdash.port", 8700, "Run appdash locally on this port.")
	lightstepToken = flag.String("lighstep.token", "", "Lightstep access token.")
)

func main() {
	flag.Parse()

	var tracer opentracing.Tracer

	// Would it make sense to embed Appdash?
	if len(*lightstepToken) > 0 {
		tracer = lightstepot.NewTracer(lightstepot.Options{AccessToken: *lightstepToken})
	} else {
		addr := startAppdashServer(*appdashPort)
		tracer = appdashot.NewTracer(appdash.NewRemoteCollector(addr))
	}

	opentracing.InitGlobalTracer(tracer)

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/home", homeHandler)
	mux.HandleFunc("/async", serviceHandler)
	mux.HandleFunc("/service", serviceHandler)
	mux.HandleFunc("/db", dbHandler)
	fmt.Printf("Go to http://localhost:%d/home to start a request!\n", *port)
	log.Fatal(http.ListenAndServe(addr, mux))
}
