![image](https://hub.steampipe.io/images/plugins/turbot/csv-social-graphic.png)

# CSV Plugin for Steampipe

Use SQL to query data from CSV files.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/csv)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/csv/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-csv/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install csv
```

Configure your [config file](https://hub.steampipe.io/plugins/turbot/csv#configuration) to include directories with CSV files. If no directory is specified, the current working directory will be used.

Run steampipe:

```shell
steampipe query
```

Run a query for the `my_users.csv` file:

```sql
select
  first_name,
  last_name
from
  my_users;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-csv.git
cd steampipe-plugin-csv
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/csv.spc
```

Try it!

```
steampipe query
> .inspect csv
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-csv/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [CSV Plugin](https://github.com/turbot/steampipe-plugin-csv/labels/help%20wanted)
