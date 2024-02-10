//go:generate oapi-codegen -generate gin -o openapi/openapi_gin.gen.go -p openapi ../../spec/openapi.yaml
//go:generate oapi-codegen -generate types -o openapi/openapi_types.gen.go -p openapi ../../spec/openapi.yaml
//go:generate oapi-codegen -generate client -o openapi/openapi_client.gen.go -p openapi ../../spec/openapi.yaml
//go:generate oapi-codegen -generate spec -o openapi/openapi_spec.gen.go -p openapi ../../spec/openapi.yaml

package api
