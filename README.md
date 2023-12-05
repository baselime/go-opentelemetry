# Go Baselime OpenTelemetry SDK
[![Documentation][docs_badge]][docs]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]

Instrument your Go applications with OpenTelemetry and send the traces to [Baselime](https://baselime.io)

## Getting Started 

Check out the [documentation](https://baselime.io/docs/sending-data/opentelemetry/).

## Example

```go
package main

import (
	"context"
	"log"

	"github.com/baselime/go-opentelemetry"
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

```

## License

&copy; Baselime Limited, 2023

Distributed under MIT License (`The MIT License`).

See [LICENSE](LICENSE) for more information.

<!-- Badges -->

[docs]: https://baselime.io/docs/
[docs_badge]: https://img.shields.io/badge/docs-reference-blue.svg?style=flat-square
[release]: https://github.com/baselime/go-opentelemetry/releases/latest
[release_badge]: https://img.shields.io/github/release/baselime/go-opentelemetry.svg?style=flat-square&ghcache=unused
[license]: https://opensource.org/licenses/MIT
[license_badge]: https://img.shields.io/github/license/baselime/go-opentelemetry.svg?color=blue&style=flat-square&ghcache=unused