---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/csv.svg"
brand_color: "#000000"
display_name: "CSV"
short_name: "csv"
description: "Steampipe plugin to query data from CSV files."
og_description: "Query CSV files with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/csv-social-graphic.png"
---

# CSV + Steampipe

A comma-separated values (CSV) file is a delimited text file that contains records of data.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query data using SQL.

Query data from the `my_users.csv` file:

```sql
select
  first_name,
  last_name
from
  my_users;
```

```sh
+------------+-----------+
| first_name | last_name |
+------------+-----------+
| Michael    | Scott     |
| Dwight     | Schrute   |
+------------+-----------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/csv/tables)**

## Get started

### Install

Download and install the latest CSV plugin:

```bash
steampipe plugin install csv
```

### Credentials

No credentials are required.

### Configuration

Installing the latest csv plugin will create a config file (`~/.steampipe/config/csv.spc`) with a single connection named `csv`:

```hcl
connection "csv" {
  plugin = "csv"

  # Paths is a list of locations to search for CSV files
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

  # Whether to use the first row as the header row when creating column names.
  # Valid values are "auto", "on", "off":
  #   - "auto": If there are no empty or duplicate values use the first row as the header; else, use generic column names, e.g., "c1", "c2".
  #   - "on": Use the first row as the header. If there are empty or duplicate values, the tables will fail to load.
  #   - "off": Do not use the first row as the header. All column names will be generic.
  # Defaults to "auto".
  # header = "auto"
}
```

- `paths` - A list of directory paths to search for CSV files. Paths are resolved relative to the current working directory. Paths may [include wildcards](https://pkg.go.dev/path/filepath#Match) and also supports `**` for recursive matching. Defaults to the current working directory.
- `separator` - Field delimiter when parsing files. Defaults to `,`.
- `comment` - Lines starting with this comment character are ignored. Disabled by default.
- `header` - Whether to use the first row as the header row when creating column names. Valid values are "auto", "on", "off". Defaults to "auto".

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-csv
- Community: [Slack Channel](https://steampipe.io/community/join)
