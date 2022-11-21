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

Query data from the `users.csv` file:

```sql
select
  first_name,
  last_name
from
  users;
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

  # Determine whether to use the first row as the header row when creating column names.
  # Valid values are "auto", "on", "off":
  #   - "auto": If there are no empty or duplicate values use the first row as the header; else, use the first row as a data row and use generic column names, e.g., "a", "b".
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

### Header row

By default, the `header` configuration argument is set to `auto`, so when CSV files are loaded, the first row will be checked if it's a valid header row, i.e., no value or duplicate values.

For instance, for the following CSV file `users.csv`:

```csv
first_name,last_name,email
Michael,Scott,mscott@dmi.com
Dwight,Schrute,dschrute@dmi.com
Pamela,Beesly,pbeesly@dmi.com
```

The CSV plugin will create a table called `users` with the header values as column names:

```bash
.inspect csv.users
+------------+------+-------------+
| column     | type | description |
+------------+------+-------------+
| first_name | text | Field 0.    |
| last_name  | text | Field 1.    |
| email      | text | Field 2.    |
+------------+------+-------------+
```

Which produces the following query results:

```bash
> select * from dmi
+------------+-----------+------------------+
| first_name | last_name | email            |
+------------+-----------+------------------+
| Dwight     | Schrute   | dschrute@dmi.com |
| Michael    | Scott     | mscott@dmi.com   |
| Pamela     | Beesly    | pbeesly@dmi.com  |
+------------+-----------+------------------+
```

However, if the `users.csv` contained the following data:

```csv
first_name,,email
Michael,Scott,mscott@dmi.com
Dwight,Schrute,dschrute@dmi.com
Pamela,Beesly,pbeesly@dmi.com
```

The CSV plugin will assume the first row is not the header row and will create a table called `users` with positional column names:

```bash
.inspect csv.users
+--------+------+-------------+
| column | type | description |
+--------+------+-------------+
| a      | text | Field 0.    |
| b      | text | Field 1.    |
| c      | text | Field 2.    |
+--------+------+-------------+
```

Which produces the following query results:

```bash
> select * from dmi
+------------+---------+------------------+
| a          | b       | c                |
+------------+---------+------------------+
| first_name |         | email            |
| Pamela     | Beesly  | pbeesly@dmi.com  |
| Dwight     | Schrute | dschrute@dmi.com |
| Michael    | Scott   | mscott@dmi.com   |
+------------+---------+------------------+
```

The `header` configuration argument can also be set to:
- `on`: This setting requires the first row to be a valid header row, else the plugin will fail to create the tables.
- `off`: This setting always assumes the first row isn't the header row and uses positional column names for all tables.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-csv
- Community: [Slack Channel](https://steampipe.io/community/join)
