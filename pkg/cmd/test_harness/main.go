package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/vilterp/go-parserlib/examples/treesql"
	"github.com/vilterp/go-parserlib/pkg/test_harness"
)

var port = flag.String("port", "9999", "port to listen on")
var host = flag.String("host", "0.0.0.0", "host to listen on")

func main() {
	flag.Parse()

	server := test_harness.NewServer(treesql.MakeLanguage(treesql.BlogSchema))

	addr := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("serving on http://%s/", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatal(err)
	}
}
