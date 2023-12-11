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
engines: ["steampipe", "sqlite", "postgres", "export"]
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
  # Paths can be configured with a local directory, a remote Git repository URL, or an S3 bucket URL
  # Refer https://hub.steampipe.io/plugins/turbot/csv#supported-path-formats for more information
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
- `header` - Whether to use the first row as the header row when creating column names. Valid values are "auto", "on", "off". Defaults to "auto". For more information, please see [Column Names](https://hub.steampipe.io/plugins/turbot/csv/tables/{csv_filename}#column-names).

### Supported Path Formats

The `paths` config argument is flexible and can search for CSV files from several sources, e.g., local directory paths, Git, S3.

The following sources are supported:

- [Local files](#configuring-local-file-paths)
- [Remote Git repositories](#configuring-remote-git-repository-urls)
- [S3](#configuring-s3-urls)

Paths may [include wildcards](https://pkg.go.dev/path/filepath#Match) and support `**` for recursive matching. For example:

```hcl
connection "csv" {
  plugin = "csv"

  paths = [
    "*.csv",
    "**/*.csv",
    "github.com/codeforamerica/ohana-api//data/sample-csv//*.csv",
    "bitbucket.org/ved_misra/sample-csv//*.csv",
    "gitlab.com/subhajit7/example-files//sample-csv-files//*.csv",
    "s3::https://bucket.s3.us-east-1.amazonaws.com/test_folder//*.csv"
  ]
}
```

**Note**: If any path matches on `*` without `.csv`, all files (including non-CSV files) in the directory will be matched, which may cause errors if incompatible file types exist.

#### Configuring Local File Paths

You can define a list of local directory paths to search for CSV files. Paths are resolved relative to the current working directory. For example:

- "*.csv" matches all CSV files in the CWD
- "*.csv.gz" matches all gzipped CSV files in the CWD
- "**/*.csv" matches all CSV files in the CWD and all sub-directories
- "../*.csv" matches all CSV files in the CWD's parent directory
- "steampipe*.csv" matches all CSV files starting with "steampipe" in the CWD
- "/path/to/dir/*.csv" matches all CSV files in a specific directory
- "/path/to/dir/custom.csv" matches a specific file

```hcl
connection "csv" {
  plugin = "csv"

  paths = [ "*.csv", "~/*.csv", "/path/to/dir/custom.csv" ]
}
```

#### Configuring Remote Git Repository URLs

You can also configure `paths` with any Git remote repository URLs, e.g., GitHub, BitBucket, GitLab. The plugin will then attempt to retrieve any CSV files from the remote repositories.

For example:

- `github.com/codeforamerica/ohana-api//data/sample-csv//*.csv` matches CSV files in the specified repository.
- `github.com/codeforamerica/ohana-api//data//**/*.csv` matches all CSV files in the specified repository and all subdirectories.
- `github.com/codeforamerica/ohana-api//data//**/*.csv?ref=v1.2.0` matches all CSV files in the specific tag of a repository.
- `github.com/codeforamerica/ohana-api//data//**/*.csv` matches all CSV files in the specified folder path.

You can specify a subdirectory after a double-slash (`//`) if you want to download only a specific subdirectory from a downloaded directory.

```hcl
connection "csv" {
  plugin = "csv"

  paths = [ "github.com/codeforamerica/ohana-api//data/sample-csv//*.csv" ]
}
```

Similarly, you can define a list of GitLab and BitBucket URLs to search for CSV files:

```hcl
connection "csv" {
  plugin = "csv"

  paths = [
    "github.com/codeforamerica/ohana-api//data/sample-csv//*.csv",
    "bitbucket.org/ved_misra/sample-csv//*.csv",
    "gitlab.com/subhajit7/example-files//sample-csv-files//*.csv"
  ]
}
```

#### Configuring S3 URLs

You can also query all CSV files stored inside an S3 bucket (public or private) using the bucket URL.

##### Accessing a Private Bucket

In order to access your files in a private S3 bucket, you will need to configure your credentials. You can use your configured AWS profile from local `~/.aws/config`, or pass the credentials using the standard AWS environment variables, e.g., `AWS_PROFILE`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`.

We recommend using AWS profiles for authentication.

**Note:** Make sure that `region` is configured in the config. If not set in the config, `region` will be fetched from the standard environment variable `AWS_REGION`.

You can also authenticate your request by setting the AWS profile and region in `paths`. For example:

```hcl
connection "csv" {
  plugin = "csv"

  paths = [
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com//*.csv?aws_profile=<AWS_PROFILE>",
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com/test_folder//*.csv?aws_profile=<AWS_PROFILE>"
  ]
}
```

**Note:**

In order to access the bucket, the IAM user or role will require the following IAM permissions:

- `s3:ListBucket`
- `s3:GetObject`
- `s3:GetObjectVersion`

If the bucket is in another AWS account, the bucket policy will need to grant access to your user or role. For example:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ReadBucketObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:user/YOUR_USER"
      },
      "Action": ["s3:ListBucket", "s3:GetObject", "s3:GetObjectVersion"],
      "Resource": ["arn:aws:s3:::test-bucket1", "arn:aws:s3:::test-bucket1/*"]
    }
  ]
}
```

##### Accessing a Public Bucket

Public access granted to buckets and objects through ACLs and bucket policies allows any user access to data in the bucket. We do not recommend making S3 buckets public, but if there are specific objects you'd like to make public, please see [How can I grant public read access to some objects in my Amazon S3 bucket?](https://aws.amazon.com/premiumsupport/knowledge-center/read-access-objects-s3-bucket/).

You can query any public S3 bucket directly using the URL without passing credentials. For example:

```hcl
connection "config" {
  plugin = "config"

  paths = [
    "s3::https://bucket-1.s3.us-east-1.amazonaws.com/test_folder//*.csv",
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com/test_folder//**/*.csv"
  ]
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-csv
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
