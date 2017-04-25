# google-taxonomy-go

## Buildtools

To download taxonomy txt:
```shell
buildtools/dl-taxonomy-txt.sh
```

To clean all generated files:
```shell
buildtools/clean.sh
```

To generate data.go:
```shell
go get -u github.com/jteeuwen/go-bindata/...
./buildtools/build.sh
```

## Data files consideration

I preseve those file in git repo to track the changes when google updates taxonomy layouts.

## Usage

See examples/main.go

```go
package main

import (
	"fmt"

	"github.com/go-google-taxonomy/google-taxonomy/taxonomy"
)

const keyLanguage = "en-US"

func main() {
	tx, err := taxonomy.NewTaxonomy(keyLanguage, []string{})
	if err != nil {
		panic(err)
	}
	infs, err := tx.GetRootsCategoryInfo(keyLanguage)
	if err != nil {
		panic(err)
	}
	for _, inf := range infs {
		fmt.Printf("%d - %s\n", inf.ID, inf.String())
	}
	inf, err := tx.GetCategoryInfo(16, keyLanguage)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d - %s\n", inf.ID, inf.String())
}

```

## Testing

```shell
go test ./taxonomy
# Or run data race test
# go test -race ./taxonomy
```

We do not test `zh-CN` for category language compatibility, as the version of `zh-CN` is
very behind the other languages.