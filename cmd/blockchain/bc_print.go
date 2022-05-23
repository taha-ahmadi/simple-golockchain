package main

import (
	"flag"
	"fmt"
	golockchain "github.com/taha-ahmadi/simple-golockchain"
)

func printBlockchain(store golockchain.Store, args ...string) error {

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		header bool
		count  int
	)
	fs.BoolVar(&header, "header", false, "only print header")
	fs.IntVar(&count, "count", -1, "how many records to show")

	fs.Parse(args[1:])

	bc, err := golockchain.OpenBlockChain(difficulty, store)
	if err != nil {
		return fmt.Errorf("open blockchain failed: %w", err)
	}

	return bc.Print(header, count)
}

func init() {
	_ = addCommand("print", "print the blockchain", printBlockchain)
}
