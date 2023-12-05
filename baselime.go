package baselime_opentelemetry

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/honeycombio/otel-config-go/otelconfig"
)

type Config struct {
	BaselimeApiKey string
	ServiceName    string
	Namespace      string
	CollectorUrl   string
	Protocol       string
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return fallback
}

func ConfigureOpenTelemetry(conf Config) (func(), error) {
	baselimeApiKey := getEnv("BASELIME_API_KEY", conf.BaselimeApiKey)

	if baselimeApiKey == "" {
		log.Printf("BASELIME_API_KEY not set, not configuring OpenTelemetry")
		return nil, nil
	}

	if conf.Protocol != "grpc" && conf.Protocol != "http" && conf.Protocol != "" {
		return nil, fmt.Errorf("protocol must be one of grpc, or http")
	}

	if conf.Protocol == "" {
		conf.Protocol = "http"
	}

	collectorUrl := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", conf.CollectorUrl)
	if collectorUrl == "" {
		if conf.Protocol == "grpc" {
			collectorUrl = "otel-grpc.baselime.io"
		}
		if conf.Protocol == "http" {
			collectorUrl = "https://otel.baselime.io"
		}
	}

	if conf.Protocol == "http" && strings.Contains(collectorUrl, "grpc") {
		return nil, fmt.Errorf("protocol is http but collector url is grpc")
	}

	if conf.Protocol == "grpc" && !strings.Contains(collectorUrl, "grpc") {
		return nil, fmt.Errorf("protocol is grpc but collector url is http")
	}

	serviceName := getEnv("OTEL_SERVICE_NAME", conf.ServiceName)

	otelConfProtocol := otelconfig.ProtocolGRPC

	if conf.Protocol == "http" {
		otelConfProtocol = otelconfig.ProtocolHTTPProto
	}

	log.Printf("Configuring OpenTelemetry with protocol %s, collector url %s, service name %s\n", conf.Protocol, collectorUrl, serviceName)

	otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithExporterProtocol(otelConfProtocol),
		otelconfig.WithExporterEndpoint(collectorUrl),
		otelconfig.WithServiceName(serviceName),
		otelconfig.WithHeaders(map[string]string{
			"x-api-key": baselimeApiKey,
		}),
	)

	if err != nil {
		return nil, err
	}

	return otelShutdown, nil
}
