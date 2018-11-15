[![Build Status](https://travis-ci.org/relvacode/rewrite.svg?branch=master)](https://travis-ci.org/relvacode/rewrite) [![GoDoc](https://godoc.org/github.com/relvacode/rewrite?status.svg)](https://godoc.org/github.com/relvacode/rewrite)

Rewrite `io.Reader` streams


```
go get github.com/relvacode/rewrite
```

```go
package main

import (
    "github.com/relvacode/rewrite"
)

func main() {
    f, _ := os.Open("file.txt")
    r := rewrite.New(f, []byte("find this"), []byte("replace with this"))

    // data from r will be rewritten
}

```