connection "csv" {
  plugin = "csv"

  # Paths is a list of locations to search for CSV files
  # Paths can be configured with a local directory, a remote Git repository URL, or an S3 bucket URL
  # Refer https://hub.steampipe.io/plugins/turbot/csv#supported-path-formats for more information
  # All paths are resolved relative to the current working directory (CWD)
  # Wildcard based searches are supported, including recursive searches

  # For example:
  #  - "*.csv" matches all CSV files in the CWD
  #  - "*.csv.gz" matches all gzipped CSV files in the CWD
  #  - "**/*.csv" matches all CSV files in the CWD and all sub-directories
  #  - "../*.csv" matches all CSV files in the CWD's parent directory
  #  - "steampipe*.csv" matches all CSV files starting with "steampipe" in the CWD
  #  - "/path/to/dir/*.csv" matches all CSV files in a specific directory
  #  - "/path/to/dir/custom.csv" matches a specific file

  # If paths includes "*", all files (including non-CSV files) in
  # the CWD will be matched, which may cause errors if incompatible file types exist

  # Defaults to CWD
  paths = [ "*.csv", "*.csv.gz" ]

  # The field delimiter character when parsing CSV files. Must be a single
  # character. Defaults to comma.
  # separator = ","

  # If set, then lines beginning with the comment character without preceding
  # whitespace are ignored. Disabled by default.
  # comment = "#"

  # Determine whether to use the first row as the header row when creating column names.
  # Valid values are "auto", "on", "off":
  #   - "auto": If there are no empty or duplicate values use the first row as the header; else, use the first row as a data row and use generic column names, e.g., "a", "b".
  #   - "on": Use the first row as the header. If there are empty or duplicate values, the tables will fail to load.
  #   - "off": Do not use the first row as the header. All column names will be generic.
  # Defaults to "auto".
  # header = "auto"
}
