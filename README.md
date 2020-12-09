# createsend

![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/xitonix/createsend?include_prereleases)
[![Build Status](https://travis-ci.com/xitonix/createsend.svg?branch=master)](https://travis-ci.com/xitonix/createsend)
[![Go Report Card](https://goreportcard.com/badge/github.com/xitonix/createsend)](https://goreportcard.com/report/github.com/xitonix/createsend)
[![codecov](https://codecov.io/gh/xitonix/createsend/branch/master/graph/badge.svg)](https://codecov.io/gh/xitonix/createsend)
[![GitHub license](https://img.shields.io/github/license/xitonix/createsend)](https://github.com/xitonix/createsend/blob/master/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/xitonix/createsend)](https://github.com/xitonix/createsend/issues)


Campaign Monitor API wrapper in Go

## Installation

```shell script
go get github.com/xitonix/createsend
```

## Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/xitonix/createsend"
)

func main() {
    client, err := createsend.New(createsend.WithAPIKey("[Your API Key]"))
    if err != nil {
        log.Fatal(err)
    }
    
    clients, err := client.Accounts().Clients()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, client := range clients {
        fmt.Printf("%s: %s\n", client.ID, client.Name)
    }
}
```

You can also use oAuth authentication token:

```go
client, err := createsend.New(createsend.WithOAuthToken("[OAuth Token]"))
if err != nil {
  log.Fatal(err)
}
```

## Contribution Guideline

The guideline can be found [here](https://github.com/xitonix/createsend/blob/master/CONTRIBUTING.md). Thank you 🥇
