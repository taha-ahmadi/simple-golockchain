package main

import (
	"fmt"
	"log"
	"time"

	sg "github.com/taha-ahmadi/simple-golockchain"
)

func main() {
	blockchain, err := sg.NewBlockChain(4, sg.NewMapStore())
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()
	blockchain.Add([]byte("test"))

	if err := blockchain.Validate(); err != nil {
		log.Fatalf(err.Error())
	}
	blockchain.Print()
}
