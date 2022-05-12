package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	block := NewBlockChain(7)
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()
	block.Add([]byte("test"))

	if err := block.Validate(); err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(block)
}
