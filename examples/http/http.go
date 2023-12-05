package main

import (
	baselime_opentelemetry "baselime-opentelemetry"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("hello-tracer")

// Implement an HTTP Handler function to be instrumented
func httpHandler(w http.ResponseWriter, r *http.Request) {

	_, span := tracer.Start(r.Context(), "hello-span")
	defer span.End()

	span.SetAttributes(
		attribute.String("foo", "bar"),
		attribute.Bool("fizz", true),
	)

	fmt.Fprintf(w, "Hello, World")
}

// Wrap the HTTP handler function with OTel HTTP instrumentation
func wrapHandler() {
	handler := http.HandlerFunc(httpHandler)
	wrappedHandler := otelhttp.NewHandler(handler, "hello")
	http.Handle("/hello", wrappedHandler)
}

func main() {

	params := baselime_opentelemetry.Config{
		ServiceName: "hello-api",
	}
	otelShutdown, err := baselime_opentelemetry.ConfigureOpenTelemetry(params)

	if err != nil {
		log.Fatalf("error setting up OTel SDK - %e", err)
	}

	defer otelShutdown()

	wrapHandler()
	log.Fatal(http.ListenAndServe(":3030", nil))
}
