package main

import (
	"flag"
	"fmt"
	golockchain "github.com/taha-ahmadi/simple-golockchain"
)

func initialize(store golockchain.Store, args ...string) error {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		genesis string
	)
	fs.StringVar(&genesis, "genesis", "Genesis block", "Genesis data to use")

	fs.Parse(args[1:])

	_, err := golockchain.NewBlockChain([]byte(genesis), difficulty, store)
	if err != nil {
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

func init() {
	addCommand("init", "Create an empty blockchain", initialize)
}
