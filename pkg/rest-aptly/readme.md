# rest-aptly

WIP

## Currently not implemented

* mirror API
* repo import/copy/move/remove/search
* snapshot verify/pull/filter
* db API
* task API

and probably some more options  

## Usage

```golang
package main

import (
    "fmt"
    aptly "raptly/pkg/rest-aptly"
)

func main() {
    client := aptly.NewClient("http://localhost:8080")

    version, err := client.Version()
    if err != nil {
        panic(err)
    }
    fmt.Println(version.Version)
}

type Context struct {
    client *aptly.Client
}
```
