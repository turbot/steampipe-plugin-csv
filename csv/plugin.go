package csv

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-csv",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		SchemaMode:       plugin.SchemaModeDynamic,
		TableMapFunc:     PluginTables,
	}
	return p
}

type key string

const (
	// keyPath has been added to avoid key collisions
	keyPath       key    = "path"
	gzipExtension string = ".gz"
)

func PluginTables(ctx context.Context, d *plugin.TableMapData) (map[string]*plugin.Table, error) {
	// Initialize tables
	tables := map[string]*plugin.Table{}

	// Search for CSV files to create as tables
	paths, err := csvList(ctx, d.Connection, d)
	if err != nil {
		return nil, err
	}
	for _, i := range paths {
		tableCtx := context.WithValue(ctx, keyPath, i)
		base := strings.TrimSuffix(filepath.Base(i), gzipExtension)

		tableData, err := tableCSV(tableCtx, d.Connection)
		if err != nil {
			plugin.Logger(ctx).Error("csv.PluginTables", "create_table_error", err, "path", i)
			return nil, err
		}

		// Skip the table if the file is empty
		if tableData != nil {
			tables[base[0:len(base)-len(filepath.Ext(base))]] = tableData
		}
	}

	return tables, nil
}

func csvList(ctx context.Context, connection *plugin.Connection, d *plugin.TableMapData) ([]string, error) {
	// Glob paths in config
	// Fail if no paths are specified
	csvConfig := GetConfig(connection)
	if csvConfig.Paths == nil {
		return nil, errors.New("paths must be configured")
	}

	// Gather file path matches for the glob
	var matches []string
	paths := csvConfig.Paths
	for _, i := range paths {
		files, err := d.GetSourceFiles(i)
		if err != nil {
			plugin.Logger(ctx).Error("csv.csvList", "failed to fetch absolute path", err, "path", i)
			return nil, err
		}
		matches = append(matches, files...)
	}

	// Sanitize the matches to ignore the directories
	var csvFilePaths []string
	for _, i := range matches {
		// Check if file or directory
		fileInfo, err := os.Stat(i)
		if err != nil {
			plugin.Logger(ctx).Error("csv.csvList", "error getting file info", err, "path", i)
			return nil, err
		}

		// Ignore directories
		if fileInfo.IsDir() {
			continue
		}

		csvFilePaths = append(csvFilePaths, i)
	}

	return csvFilePaths, nil
}
