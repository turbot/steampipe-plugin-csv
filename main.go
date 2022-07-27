package main

import (
	"github.com/turbot/steampipe-plugin-csv/csv"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: csv.Plugin})
}
