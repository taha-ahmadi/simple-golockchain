package main

import (
	"flag"
	golockchain "github.com/taha-ahmadi/simple-golockchain"
	"log"
	"os"
)

const (
	difficulty = 2
)

func main() {
	var store string
	flag.StringVar(&store, "store", os.Getenv("BC_STORE"), "The storage to use")
	flag.Usage = usage
	flag.Parse()

	s := golockchain.NewFolderStore(store)
	if err := dispatch(s, flag.Args()...); err != nil {
		log.Fatal(err.Error())
	}
}
