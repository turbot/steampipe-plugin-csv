package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/dimchansky/utfbom"
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

	// Some CSV files have a non-standard Byte Order Mark (BOM) at the start
	// of the file - for example, UTF-8 encoded CSV files from Excel. This
	// messes up the first column name, so skip the BOM if found.
	csvFileWithoutBom, enc := utfbom.Skip(csvFile)
	plugin.Logger(ctx).Debug("csv.tableCSV", "path", path, "detected_encoding", enc)

	r := csv.NewReader(csvFileWithoutBom)

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

	// Conclude to use the default column names or not
	isHeader := true
	keys := make(map[string]bool)
	for _, i := range header {
		// Check the empty column name
		if len(i) == 0 {
			isHeader = false
			break
		}
		// Check the duplicated column name
		_, ok := keys[i]
		if ok {
			isHeader = false
			break
		} else {
			keys[i] = true
		}
	}

	cols := []*plugin.Column{}
	colNames := []string{}
	for idx, i := range header {
		// Set the default column name
		if !isHeader {
			i = fmt.Sprintf("c%d", idx)
		}
		colNames = append(colNames, i)
		cols = append(cols, &plugin.Column{Name: i, Type: proto.ColumnType_STRING, Transform: transform.FromField(helpers.EscapePropertyName(i)), Description: fmt.Sprintf("Field %d.", idx)})
	}

	return &plugin.Table{
		Name:        path,
		Description: fmt.Sprintf("CSV file at %s", path),
		List: &plugin.ListConfig{
			Hydrate: listCSVWithPath(path, isHeader, colNames),
		},
		Columns: cols,
	}, nil
}

func listCSVWithPath(path string, isHeader bool, colNames []string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

		csvFile, err := os.Open(path)
		if err != nil {
			plugin.Logger(ctx).Error("csv.listCSVWithPath", "os_open_error", err, "path", path)
			return nil, err
		}

		// Some CSV files have a non-standard Byte Order Mark (BOM) at the start
		// of the file - for example, UTF-8 encoded CSV files from Excel. This
		// messes up the first column name, so skip the BOM if found.
		csvFileWithoutBom, enc := utfbom.Skip(csvFile)
		plugin.Logger(ctx).Debug("csv.listCSVWithPath", "path", path, "detected_encoding", enc)

		r := csv.NewReader(csvFileWithoutBom)

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

		// Header row consume or not
		if isHeader {
			header, err := r.Read()
			if err != nil {
				plugin.Logger(ctx).Error("csv.listCSVWithPath", "header_parse_error", err, "path", path, "header", header)
				return nil, err
			}
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
				row[colNames[idx]] = j
			}
			d.StreamListItem(ctx, row)
		}

		return nil, nil
	}
}
