module github.com/communitybridge/easycla/cla-backend-go

go 1.14

replace github.com/awslabs/aws-lambda-go-api-proxy => github.com/LF-Engineering/aws-lambda-go-api-proxy v0.3.2

require (
	github.com/LF-Engineering/aws-lambda-go-api-proxy v0.3.2
	github.com/LF-Engineering/lfx-kit v0.1.10
	github.com/LF-Engineering/lfx-models v0.3.9
	github.com/aws/aws-lambda-go v1.16.0
	github.com/aws/aws-sdk-go v1.30.14
	github.com/aymerick/raymond v2.0.2+incompatible
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fnproject/fdk-go v0.0.2
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime v0.19.15
	github.com/go-openapi/spec v0.19.7
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/go-openapi/validate v0.19.8
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/gorilla/sessions v1.2.0 // indirect
	github.com/imroc/req v0.3.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lytics/logrus v0.0.0-20170528191427-4389a17ed024
	github.com/rs/cors v1.7.0
	github.com/savaki/dynastore v0.0.0-20171109173440-28d8558bb429
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tencentyun/scf-go-lib v0.0.0-20200116145541-9a6ea1bf75b8
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
