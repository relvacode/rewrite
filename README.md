# Rewrite

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