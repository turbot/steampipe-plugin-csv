---
title: "Steampipe Table: csv_filename - Query OCI CSV Files using SQL"
description: "Allows users to query CSV Files in OCI, specifically to extract, transform, and load data from CSV files for analysis and reporting."
---

# Table: csv_filename - Query OCI CSV Files using SQL

Oracle Cloud Infrastructure (OCI) CSV Files are a type of structured data file format that stores tabular data, such as a spreadsheet or database in plain text. These files can be easily imported and exported from programs that store data in tables, such as Microsoft Excel or Open Office Calc. CSV files are primarily used to transport data between applications that handle lots of data, and they support a wide array of data types and are flexible in terms of what type of data they can hold.

## Table Usage Guide

The `csv_filename` table provides insights into the data stored within CSV files in OCI. As a data analyst, you can leverage this table to extract, transform, and load data from CSV files for in-depth analysis and reporting. Use it to uncover valuable insights from your data, such as identifying trends, patterns, and correlations. 

Schema link: [https://hub.steampipe.io/plugins/turbot/csv/tables/csv_filename](https://hub.steampipe.io/plugins/turbot/csv/tables/csv_filename)

## Examples

**Note**: All examples in this section assume the `header` configuration argument is set to `auto` (default value). For more information on how column names are created, please see [Column Names](https://hub.steampipe.io/plugins/turbot/csv/tables/{csv_filename}#column-names).

### Inspect the table structure

Assuming your connection is called `csv` (the default), list all tables with:

```bash
.inspect csv
+----------+--------------------------------------------+
| table    | description                                |
+----------+--------------------------------------------+
| products | CSV file at /Users/myuser/csv/products.csv |
| users    | CSV file at /Users/myuser/csv/users.csv    |
+----------+--------------------------------------------+
```

To get defails for a specific table, inspect it by name:

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

### Query a simple file
Explore all user data to gain a comprehensive understanding of your user base. This could be beneficial in identifying trends, understanding user behavior, and informing strategic decisions.
Given the file `users.csv`, the query is:


```sql+postgres
select
  *
from
  users;
```

```sql+sqlite
select
  *
from
  users;
```

### Query a complex file name
Explore the intricacies of a complex file by analyzing its various attributes. This is useful for understanding the file's structure and content in a comprehensive manner.
Given the file `My complex file-name.csv`, the query uses identifier quotes:


```sql+postgres
select
  *
from
  "My complex file-name"
```

```sql+sqlite
select
  *
from
  `My complex file-name`
```

### Query specific columns
Determine the areas in which you want to focus by selecting specific user details. This is useful when you want to narrow down your data analysis to specific user attributes.
Columns are always in text form when read from the CSV file. The column names come from the first row of the file.


```sql+postgres
select
  first_name,
  last_name
from
  users
```

```sql+sqlite
select
  first_name,
  last_name
from
  users
```

If your column names are complex, use identifier quotes:

```sql+postgres
select
  "First Name",
  "Last Name"
from
  users
```

```sql+sqlite
select
  "First Name",
  "Last Name"
from
  users
```

### Casting column data for analysis
Analyze the settings to understand which users are older than 25. This can be useful in tailoring age-specific content or offers.
Text columns can be easily cast to other types:


```sql+postgres
select
  first_name,
  age::int as iage
from
  users
where
  iage > 25
```

```sql+sqlite
select
  first_name,
  CAST(age as INTEGER) as iage
from
  users
where
  CAST(age as INTEGER) > 25
```

### Query multiple CSV files
Determine the contents of multiple CSV files in a unified view. This is beneficial when needing to analyze or compare data from multiple sources simultaneously.
Given this data:

ips1.csv:

```csv
service,ip_addr
service1,54.176.63.153
service2,222.236.38.99
```

ips2.csv:

```csv
service,ip_addr
service3,41.65.221.12
service4,83.151.87.112
service5,85.188.10.179
```

You can query both files like so:


```sql+postgres
create view all_ips as select * from ips1 union select * from ips2;
select * from all_ips
```

```sql+sqlite
create view all_ips as select * from ips1 union select * from ips2;
select * from all_ips
```

## Column Names

By default, the `header` configuration argument is set to `auto`, so when CSV files are loaded, the first row will be checked if it's a valid header row, i.e., no missing or duplicate values.

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

However, if the first row in `users.csv` was missing a value:

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