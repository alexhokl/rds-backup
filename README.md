# AWS RDS MSSQL Backup Manager [![Build Status](https://travis-ci.org/alexhokl/rds-backup.svg?branch=master)](https://travis-ci.org/alexhokl/rds-backup) [![Coverage Status](https://img.shields.io/coveralls/alexhokl/rds-backup.svg)](https://coveralls.io/r/alexhokl/rds-backup?branch=master)

A command line tool to manage SQL backups on AWS RDS

##### Prerequisites

- [Sqlcmd](https://docs.microsoft.com/en-us/sql/tools/sqlcmd-utility) or [Docker](https://www.docker.com/) installed
- [AWS CLI](https://aws.amazon.com/cli/) installed and configured

##### Download

- Feel free to download the latest version from [release page](https://github.com/alexhokl/rds-backup/releases), or
- use `go get -u github.com/alexhokl/rds-backup` if you have Go installed

##### Examples

###### To create a backup

```sh
rds-backup create -w --bucket your-s3-bucket-name --database your-database-name --password your-database-password --server your-rds-server --username your-rds-sql-server-login --filename filename-on-s3.bak
```

###### To create a backup and restore in a Docker container on your local machine

```sh
rds-backup create -r --bucket your-s3-bucket-name --database your-database-name --password your-database-password --server your-rds-server --username your-rds-sql-server-login --filename filename-on-s3.bak --container your-container-name --restore-password your-container-sql-password
```

###### To create a backup and restore in a native MSSQL server on your local machine

```sh
rds-backup create -r -n --bucket your-s3-bucket-name --database your-database-name --password your-database-password --server your-rds-server --username your-rds-sql-server-login --filename filename-on-s3.bak --restore-password your-container-sql-password
```

##### Tricks

You can avoid specifying some of the parameters every time by using a configuration file or environment variables or a combination of both.

###### Configuration file

```yaml
bucket: your-s3-bucket-name
database: your-database-name
password: your-database-password
server: your-rds-server
username: your-rds-sql-server-login
filename: filename-on-s3.bak
```

###### Environment variables

```sh
export bucket=your-s3-bucket-name
export database=your-database-name
export password=your-database-password
export server=your-rds-server
export username=your-rds-sql-server-login
export filename=filename-on-s3.bak
```

