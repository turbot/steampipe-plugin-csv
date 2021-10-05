connection "csv" {
  plugin = "csv"

  # Paths is a list of locations to search for CSV files. Each file will be
  # converted to a table. Wildcards are supported per
  # https://golang.org/pkg/path/filepath/#Match
  # Exact file paths can have any name. Wildcard based matches must have an
  # extension of .csv (case insensitive).
  # paths = [ "/path/to/dir/*", "/path/to/exact/custom.csv" ]

  # The field delimiter character when parsing CSV files. Must be a single
  # character. Defaults to comma.
  # separator = ","

  # If set, then lines beginning with the comment character without preceding
  # whitespace are ignored. Disabled by default.
  # comment = "#"
}
