module github.com/devingen/sepet-cdn

go 1.12

//replace github.com/devingen/api-core => ../api-core

require (
	github.com/aws/aws-sdk-go v1.0.0
	github.com/devingen/api-core v0.0.21
	github.com/go-ini/ini v1.62.0 // indirect
	github.com/go-resty/resty/v2 v2.4.0
	github.com/gorilla/mux v1.7.4
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	go.elastic.co/apm/module/apmhttp v1.9.0
	go.mongodb.org/mongo-driver v1.3.2
)
