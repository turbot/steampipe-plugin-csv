package csv

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
	// KeyValue has been added to avoid key collisions
    KeyValue key = "path"
)


func PluginTables(ctx context.Context, p *plugin.Plugin) (map[string]*plugin.Table, error) {
	// Initialize tables
	tables := map[string]*plugin.Table{}

	// Search for CSV files to create as tables
	paths, err := csvList(ctx, p)
	if err != nil {
		return nil, err
	}
	for _, i := range paths {
		tableCtx := context.WithValue(ctx, KeyValue, i)
		base := filepath.Base(i)
		tables[base[0:len(base)-len(filepath.Ext(base))]], err = tableCSV(tableCtx, p)
		if err != nil {
			plugin.Logger(ctx).Error("csv.PluginTables", "create_table_error", err, "path", i)
			return nil, err
		}
	}

	return tables, nil
}

func csvList(ctx context.Context, p *plugin.Plugin) ([]string, error) {
	// Glob paths in config
	// Fail if no paths are specified
	csvConfig := GetConfig(p.Connection)
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

	// Sanitize the matches to likely csvfiles
	var csvFilePaths []string
	for _, i := range matches {
		// Check if file or directory
		fileInfo, err := os.Stat(i)
		if err != nil {
			plugin.Logger(ctx).Error("utils.tfConfigList", "error getting file info", err, "path", i)
			return nil, err
		}

		// Ignore directories
		if fileInfo.IsDir() {
			continue
		}

		// If the file path is an exact match to a matrix path then it's always
		// treated as a match - it was requested exactly
		hit := false
		for _, j := range paths {
			if i == j {
				hit = true
				break
			}
		}
		if hit {
			csvFilePaths = append(csvFilePaths, i)
			continue
		}
		csvFilePaths = append(csvFilePaths, i)
	}

	return csvFilePaths, nil
}
