module go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp

go 1.20

replace go.opentelemetry.io/contrib => github.com/coderpoet/opentelemetry-go-contrib v1.3.1-0.20231226171440-1a616f130d95

require (
	github.com/felixge/httpsnoop v1.0.4
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/contrib v1.3.1
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/metric v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
