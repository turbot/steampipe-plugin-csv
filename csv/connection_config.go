package csv

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type csvConfig struct {
	Paths     []string `hcl:"paths,optional" steampipe:"watch"`
	Separator *string  `hcl:"separator"`
	Comment   *string  `hcl:"comment"`
	Header    *string  `hcl:"header"`
}

// var ConfigSchema = map[string]*schema.Attribute{
// 	"paths": {
// 		Type: schema.TypeList,
// 		Elem: &schema.Attribute{Type: schema.TypeString},
// 	},
// 	"separator": {
// 		Type: schema.TypeString,
// 	},
// 	"comment": {
// 		Type: schema.TypeString,
// 	},
// 	"header": {
// 		Type: schema.TypeString,
// 	},
// }

func ConfigInstance() interface{} {
	return &csvConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) csvConfig {
	if connection == nil || connection.Config == nil {
		return csvConfig{}
	}
	config, _ := connection.Config.(csvConfig)
	return config
}
