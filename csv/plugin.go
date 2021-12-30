package csv

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func PluginTables(ctx context.Context, p *plugin.Plugin) (map[string]*plugin.Table, error) {

	// Initialize tables
	tables := map[string]*plugin.Table{}

	// Search for CSV files to create as tables
	paths, err := csvList(ctx, p)
	if err != nil {
		return nil, err
	}
	for _, i := range paths {
		tableCtx := context.WithValue(ctx, "path", i)
		base := filepath.Base(i)
		tables[base[0:len(base)-len(filepath.Ext(base))]] = tableCSV(tableCtx, p)
	}

	return tables, nil
}

func csvList(ctx context.Context, p *plugin.Plugin) ([]string, error) {

	var csvFilePaths []string
	var paths []string

	// fetch current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Glob paths in config
	// set current working directory if no paths are specified
	csvConfig := GetConfig(p.Connection)
	if &csvConfig == nil || csvConfig.Paths == nil {
		paths = append(paths, cwd+"/*")
	} else {
		paths = csvConfig.Paths
	}

	// File system context
	home, err := os.UserHomeDir()
	if err != nil {
		plugin.Logger(ctx).Error("csvList", "os.UserHomeDir error. ~ will not be expanded in paths.", err)
	}

	// Gather file path matches for the glob
	var matches []string

	for _, i := range paths {

		// Resolve ~ to home dir
		if home != "" {
			if i == "~" {
				i = home
			} else if strings.HasPrefix(i, "~/") {
				i = filepath.Join(home, i[2:])
			}
		}

		// Expand globs
		iMatches, err := filepath.Glob(i)
		if err != nil {
			// Fail if any path is an invalid glob
			return matches, fmt.Errorf("path is not a valid glob: %s", i)
		}

		matches = append(matches, iMatches...)
	}

	// Sanitize the matches to likely csvfiles
	for _, i := range matches {

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

		// This file was expanded from the glob, so check it's likely to be
		// of the right type based on the name / extension.
		ext := strings.ToLower(filepath.Ext(i))
		if ext == ".csv" {
			csvFilePaths = append(csvFilePaths, i)
		}
	}

	return csvFilePaths, nil
}
