package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableCSV(ctx context.Context, p *plugin.Plugin) *plugin.Table {

	path := ctx.Value("path").(string)
	csvFile, err := os.Open(path)
	if err != nil {
		plugin.Logger(ctx).Error("Could not open CSV file", "path", path)
		panic(err)
	}

	r := csv.NewReader(csvFile)

	csvConfig := GetConfig(p.Connection)
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

	// Read the header to peak at the column names
	header, err := r.Read()
	if err != nil {
		plugin.Logger(ctx).Error("Error parsing CSV header:", "path", path, "header", header, "err", err)
		panic(err)
	}
	cols := []*plugin.Column{}
	for idx, i := range header {
		cols = append(cols, &plugin.Column{Name: i, Type: proto.ColumnType_STRING, Transform: transform.FromField(i), Description: fmt.Sprintf("Field %d.", idx)})
	}

	return &plugin.Table{
		Name:        path,
		Description: fmt.Sprintf("CSV file at %s", path),
		List: &plugin.ListConfig{
			Hydrate: listCSVWithPath(path),
		},
		Columns: cols,
	}
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
			plugin.Logger(ctx).Error("Error parsing CSV header:", "path", path, "header", header, "err", err)
			return nil, err
		}

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				plugin.Logger(ctx).Error("Error parsing CSV record:", "path", path, "record", record, "err", err)
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
