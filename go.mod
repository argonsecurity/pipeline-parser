module github.com/argonsecurity/pipeline-parser

go 1.21

toolchain go1.21.9

require (
	github.com/go-test/deep v1.0.8
	github.com/imroc/req/v3 v3.42.3
	github.com/mitchellh/mapstructure v1.4.3
	github.com/pkg/errors v0.9.1
	github.com/r3labs/diff/v3 v3.0.0
	github.com/spf13/cobra v1.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/andybalholm/brotli v1.0.6 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/pprof v0.0.0-20231229205709-960ae82b1e42 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/onsi/ginkgo/v2 v2.13.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/quic-go v0.40.1 // indirect
	github.com/refraction-networking/utls v1.6.0 // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	golang.org/x/exp v0.0.0-20240103183307-be819d1f06fc
	golang.org/x/net v0.23.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
)

// fix CVE-2024-22189
replace github.com/quic-go/quic-go v0.40.1 => github.com/quic-go/quic-go v0.42.0
