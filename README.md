[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/asachs01/sensu-go-uptime-checks)
[![Build Status](https://travis-ci.org/asachs01/sensu-go-uptime-checks.svg?branch=master)](https://travis-ci.org/asachs01/sensu-go-uptime-checks)

## Sensu Go Uptime Checks Plugin

- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Sensu Go](#sensu-go)
    - [Asset registration](#asset-registration)
    - [Asset definition](#asset-definition)
    - [Handler definition](#handler-definition)
  - [Sensu Core](#sensu-core)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

### Overview

This plugin provides a check for system uptime for Sensu Go. The `sensu-go-uptime-check` check takes the flags `-w` (warning) and `-c` (critical) and a time duration after each flag (e.g. 24h). The time values can be represented in seconds (s), minutes (m), or hours (h). The plugin uses a default warning value of 72h and a default critical value of 168h (1 week).

### Files

N/A

## Usage examples

### Help

```
The Sensu Go check for system uptime

Usage:
  sensu-go-uptime-check [flags]

Flags:
  -w, --warning (time in s,m,h)   Warning value in seconds, minutes, or hours, default is 72 hours (72h)
  -c, --critical (time in s,m,h)   Warning value in seconds, minutes, or hours default is 1 week (168h)
  -h, --help         help for sensu-go-uptime-status
```

## Configuration
### Sensu Go
#### Asset registration

Assets are the best way to make use of this plugin. If you're not using an asset, please consider doing so! If you're using sensuctl 5.13 or later, you can use the following command to add the asset: 

`sensuctl asset add sensu-plugins/sensu-go-uptime-checks`

If you're using an earlier version of sensuctl, you can download the asset definition from [this project's Bonsai asset index page][2], download a copy of the handler plugin from [releases][1], or create an executable script from this source.

From the local path of the sensu-go-uptime-checks repository:

```
go build -o /usr/local/bin/sensu-go-uptime-checks main.go
```

#### Asset definition

```yaml
---
type: Asset
api_version: core/v2
metadata:
  name: sensu-go-uptime-checks
spec:
  url: https://assets.bonsai.sensu.io/1f967e65880ead0b4bfe13e331b3fc6f26ebfed2/sensu-go-uptime-checks_1.0.1_linux_amd64.tar.gz
  sha512: 3d732f21611bb03dddc529c46b1bffd2b97f2a5b5ae7a935964a179adb26d5ccffa3ffaf0662380be4485d6d5e4295b0812ef1d22919e02ebb412c4eef1aff24
```

**NOTE**: Make sure to update the URL and SHA512 before you use the asset. If you don't, you might be using an older version.

#### Handler definition

```yaml
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-go-uptime-check
  namespace: CHANGEME
spec:
  command: sensu-go-uptime-check
  runtime_assets:
  - sensu-go-uptime-checks
  interval: 60
  publish: true
  output_metric_format: nagios_perfdata
  output_metric_handlers:
  - infuxdb
  handlers:
  - slack
  subscriptions:
  - system
```

### Sensu Core

N/A

## Installation from source

### Sensu Go

See the instructions above for [asset registration](#asset-registration).

### Sensu Core

Install and setup plugins on [Sensu Core](https://docs.sensu.io/sensu-core/latest/installation/installing-plugins/).

## Additional notes

### Supported Operating Systems

This project uses `gopsutil` and thus depends on the systems that it supports. For this plugin, the following operating systems are supported:

* Linux
* FreeBSD
* OpenBSD
* Mac OS X
* Windows (states not supported, but I've confirmed that it is)
* Solaris

### Example output

![screenshot_of_check_result](http://share.sachshaus.net/ddbeec586345/Screen%252520Shot%2525202019-07-29%252520at%25252011.05.48%252520PM.png)

## Contributing


See the [Sensu Go repository CONTRIBUTING.md][3] for information about contributing to this plugin. 

[1]: https://github.com/asachs01/sensu-go-uptime-checks/releases
[2]: https://bonsai.sensu.io/assets/asachs01/sensu-go-uptime-checks
[3]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
