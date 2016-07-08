package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<a href="/home"> Click here to start a request </a>`))
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Request started"))
	sp := opentracing.StartSpan("GET /home") // Start a new root span.
	defer sp.Finish()

	asyncReq, _ := http.NewRequest("GET", "http://localhost:8080/async", nil)
	// Inject the trace information into the HTTP Headers.
	err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeaderTextMapCarrier(asyncReq.Header))
	if err != nil {
		log.Fatalf("%s: Couldn't inject headers (%v)", r.URL.Path, err)
	}

	go func() {
		sleepMilli(50)
		if _, err := http.DefaultClient.Do(asyncReq); err != nil {
			log.Printf("%s: Async call failed (%v)", r.URL.Path, err)
		}
	}()

	sleepMilli(10)
	syncReq, _ := http.NewRequest("GET", "http://localhost:8080/service", nil)
	// Inject the trace info into the headers.
	err = sp.Tracer().Inject(sp.Context(),
		opentracing.TextMap,
		opentracing.HTTPHeaderTextMapCarrier(syncReq.Header))
	if err != nil {
		log.Fatalf("%s: Couldn't inject headers (%v)", r.URL.Path, err)
	}
	if _, err = http.DefaultClient.Do(syncReq); err != nil {
		log.Printf("%s: Synchronous call failed (%v)", r.URL.Path, err)
		return
	}
	w.Write([]byte("... done!"))
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	opName := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeaderTextMapCarrier(r.Header))
	if err == nil {
		sp = opentracing.StartSpan(opName, opentracing.ChildOf(spCtx))
	} else {
		sp = opentracing.StartSpan(opName)
	}
	defer sp.Finish()

	sleepMilli(50)

	dbReq, _ := http.NewRequest("GET", "http://localhost:8080/db", nil)
	err = sp.Tracer().Inject(sp.Context(),
		opentracing.TextMap,
		opentracing.HTTPHeaderTextMapCarrier(dbReq.Header))
	if err != nil {
		log.Fatalf("%s: Couldn't inject headers (%v)", r.URL.Path, err)
	}

	if _, err := http.DefaultClient.Do(dbReq); err != nil {
		sp.LogEventWithPayload("db request error", err)
	}
}

func dbHandler(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span

	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeaderTextMapCarrier(r.Header))
	if err != nil {
		log.Println("%s: Could not join trace (%v)", r.URL.Path, err)
		return
	}
	if err == nil {
		sp = opentracing.StartSpan("GET /db", opentracing.ChildOf(spanCtx))
	} else {
		sp = opentracing.StartSpan("GET /db")
	}
	defer sp.Finish()
	sleepMilli(25)
}

func sleepMilli(min int) {
	time.Sleep(time.Millisecond * time.Duration(min+rand.Intn(100)))
}
