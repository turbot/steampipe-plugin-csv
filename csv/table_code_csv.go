package csv

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dimchansky/utfbom"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func readCSV(ctx context.Context, connection *plugin.Connection, path string) (*csv.Reader, error) {

	file, err := os.Open(path)
	if err != nil {
		plugin.Logger(ctx).Error("csv.readCSV", "os_open_error", err, "path", path)
		return nil, err
	}

	var csvFile io.Reader
	if strings.HasSuffix(path, gzipExtension) {
		gzipFile, err := gzip.NewReader(file)
		if err != nil {
			plugin.Logger(ctx).Error("csv.readCSV", "gzip_open_error", err, "path", path)
			return nil, err
		}
		csvFile = gzipFile
	} else {
		csvFile = file
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

	return r, nil
}

func tableCSV(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {

	path := ctx.Value(keyPath).(string)
	r, err := readCSV(ctx, connection, path)
	if err != nil {
		plugin.Logger(ctx).Error("csv.tableCSV", "read_csv_error", err, "path", path)
		return nil, fmt.Errorf("failed to load csv file %s: %v", path, err)
	}

	// Read the header to peek at the column names
	header, err := r.Read()
	if err != nil {
		plugin.Logger(ctx).Error("csv.tableCSV", "header_parse_error", err, "path", path, "header", header)
		return nil, fmt.Errorf("failed to parse file header %s: %v", path, err)
	}

	csvConfig := GetConfig(connection)
	setDefaultHeader := checkCSVWithHeaderOption(ctx, *csvConfig.Header, header)

	cols := []*plugin.Column{}
	colNames := []string{}
	keys := make(map[string]bool)
	for idx, i := range header {
		// Set the default column name
		if setDefaultHeader {
			i = fmt.Sprintf("c%d", idx)
		// Table column names cannot be empty strings in default
		} else if len(i) == 0 {
			plugin.Logger(ctx).Error("csv.tableCSV", "empty_header_error", "header row has empty value", "path", path, "field", idx)
			return nil, fmt.Errorf("%s header row has empty value in field %d", path, idx)
		// Table column names cannot be duplicated in default
		} else {
			_, ok := keys[i]
			if !ok {
				keys[i] = true
			} else {
				plugin.Logger(ctx).Error("csv.tableCSV", "duplicated_header_error", "header row has duplicated value", "path", path, "field", idx)
				return nil, fmt.Errorf("%s header row has duplicated value in field %d", path, idx)
			}
		}
		colNames = append(colNames, i)
		cols = append(cols, &plugin.Column{Name: i, Type: proto.ColumnType_STRING, Transform: transform.FromField(helpers.EscapePropertyName(i)), Description: fmt.Sprintf("Field %d.", idx)})
	}

	return &plugin.Table{
		Name:        path,
		Description: fmt.Sprintf("CSV file at %s", path),
		List: &plugin.ListConfig{
			Hydrate: listCSVWithPath(path, setDefaultHeader, colNames),
		},
		Columns: cols,
	}, nil
}

func listCSVWithPath(path string, setDefaultHeader bool, colNames []string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

		r, err := readCSV(ctx, d.Connection, path)
		if err != nil {
			plugin.Logger(ctx).Error("csv.listCSVWithPath", "read_csv_error", err, "path", path)
			return nil, fmt.Errorf("failed to load csv file %s: %v", path, err)
		}

		// Header row consume or not
		if !setDefaultHeader {
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

func checkCSVWithHeaderOption(ctx context.Context, headerOption string, header []string) bool {

	// Conclude to use the default column names or not
	setDefaultHeader := false
	switch headerOption {
	case "auto":
		keys := make(map[string]bool)
		for _, i := range header {
			// Check the empty column name
			if len(i) == 0 {
				setDefaultHeader = true
				break
			}
			// Check the duplicated column name
			_, ok := keys[i]
			if ok {
				setDefaultHeader = true
				break
			} else {
				keys[i] = true
			}
		}
	case "off":
		setDefaultHeader = true
	case "on":
	default:
		plugin.Logger(ctx).Warn("csv.headerCSV", "unknown_header_option", "headerOption", headerOption)
	}

	return setDefaultHeader
}
