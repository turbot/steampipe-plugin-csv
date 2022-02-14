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
