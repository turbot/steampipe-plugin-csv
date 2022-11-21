# Table: {csv_filename}

Query data from CSV files. A table is automatically created to represent each
CSV file found in the configured `paths`.

**Note**: All examples in this document assume the `header` configuration argument is set to `auto` (default value). For more information on how column names are created, please see [Header row](https://hub.steampipe.io/plugins/turbot/csv#header-row).

For instance, if `paths` is set to `/Users/myuser/csv/*.csv`, and that directory contains:

- products.csv:
  ```csv
  product_name,sku
  Paper,P001
  Better Paper,P002
  ```
- users.csv:
  ```csv
  first_name,last_name,email
  Michael,Scott,mscott@dmi.com
  Dwight,Schrute,dschrute@dmi.com
  Pamela,Beesly,pbeesly@dmi.com
  ```

This plugin will create 2 tables:

- products
- users

Which you can then query directly:

```sql
select
  *
from
  users;
```

All column values are returned as text data type.

## Examples

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

Given the file `users.csv`, the query is:

```sql
select
  *
from
  users;
```

### Query a complex file name

Given the file `My complex file-name.csv`, the query uses identifier quotes:

```sql
select
  *
from
  "My complex file-name"
```

### Query specific columns

Columns are always in text form when read from the CSV file. The column names come from the first row of the file.

```sql
select
  first_name,
  last_name
from
  users
```

If your column names are complex, use identifier quotes:

```sql
select
  "First Name",
  "Last Name"
from
  users
```

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
select
  first_name,
  age::int as iage
from
  users
where
  iage > 25
```

### Query multiple CSV files

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

```sql
create view all_ips as select * from ips1 union select * from ips2;
select * from all_ips
```
