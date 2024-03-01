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
			collectorUrl = "otel-ingest.baselime.io"
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
		otelconfig.WithResourceAttributes(getResourceAttributes()),
	)

	if err != nil {
		return nil, err
	}

	return otelShutdown, nil
}

func getResourceAttributes() map[string]string {
	resAttr := map[string]string{}
	var attrMapping map[string]string
	if _, ok := os.LookupEnv("FLY_APP_NAME"); ok {
		resAttr["cloud.provider"] = "fly"
		attrMapping = map[string]string{
			"FLY_APP_NAME":        "fly.app.name",
			"FLY_MACHINE_ID":      "fly.machine.id",
			"FLY_ALLOC_ID":        "fly.allocId",
			"FLY_REGION":          "fly.region",
			"FLY_PUBLIC_IP":       "fly.publicIp",
			"FLY_IMAGE_REF":       "fly.imageRef",
			"FLY_MACHINE_VERSION": "fly.machine.version",
			"FLY_PRIVATE_IP":      "fly.privateIp",
			"FLY_PROCESS_GROUP":   "fly.processGroup",
			"FLY_VM_MEMORY_MB":    "fly.machine.vmMemoryMb",
			"PRIMARY_REGION":      "fly.primaryRegion",
		}
	} else if _, ok := os.LookupEnv("KOYEB_APP_NAME"); ok {
		resAttr["cloud.provider"] = "koyeb"
		attrMapping = map[string]string{
			"KOYEB_APP_NAME":               "koyeb.app.name",
			"KOYEB_APP_ID":                 "koyeb.app.id",
			"KOYEB_ORGANIZATION_NAME":      "koyeb.organization.name",
			"KOYEB_ORGANIZATION_ID":        "koyeb.organization.id",
			"KOYEB_SERVICE_NAME":           "koyeb.service.name",
			"KOYEB_SERVICE_ID":             "koyeb.service.id",
			"KOYEB_SERVICE_PRIVATE_DOMAIN": "koyeb.service.privateDomain",
			"KOYEB_PUBLIC_DOMAIN":          "koyeb.publicDomain",
			"KOYEB_REGION":                 "koyeb.region",
			"KOYEB_REGIONAL_DEPLOYMENT_ID": "koyeb.regionalDeploymentId",
			"KOYEB_INSTANCE_ID":            "koyeb.instance.id",
			"KOYEB_INSTANCE_TYPE":          "koyeb.instance.type",
			"KOYEB_INSTANCE_MEMORY_MB":     "koyeb.instance.memory",
			"KOYEB_PRIVILEGED":             "koyeb.privileged",
			"KOYEB_HYPERVISOR_ID":          "koyeb.hypervisor.id",
			"KOYEB_DC":                     "koyeb.dc",
			"KOYEB_DOCKER_REF":             "koyeb.docker.ref",
			"KOYEB_GIT_SHA":                "koyeb.git.sha",
			"KOYEB_GIT_BRANCH":             "koyeb.git.branch",
			"KOYEB_GIT_COMMIT_AUTHOR":      "koyeb.git.commit.author",
			"KOYEB_GIT_COMMIT_MESSAGE":     "koyeb.git.commit.message",
			"KOYEB_GIT_REPOSITORY":         "koyeb.git.repository",
		}
	} else if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		resAttr["cloud.provider"] = "aws"
		attrMapping = map[string]string{
			"AWS_REGION":                      "aws.region",
			"AWS_DEFAULT_REGION":              "aws.defaultRegion",
			"AWS_EXECUTION_ENV":               "aws.lambda.executionEnv",
			"AWS_LAMBDA_FUNCTION_NAME":        "aws.lambda.functionName",
			"AWS_LAMBDA_FUNCTION_VERSION":     "aws.lambda.functionVersion",
			"AWS_LAMBDA_FUNCTION_MEMORY_SIZE": "aws.lambda.functionMemorySize",
			"AWS_LAMBDA_INITIALIZATION_TYPE":  "aws.lambda.initializationType",
			"AWS_LAMBDA_RUNTIME_API":          "aws.lambda.runtimeApi",
		}
	} else if _, ok := os.LookupEnv("VERCEL"); ok {
		resAttr["cloud.provider"] = "vercel"
		attrMapping = map[string]string{
			"VERCEL_REGION":                 "vercel.region",
			"VERCEL_ENV":                    "vercel.env",
			"VERCEL_URL":                    "vercel.url",
			"VERCEL_BRANCH_URL":             "vercel.branch.url",
			"VERCEL_GIT_PROVIDER":           "vercel.git.provider",
			"VERCEL_GIT_REPO_SLUG":          "vercel.git.repoSlug",
			"VERCEL_GIT_COMMIT_SHA":         "vercel.git.commitSha",
			"VERCEL_GIT_COMMIT_MESSAGE":     "vercel.git.commitMessage",
			"VERCEL_GIT_COMMIT_AUTHOR_NAME": "vercel.git.commitAuthorName",
		}
	} else {
		attrMapping = map[string]string{}
	}
	for envName, attrName := range attrMapping {
		if value, ok := os.LookupEnv(envName); ok {
			resAttr[attrName] = value
		}
	}
	return resAttr
}
