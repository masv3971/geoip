module geoip

go 1.17

require (
	github.com/biter777/countries v1.3.4
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-redis/redis/v8 v8.11.4
	github.com/google/uuid v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe
	github.com/moogar0880/problems v0.1.1
	github.com/oschwald/geoip2-golang v1.6.1
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/mongo-driver v1.8.3
	go.opentelemetry.io/otel/trace v1.4.1
	go.uber.org/zap v1.21.0
	gopkg.in/yaml.v2 v2.4.0
	inet.af/netaddr v0.0.0-20211027220019-c74959edd3b6
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.14.4 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/oschwald/maxminddb-golang v1.8.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.0 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.opentelemetry.io/otel v1.4.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go4.org/intern v0.0.0-20211027215823-ae77deb06f29 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20211027215541-db492cf91b37 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

//replace github.com/tkanos/go-dtree => github.com/masv3971/go-dtree v0.5.1-0.20220215152355-4680398bc5f1
replace github.com/tkanos/go-dtree => ../go-dtree
