package csv

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
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
	paths, err := csvList(ctx, d.Connection)
	if err != nil {
		return nil, err
	}
	for _, i := range paths {
		tableCtx := context.WithValue(ctx, keyPath, i)
		base := strings.TrimSuffix(filepath.Base(i), gzipExtension)
		tables[base[0:len(base)-len(filepath.Ext(base))]], err = tableCSV(tableCtx, d.Connection)
		if err != nil {
			plugin.Logger(ctx).Error("csv.PluginTables", "create_table_error", err, "path", i)
			return nil, err
		}
	}

	return tables, nil
}

func csvList(ctx context.Context, connection *plugin.Connection) ([]string, error) {
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
		// Check to resolve ~ to home dir
		if strings.HasPrefix(i, "~") {
			// File system context
			home, err := os.UserHomeDir()
			if err != nil {
				plugin.Logger(ctx).Error("csv.csvList", "os.UserHomeDir error. ~ will not be expanded in paths.", err)
			}

			// Resolve ~ to home dir
			if home != "" {
				if i == "~" {
					i = home
				} else if strings.HasPrefix(i, "~/") {
					i = filepath.Join(home, i[2:])
				}
			}
		}

		// Get full path
		fullPath, err := filepath.Abs(i)
		if err != nil {
			plugin.Logger(ctx).Error("csv.csvList", "failed to fetch absolute path", err, "path", i)
			return nil, err
		}

		// Expand globs
		iMatches, err := doublestar.Glob(fullPath)
		if err != nil {
			// Fail if any path is an invalid glob
			return nil, fmt.Errorf("Path is not a valid glob: %s", i)
		}
		matches = append(matches, iMatches...)
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
