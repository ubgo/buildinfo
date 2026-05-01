module github.com/ubgo/buildinfo/contrib/buildinfo-otel

go 1.24

require (
	github.com/ubgo/buildinfo v0.0.0
	go.opentelemetry.io/otel v1.32.0
)

require github.com/stretchr/testify v1.10.0 // indirect

replace github.com/ubgo/buildinfo => ../..
