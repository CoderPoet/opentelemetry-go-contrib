module go.opentelemetry.io/contrib/instrumentation/github.com/apache/rocketmq/otelrocketmq

go 1.16

replace go.opentelemetry.io/contrib => ../../../../..

require (
	github.com/apache/rocketmq-client-go/v2 v2.1.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/trace v1.3.0
)

replace github.com/apache/rocketmq-client-go/v2 => /Users/bytedance/opensource/rocketmq-client-go
