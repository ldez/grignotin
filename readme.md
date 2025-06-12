# Grignotin

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/ldez/grignotin.svg)](https://github.com/ldez/grignotin/releases)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ldez/grignotin)](https://pkg.go.dev/github.com/ldez/grignotin)
[![Build Status](https://github.com/ldez/grignotin/actions/workflows/ci.yml/badge.svg)](https://github.com/ldez/grignotin/actions)

A collection of small helpers around Go proxy, Go meta information, etc.

## goproxy

A small Go proxy client to get information about a module from a Go proxy.

<details><summary>Example</summary>

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

</details>

## metago

A small lib to fetch meta information (`go-import`, `go-source`) for a module.

<details><summary>Example</summary>

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

</details>

## Version

Gets information about releases and build. 

<details><summary>Examples</summary>

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

</details>

## SumDB

- I recommend using the package [sumdb](https://pkg.go.dev/golang.org/x/mod/sumdb?tab=doc)


## gomod

A set of functions to get information about module (`go list`/`go env`).

<details><summary>Examples</summary>

```go
package main

import (
	"context"
	"fmt"

	"github.com/ldez/grignotin/gomod"
)

func main() {
	info, err := gomod.GetModuleInfo(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(info)
}
```

```go
package main

import (
	"context"
	"fmt"

	"github.com/ldez/grignotin/gomod"
)

func main() {
	modpath, err := gomod.GetModulePath(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(modpath)
}

```

</details>

## goenv

A set of functions to get information from `go env`.

<details><summary>Examples</summary>

```go
package main

import (
	"context"
	"fmt"

	"github.com/ldez/grignotin/goenv"
)

func main() {
	value, err := goenv.GetOne(context.Background(), goenv.GOMOD)
	if err != nil {
		panic(err)
	}

	fmt.Println(value)
}
```

```go
package main

import (
	"context"
	"fmt"

	"github.com/ldez/grignotin/goenv"
)

func main() {
	values, err := goenv.Get(context.Background(), goenv.GOMOD, goenv.GOMODCACHE)
	if err != nil {
		panic(err)
	}

	fmt.Println(values)
}
```

```go
package main

import (
	"context"
	"fmt"

	"github.com/ldez/grignotin/goenv"
)

func main() {
	values, err := goenv.GetAll(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(values)
}
```

</details>
