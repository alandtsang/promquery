# promquery

Query library for prometheus metrics

## Document

## Usage

```go
package main

import (
    "log"

    "github.com/alandtsang/promquery/pkg/promquery"
    v1 "github.com/alandtsang/promquery/pkg/v1"
)

func main() {
    addr := "http://prometheus.zxl.test.com:32219"
    node := "docker-desktop"
    queryRange := "1h"
    namespace := "default"
    pod := "nginx-6cc5b96679-ts7lq"
    container := "nginx"

    query, err := promquery.New(addr)
    if err != nil {
        log.Fatal(err)
    }

    params := promquery.BuildQueryParams(node, namespace, pod, container, queryRange)
    v1.CPUQueryRange(query, params)
}
```

## LICENCE
[LICENCE](https://raw.githubusercontent.com/alandtsang/promquery/main/LICENSE)