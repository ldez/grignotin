# Grignotin

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/ldez/grignotin.svg)](https://github.com/ldez/grignotin/releases)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ldez/grignotin)](https://pkg.go.dev/github.com/ldez/grignotin)
[![Build Status](https://github.com/ldez/grignotin/actions/workflows/ci.yml/badge.svg)](https://github.com/ldez/grignotin/actions)

A collection of small helpers around Go proxy, Go meta information, etc.

## goproxy

A small Go proxy client to get information about a module from a Go proxy.

```go
package main

import (
	"fmt"

	"github.com/ldez/grignotin/goproxy"
)

func main() {
	client := goproxy.NewClient("")

	versions, err := client.GetVersions("github.com/ldez/grignotin")
	if err != nil {
		panic(err)
	}

	fmt.Println(versions)
}
```

## metago

A small lib to fetch meta information (`go-import`, `go-source`) for a module.

```go
package main

import (
	"fmt"

	"github.com/ldez/grignotin/metago"
)

func main() {
	meta, err := metago.Get("github.com/ldez/grignotin")
	if err != nil {
		panic(err)
	}

	fmt.Println(meta)
}
```

## Version

Gets information about releases and build. 

```go
package main

import (
	"fmt"

	"github.com/ldez/grignotin/version"
)

func main() {
    includeAll := false
	releases, err := version.GetReleases(includeAll)
	if err != nil {
		panic(err)
	}

	fmt.Println(releases)
}
```

```go
package main

import (
	"fmt"

	"github.com/ldez/grignotin/version"
)

func main() {
	build, err := version.GetBuild()
	if err != nil {
		panic(err)
	}

	fmt.Println(build)
}
```

## SumDB

- I recommend using the package [sumdb](https://pkg.go.dev/golang.org/x/mod/sumdb?tab=doc)
