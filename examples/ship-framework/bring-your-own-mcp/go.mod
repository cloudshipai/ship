module bring-your-own-mcp

go 1.23.2

require (
	github.com/cloudshipai/ship v0.6.7
	github.com/mark3labs/mcp-go v0.37.0
)

replace github.com/cloudshipai/ship => ../../../

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/invopop/jsonschema v0.13.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
