package main

import (
	baselime_opentelemetry "baselime-opentelemetry"
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("hello-tracer")

func main() {
	params := baselime_opentelemetry.Config{
		ServiceName: "hello-basic",
	}
	otelShutdown, err := baselime_opentelemetry.ConfigureOpenTelemetry(params)

	if err != nil {
		log.Fatalf("error setting up OTel SDK - %e", err)
	}

	defer otelShutdown()
	ctx := context.TODO()
	_, span := tracer.Start(ctx, "hello-span")
	defer span.End()

	span.SetAttributes(
		attribute.String("foo", "bar"),
		attribute.Bool("fizz", true),
	)
}
