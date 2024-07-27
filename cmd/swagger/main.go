package main

//go:generate swag init -g main.go -o ../../docs --ot yaml

// @title    ${service_title}
// @version  1.0
// @host     ${service_name}
// @Accept   json
// @Produce  json

// @x-google-endpoints  [{"name": "${service_name}", "target": "${ip_address}", "allowCors": true}]
// @schemes             http https

// @securityDefinitions.oauth2.implicit  firebase
// @authorizationUrl
// @x-google-issuer     ${auth.firebase_issuer}
// @x-google-jwks_uri   ${auth.firebase_jwks_uri}
// @x-google-audiences  ${auth.firebase_audiences}

// @securityDefinitions.oauth2.implicit  sandstorm-backend
// @authorizationUrl
// @x-google-issuer     ${auth.sandstorm_backend_x_google_issuer}
// @x-google-jwks_uri   ${auth.sandstorm_backend_x_google_jwks_uri}
// @x-google-audiences  ${auth.sandstorm_backend_x_google_audiences}

func main() {
	// this module is needed only for building swagger documentation
	// and main function is nessesary for prebuild application in pipelines
}
