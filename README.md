# Pipeline Parser

[![Test Pipeline Parser](https://github.com/argonsecurity/pipeline-parser/actions/workflows/test.yml/badge.svg)](https://github.com/argonsecurity/pipeline-parser/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://github.com/argonsecurity/pipeline-parser/blob/main/LICENSE)
[![go-report-card][go-report-card]](https://goreportcard.com/report/github.com/argonsecurity/pipeline-parser)

[go-report-card]: https://goreportcard.com/badge/github.com/argonsecurity/pipeline-parser?style=flat-square

## Description

Pipeline Parser is Argon's solution for parsing and analyzing pipeline files of popular CI yaml files in order to create a generic pipeline entity that can be used across platforms.

#### Supported Platforms:

| Platform
| :---:
| Github Workflows
| Gitlab CI

## Local Development

First, execute the following command to enable the client's git hooks:

```
git config core.hooksPath .githooks
```
