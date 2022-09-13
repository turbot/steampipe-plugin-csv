package csv

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type csvConfig struct {
	Paths     []string `cty:"paths"`
	Separator *string  `cty:"separator"`
	Comment   *string  `cty:"comment"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"paths": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
	"separator": {
		Type: schema.TypeString,
	},
	"comment": {
		Type: schema.TypeString,
	},
}

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
