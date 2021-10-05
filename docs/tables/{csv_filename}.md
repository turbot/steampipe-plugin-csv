# Table: {csv_filename}

Query data from CSV files. A table is automatically created to represent each CSV file
found in the configured `paths`.

## Examples

### Inspect the table structure

Assuming your connection is called `csv` (the default), list all tables with:
```sql
.inspect csv
```

To get defails for a specific table, inspect it by name:
```sql
.inspect csv.users
```

### Query a simple file

Given the file `users.csv`, the query is:

```sql
select
  *
from
  users
```

### Query a complex file name

Given the file `My complex file-name.csv`, the query uses identifier quotes:

```sql
select
  *
from
  "My complex file-name"
```
