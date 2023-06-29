package csv

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/turbot/go-kit/helpers"
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

	tableNames := []string{}
	re := regexp.MustCompile(`[-.]`)
	for _, i := range paths {
		tableCtx := context.WithValue(ctx, keyPath, i)
		base := strings.TrimSuffix(filepath.Base(i), gzipExtension)

		// check if the base name is already added as table name
		// then add the parent folder as prefix to the base for table name
		// if 2 same name file is present in 2 different path - /Users/dspl_countries.csv and /Users/test/dspl_countries.csv
		// there will be 2 tables - dspl_countries and test_dspl_countries
		if helpers.StringSliceContains(tableNames, base[0:len(base)-len(filepath.Ext(base))]) {
			folder_path := re.ReplaceAllString(strings.Split(i, "/")[len(strings.Split(i, "/"))-2], "_")
			tables[folder_path+"_"+base[0:len(base)-len(filepath.Ext(base))]], err = tableCSV(tableCtx, d.Connection)
		} else {
			tables[base[0:len(base)-len(filepath.Ext(base))]], err = tableCSV(tableCtx, d.Connection)
			tableNames = append(tableNames, base[0:len(base)-len(filepath.Ext(base))])
		}
		if err != nil {
			plugin.Logger(ctx).Error("csv.PluginTables", "create_table_error", err, "path", i)
			return nil, err
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
