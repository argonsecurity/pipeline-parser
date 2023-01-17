# Pipeline Parser

[![Test Pipeline Parser](https://github.com/argonsecurity/pipeline-parser/actions/workflows/test.yml/badge.svg)](https://github.com/argonsecurity/pipeline-parser/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/argonsecurity/pipeline-parser/blob/main/LICENSE)
[![go-report-card][go-report-card]](https://goreportcard.com/report/github.com/argonsecurity/pipeline-parser)
![coverage report](https://img.shields.io/codecov/c/github/argonsecurity/pipeline-parser)

[go-report-card]: https://goreportcard.com/badge/github.com/argonsecurity/pipeline-parser

## Description

Pipeline Parser is Argon's solution for parsing and analyzing pipeline files of popular CI yaml files in order to create a generic pipeline entity that can be used across platforms.

#### Supported Platforms:

| Platform
| :---:
| GitHub Workflows
| GitLab CI
| Azure Pipelines
| Bitbucket Pipelines

## Usage

### Package Usage

```golang
import (
    "os"

    "github.com/argonsecurity/pipeline-parser/pkg/handler"
    "github.com/argonsecurity/pipeline-parser/pkg/consts"
)

// Read the pipeline data as bytes array
buf, err := os.ReadFile("/path/to/workflow.yml")
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

#### Parse Azure Pipelines yaml

```bash
pipeline-parser -p azure .azure-pipelines.yml
```

#### Parse Bitbucket Pipelines yaml

```bash
pipeline-parser -p bitbucket .bitbucket-pipelines.yml
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
