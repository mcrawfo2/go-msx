module cto-github.cisco.com/NFV-BU/go-msx

go 1.14

require (
	cto-github.cisco.com/NFV-BU/msx-platform-go-client v1.0.4
	github.com/AlecAivazis/survey/v2 v2.0.5
	github.com/Shopify/sarama v1.26.1
	github.com/ThreeDotsLabs/watermill v1.0.2
	github.com/ThreeDotsLabs/watermill-kafka/v2 v2.2.0
	github.com/ThreeDotsLabs/watermill-sql v1.3.4
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535 // indirect
	github.com/benbjohnson/clock v1.0.3
	github.com/bmatcuk/doublestar v1.1.5
	github.com/dave/jennifer v1.4.1
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/doug-martin/goqu/v9 v9.9.0
	github.com/emicklei/go-restful v2.14.3+incompatible
	github.com/emicklei/go-restful-openapi v1.2.0
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813
	github.com/getkin/kin-openapi v0.20.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-ini/ini v1.48.0
	github.com/go-openapi/spec v0.20.3
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/gocql/gocql v0.0.0-20191013011951-93ce931da9e1
	github.com/google/uuid v1.1.1
	github.com/hashicorp/consul v1.4.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/vault/api v1.0.4
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0
	github.com/magiconair/properties v1.8.1
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/opentracing/opentracing-go v1.1.0
	github.com/otiai10/copy v1.0.2
	github.com/pavel-v-chernykh/keystore-go v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang v1.2.1
	github.com/radovskyb/watcher v1.0.7
	github.com/robfig/cron/v3 v3.0.1
	github.com/scylladb/go-reflectx v1.0.1
	github.com/scylladb/gocqlx v1.3.1
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/thejerf/abtime v1.0.3
	github.com/tidwall/gjson v1.6.4
	github.com/uber/jaeger-client-go v2.19.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible
	go.uber.org/atomic v1.6.0
	gopkg.in/pipe.v2 v2.0.0-20140414041502-3c2ca4d52544
	gopkg.in/yaml.v2 v2.4.0
	vitess.io/vitess v0.0.0-20191026003914-d26b6c7975b1
)

replace (
	github.com/ThreeDotsLabs/watermill-kafka/v2 v2.2.0 => cto-github.cisco.com/NFV-BU/watermill-kafka/v2 v2.2.1
	github.com/dave/jennifer => github.com/mcrawfo2/jennifer v1.4.2
)
