package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"ecosia/tree"
)

var conf struct {
	listenAddr string
	urlKey     string
}

func init() {
	conf.listenAddr = *flag.String("addr", ":8000", "listen address")
	conf.urlKey = *flag.String("url-key", "favoriteTree", "url query key for the input")
	flag.Parse()
}

func main() {
	http.Handle("/", tree.NewHandler(conf.urlKey, tree.NewTemplateView(nil)))
	fmt.Printf("listening on %s for query key '%s'\n\n", conf.listenAddr, conf.urlKey)
	log.Fatalf("server crashed: %s", http.ListenAndServe(conf.listenAddr, nil))
}
