module github.com/coreos/etcd-operator

go 1.13

require (
	cloud.google.com/go/bigquery v1.4.0 // indirect
	cloud.google.com/go/storage v1.5.0
	github.com/Azure/azure-sdk-for-go v43.3.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.16
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190125095113-2b29687e15f2
	github.com/aws/aws-sdk-go v1.36.25
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/coreos/bbolt v1.3.4 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dnaeon/go-vcr v1.1.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/go-cmp v0.4.1 // indirect
	github.com/googleapis/gnostic v0.4.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.5 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	go.opencensus.io v0.22.3 // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20200601175630-2caf76543d99 // indirect
	google.golang.org/api v0.15.1
	google.golang.org/appengine v1.6.6 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
	k8s.io/api v0.16.10
	k8s.io/apiextensions-apiserver v0.16.10
	k8s.io/apimachinery v0.16.10
	k8s.io/client-go v0.20.2
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1

replace k8s.io/client-go => k8s.io/client-go v0.16.10

replace github.com/coreos/bbolt => github.com/coreos/bbolt v1.3.1-coreos.6
