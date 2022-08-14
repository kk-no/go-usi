# go-usi

This repository is a library for using the USI protocol from Go.
For more information on the USI protocol, please see.
- http://shogidokoro.starfree.jp/usi.html
- http://shogidokoro.starfree.jp/usi.html#GameExample

## Example

main.go
```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/kk-no/go-usi"
)

func main() {
	ctx := context.Background()

	enginePath := os.Getenv("ENGINE_PATH")
	engine, err := usi.New(enginePath)
	if err != nil {
		log.Fatalln(err)
	}

	if err := engine.Connect(ctx); err != nil {
		os.Exit(0)
	}
}
```

shell
```shell
$ go run main.go
$ isready
$ go infinite
```