module cto-github.cisco.com/NFV-BU/go-msx

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.0.5
	github.com/AlekSi/gocov-xml v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/Shopify/sarama v1.26.1
	github.com/ThreeDotsLabs/watermill v1.0.2
	github.com/ThreeDotsLabs/watermill-kafka/v2 v2.2.0
	github.com/ThreeDotsLabs/watermill-sql v1.3.4
	github.com/armon/go-metrics v0.0.0-20190430140413-ec5e00d3c878 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535 // indirect
	github.com/axw/gocov v1.1.0 // indirect
	github.com/benbjohnson/clock v1.0.3
	github.com/bmatcuk/doublestar v1.1.5
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/dave/jennifer v1.4.1
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/doug-martin/goqu/v9 v9.9.0
	github.com/elastic/go-seccomp-bpf v1.2.0
	github.com/emicklei/go-restful v2.14.3+incompatible
	github.com/emicklei/go-restful-openapi v1.2.0
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813
	github.com/getkin/kin-openapi v0.20.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-ini/ini v1.48.0
	github.com/go-openapi/spec v0.20.3
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-redis/redis/v7 v7.4.1
	github.com/gocql/gocql v0.0.0-20191013011951-93ce931da9e1
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/uuid v1.1.1
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/hashicorp/consul v1.4.0
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/hashicorp/memberlist v0.1.4 // indirect
	github.com/hashicorp/serf v0.0.0-20161207011743-d3a67ab21bc8 // indirect
	github.com/hashicorp/vault/api v1.0.4
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/jackc/puddle v1.2.1
	github.com/jackpal/gateway v1.0.7
	github.com/jmoiron/sqlx v1.2.0
	github.com/jstemmer/go-junit-report v1.0.0 // indirect
	github.com/kennygrant/sanitize v1.2.4
	github.com/lib/pq v1.3.0
	github.com/magiconair/properties v1.8.1
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mitchellh/mapstructure v1.4.2
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/opentracing/opentracing-go v1.1.0
	github.com/otiai10/copy v1.0.2
	github.com/pavel-v-chernykh/keystore-go v2.1.0+incompatible
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang v1.2.1
	github.com/radovskyb/watcher v1.0.7
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/robfig/cron/v3 v3.0.1
	github.com/scylladb/go-reflectx v1.0.1
	github.com/scylladb/gocqlx v1.3.1
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.7.0
	github.com/smartystreets/assertions v0.0.0-20190116191733-b6c0e53d7304 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.4
	github.com/thejerf/abtime v1.0.3
	github.com/tidwall/gjson v1.6.4
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/uber/jaeger-client-go v2.19.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible
	go.uber.org/atomic v1.6.0
	golang.org/x/mod v0.5.1
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/tools v0.1.0 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.33.0
	gopkg.in/ini.v1 v1.41.0 // indirect
	gopkg.in/pipe.v2 v2.0.0-20140414041502-3c2ca4d52544
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/gotestsum v1.8.1 // indirect
	moul.io/banner v1.0.1
)

replace (
	github.com/ThreeDotsLabs/watermill-kafka/v2 v2.2.0 => cto-github.cisco.com/NFV-BU/watermill-kafka/v2 v2.2.1
	github.com/dave/jennifer => github.com/mcrawfo2/jennifer v1.4.2
	github.com/rcrowley/go-metrics => github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
)
