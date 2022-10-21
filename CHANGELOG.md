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
