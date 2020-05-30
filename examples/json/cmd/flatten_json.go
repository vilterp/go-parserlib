package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vilterp/go-parserlib/examples/json"
	parserlib "github.com/vilterp/go-parserlib/pkg"
)

func main() {
	parser, err := parserlib.NewStreamingParser(os.Stdin, json.Grammar, "value")
	if err != nil {
		log.Fatal(err)
	}
	for {
		evt, err := parser.NextEvent()
		if err == io.EOF {
			break
		}
		fmt.Printf("evt %#v\n", evt)
	}
}
