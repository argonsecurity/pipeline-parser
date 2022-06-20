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
| GitHub Workflows
| GitLab CI

## Usage

### Package Usage

```golang
import "github/argonsecurity/pipeline-parser/pkg/handler"

// Read the pipeline data as bytes array
buf, err := ioutil.ReadFile("/path/to/workflow.yml")
if err != nil {
    return nil
}

// Parse the pipeline from the specific platform to the common pipeline object
pipeline, err := handler.Handle(buf, consts.GitHubPlatform)
```

### CLI Usage

#### Parse GitHub Workflow yaml

```bash
pipeline-parser -p github workflow.yml
```

#### Parse GitLab CI yaml

```bash
pipeline-parser -p gitlab .gitlab-ci.yml
```

#### Parse multiple files in one execution

```bash
pipeline-parser -p github workflow-1.yml workflow-2.yml workflow-3.yml
```

## Local Development

First, execute the following command to enable the client's git hooks:

```
git config core.hooksPath .githooks
```
