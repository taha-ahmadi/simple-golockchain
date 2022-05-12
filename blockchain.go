package main

import (
	"bytes"
	"fmt"
)

type BlockChain struct {
	Difficulty int
	Mask       []byte
	Block      []*Block
}

// Add a block to the BlockChain
func (bc *BlockChain) Add(data []byte) error {
	lenBlockchain := len(bc.Block)

	if lenBlockchain == 0 {
		return fmt.Errorf("you should have a block chain first")
	}

	bc.Block = append(bc.Block, NewBlock(data, bc.Mask, bc.Block[lenBlockchain-1].Hash))

	return nil
}

func (bc *BlockChain) String() string {
	var str string
	for _, v := range bc.Block {
		str += v.String()
	}

	return str
}

// NewBlockChain returns a new BlockChain
func NewBlockChain(difficulty int) *BlockChain {
	mask := GenerateMask(difficulty)
	blockChain := BlockChain{
		Difficulty: difficulty,
		Mask:       mask,
	}
	blockChain.Block = []*Block{
		NewBlock([]byte("Genesis Block"), blockChain.Mask,
			GenerateHash([]byte("Genesis Block"))),
	}

	return &blockChain
}

// Validate checks the PrevHash and hash of the blocks
func (bc *BlockChain) Validate() error {

	for i := range bc.Block {

		if err := bc.Block[i].Validate(bc.Mask); err != nil {
			return fmt.Errorf("blockchain is not valid: %v", err)
		}

		if i == 0 {
			continue
		}

		if !bytes.Equal(bc.Block[i].PrevHash, bc.Block[i-1].Hash) {
			return fmt.Errorf("the order is invalid, it should be %x but it is %x", bc.Block[i-1].Hash, bc.Block[i].PrevHash)
		}
	}

	return nil
}
