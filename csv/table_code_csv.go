package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableCSV(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {

	path := ctx.Value(keyPath).(string)
	csvFile, err := os.Open(path)
	if err != nil {
		plugin.Logger(ctx).Error("csv.tableCSV", "os_open_error", err, "path", path)
		return nil, err
	}

	r := csv.NewReader(csvFile)

	csvConfig := GetConfig(connection)
	if csvConfig.Separator != nil && *csvConfig.Separator != "" {
		r.Comma = rune((*csvConfig.Separator)[0])
	}
	if csvConfig.Comment != nil {
		if *csvConfig.Comment == "" {
			// Disable comments
			r.Comment = 0
		} else {
			// Set the comment character
			r.Comment = rune((*csvConfig.Comment)[0])
		}
	}

	// Read the header to peek at the column names
	header, err := r.Read()
	if err != nil {
		plugin.Logger(ctx).Error("csv.tableCSV", "header_parse_error", err, "path", path, "header", header)
		return nil, fmt.Errorf("failed to parse file header %s: %v", path, err)
	}

	cols := []*plugin.Column{}
	for idx, i := range header {
		// Table column names cannot be empty strings
		if len(i) == 0 {
			plugin.Logger(ctx).Error("csv.tableCSV", "empty_header_error", "header row has empty value", "path", path, "field", idx)
			return nil, fmt.Errorf("%s header row has empty value in field %d", path, idx)
		}
		cols = append(cols, &plugin.Column{Name: i, Type: proto.ColumnType_STRING, Transform: transform.FromField(helpers.EscapePropertyName(i)), Description: fmt.Sprintf("Field %d.", idx)})
	}

	return &plugin.Table{
		Name:        path,
		Description: fmt.Sprintf("CSV file at %s", path),
		List: &plugin.ListConfig{
			Hydrate: listCSVWithPath(path),
		},
		Columns: cols,
	}, nil
}

func listCSVWithPath(path string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

		csvFile, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		r := csv.NewReader(csvFile)

		csvConfig := GetConfig(d.Connection)
		if csvConfig.Separator != nil && *csvConfig.Separator != "" {
			r.Comma = rune((*csvConfig.Separator)[0])
		}
		if csvConfig.Comment != nil {
			if *csvConfig.Comment == "" {
				// Disable comments
				r.Comment = 0
			} else {
				// Set the comment character
				r.Comment = rune((*csvConfig.Comment)[0])
			}
		}

		header, err := r.Read()
		if err != nil {
			plugin.Logger(ctx).Error("csv.listCSVWithPath", "header_parse_error", err, "path", path, "header", header)
			return nil, err
		}

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				plugin.Logger(ctx).Error("csv.listCSVWithPath", "record_parse_error", err, "path", path, "record", record)
				continue
			}
			row := map[string]string{}
			for idx, j := range record {
				row[header[idx]] = j
			}
			d.StreamListItem(ctx, row)
		}

		return nil, nil
	}
}
