## v0.11.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters.
- Recompiled plugin with Go version `1.21`.

## v0.10.0 [2023-09-14]

_Bug fixes_

- Added the missing [S3 go-getter](https://hub.steampipe.io/plugins/turbot/csv#configuring-s3-urls) examples in the `docs/index.md` file.

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v551-2023-07-26). ([#70](https://github.com/turbot/steampipe-plugin-csv/pull/70))
- Recompiled plugin with `github.com/turbot/go-kit v0.8.0`. ([#71](https://github.com/turbot/steampipe-plugin-csv/pull/71))

## v0.9.1 [2023-07-05]

_Bug fixes_

- Fixed the plugin to return a warning message if the configured `paths` argument has empty CSV file(s), instead of returning a plugin initialization error. ([#64](https://github.com/turbot/steampipe-plugin-csv/pull/64))

## v0.9.0 [2023-06-20]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. This update significantly lowers the plugin initialization time of dynamic plugins by avoiding recursing into child folders when not necessary. ([#61](https://github.com/turbot/steampipe-plugin-csv/pull/61))

## v0.8.0 [2023-05-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05) which fixes increased plugin initialization time due to multiple connections causing the schema to be loaded repeatedly. ([#58](https://github.com/turbot/steampipe-plugin-csv/pull/58))

## v0.7.0 [2023-03-20]

_What's new?_

- Added support for CSV files from remote Git repositories and S3 buckets. For more information, please see [Supported Path Formats](https://hub.steampipe.io/plugins/turbot/csv#supported-path-formats). ([#55](https://github.com/turbot/steampipe-plugin-csv/pull/55))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes support for fetching remote files with go-getter for dynamic plugins. ([#55](https://github.com/turbot/steampipe-plugin-csv/pull/55))

## v0.6.0 [2023-03-09]

_What's new?_

- Added file watching support for files included in the `paths` config argument. ([#53](https://github.com/turbot/steampipe-plugin-csv/pull/53))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.2.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v520-2023-03-02) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables.([#53](https://github.com/turbot/steampipe-plugin-csv/pull/53))

## v0.5.0 [2022-11-21]

_What's new?_

- Added support to read gzipped CSV files. ([#41](https://github.com/turbot/steampipe-plugin-csv/pull/41)) (Thanks to [@daeho-ro](https://github.com/daeho-ro) for the new feature!)
- Added `header` config argument, which allows you to set if the first row should be used a header row when determining each table's column names. By default, `header` is set to `auto`, so if the first row is a valid header row, i.e., no missing or duplicate values, its values will be used for column names; else, generic column names, e.g., "a", "b", will be used. For more information, please see [Column Names](https://hub.steampipe.io/plugins/turbot/csv/tables/{csv_filename}#column-names). ([#42](https://github.com/turbot/steampipe-plugin-csv/pull/42)) ([#47](https://github.com/turbot/steampipe-plugin-csv/pull/47)) ([#48](https://github.com/turbot/steampipe-plugin-csv/pull/48)) (Thanks to [@daeho-ro](https://github.com/daeho-ro) for another new feature!)

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.8](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v418-2022-09-08) which increases the default open file limit. ([#46](https://github.com/turbot/steampipe-plugin-csv/pull/46))

## v0.4.1 [2022-10-19]

_Bug fixes_

- Fixed the plugin to skip the non-standard Byte Order Mark (BOM) at the start of some CSV files (for instance, UTF-8 encoded CSV files from Excel) which would make the first column of such CSV files not queryable.

## v0.4.0 [2022-09-28]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#34](https://github.com/turbot/steampipe-plugin-csv/pull/34))
- Recompiled plugin with Go version `1.19`. ([#34](https://github.com/turbot/steampipe-plugin-csv/pull/34))

## v0.3.2 [2022-06-23]

_Bug fixes_

- Fixed an issue in the `{csv_filename}` tables where values would always be `null` for header columns with a period in their name. ([#25](https://github.com/turbot/steampipe-plugin-csv/pull/25))

## v0.3.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#22](https://github.com/turbot/steampipe-plugin-csv/pull/22))

## v0.3.0 [2022-04-27]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#16](https://github.com/turbot/steampipe-plugin-csv/pull/16))
- Added support for native Linux ARM and Mac M1 builds. ([#15](https://github.com/turbot/steampipe-plugin-csv/pull/15))

## v0.2.0 [2022-02-14]

_What's new?_

- File loading and matching through the `paths` argument has been updated to make the plugin easier to use:
  - The `paths` argument is no longer commented out by default for new plugin installations and now defaults to the current working directory
  - Recursive directory searching (`**`) is now supported
- Previously, when using wildcard matching (`*`), non-CSV files were automatically excluded to prevent parsing errors. These files are no longer automatically excluded to allow for a wider range of matches. If your current configuration uses wildcard matching, e.g., `paths = [ "/path/to/my/files/*" ]`, please update it to include the file extension, e.g., `paths = [ "/path/to/my/files/*.csv" ]`.

## v0.1.0 [2021-12-08]

_Enhancements_

- Recompiled plugin with Go version 1.17 ([#9](https://github.com/turbot/steampipe-plugin-csv/pull/9))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#8](https://github.com/turbot/steampipe-plugin-csv/pull/8))

## v0.0.4 [2021-11-16]

_Enhancements_

- Updated: Recompiled plugin with [steampipe-plugin-sdk v1.7.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v173--2021-11-08)

## v0.0.3 [2021-10-18]

_Bug fixes_

- Fixed: Plugin description in docs/index.md is now correct

## v0.0.2 [2021-10-18]

_Enhancements_

- Updated: Added more usage information to `{csv_filename}` table document

_Bug fixes_

- Fixed: Brand colour is now correct

## v0.0.1 [2021-10-18]

_What's new?_

- New tables added
  - [{csv_filename}](https://hub.steampipe.io/plugins/turbot/csv/tables/{csv_filename})
