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
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func readCSV(ctx context.Context, connection *plugin.Connection, path string) (*csv.Reader, error) {

	// Only allow parsing of one file at a time to prevent concurrent map read
	// and write errors
	parseMutex.Lock()
	defer parseMutex.Unlock()

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
		return nil, fmt.Errorf("failed to load CSV file %s: %v", path, err)
	}

	// Read the header to peek at the column names
	header, err := r.Read()
	if err != nil {
		// Parse errors
		if len(header) > 0 {
			plugin.Logger(ctx).Error("csv.tableCSV", "header_parse_error", err, "path", path, "header", header)
			return nil, fmt.Errorf("failed to parse file header %s: %v", path, err)
		}
		
		// Return nil if the given file is empty, also add a log message to inform that the file is empty
		plugin.Logger(ctx).Error("csv.tableCSV", "skipping the file since empty", path)
		return nil, nil
	}

	// Determine whether to use the first row as the header row when creating column names:
	// - "auto": If there are no empty or duplicate values use the first row as the header. Else, use generic column names, e.g., "a", "b".
	// - "on": Use the first row as the header. If there are empty or duplicate values, the tables will fail to load.
	// - "off": Do not use the first row as the header. All column names will be generic.
	csvConfig := GetConfig(connection)
	headerMode := "auto"

	if csvConfig.Header != nil {
		headerMode = *csvConfig.Header
	}

	cols := []*plugin.Column{}
	colNames := []string{}
	var headerValue string

	// TODO: Can we read the header just once, collecting column names and rows
	// along the way?

	// If header mode is "off", no need to check if header is valid since it's
	// not used
	var isValidHeader bool
	var invalidReason string
	if headerMode == "auto" || headerMode == "on" {
		isValidHeader, invalidReason = validHeader(ctx, header)
	}

	useHeaderRow := true

	// Check if we should use header row
	switch headerMode {
	case "auto":
		if !isValidHeader {
			useHeaderRow = false
		}
	case "off":
		useHeaderRow = false
	case "on":
		if !isValidHeader {
			plugin.Logger(ctx).Error("csv.tableCSV", "invalid_header_error", invalidReason, "path", path)
			return nil, fmt.Errorf(invalidReason)
		}
	}

	for idx, i := range header {
		if useHeaderRow {
			headerValue = i
		} else {
			headerValue = intToLetters(idx + 1)
		}

		colNames = append(colNames, headerValue)
		cols = append(cols, &plugin.Column{Name: headerValue, Type: proto.ColumnType_STRING, Transform: transform.FromField(helpers.EscapePropertyName(headerValue)), Description: fmt.Sprintf("Field %d.", idx)})
	}

	return &plugin.Table{
		Name:        path,
		Description: fmt.Sprintf("CSV file at %s", path),
		List: &plugin.ListConfig{
			Hydrate: listCSVWithPath(path, useHeaderRow, colNames),
		},
		Columns: cols,
	}, nil
}

func listCSVWithPath(path string, useHeaderRow bool, colNames []string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

		r, err := readCSV(ctx, d.Connection, path)
		if err != nil {
			plugin.Logger(ctx).Error("csv.listCSVWithPath", "read_csv_error", err, "path", path)
			return nil, fmt.Errorf("failed to load CSV file %s: %v", path, err)
		}

		// Header rows should not be used as a data row
		if useHeaderRow {
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

// A valid header row has no empty values or duplicate values
func validHeader(ctx context.Context, header []string) (bool, string) {
	keys := make(map[string]bool)
	for idx, i := range header {
		// Check for empty column names
		if len(i) == 0 {
			return false, fmt.Sprintf("header row has empty value in field %d", idx)
		}
		// Check for duplicate column names
		_, ok := keys[i]
		if ok {
			return false, fmt.Sprintf("header row has duplicate value in field %d", idx)
		}
		keys[i] = true
	}

	// No empty or duplicate column names found
	return true, ""
}
