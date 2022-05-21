package golockchain

import (
	"errors"
	"fmt"
)

type BlockChain struct {
	Difficulty int
	Mask       []byte

	store Store
}

// Add a block to the BlockChain
func (bc *BlockChain) Add(data []byte) (*Block, error) {
	lastBlockHash, err := bc.store.LastHash()
	if err != nil {
		return nil, fmt.Errorf("failed getting last block hash %w", err)
	}
	block, err := bc.store.Append(NewBlock(data, bc.Mask, lastBlockHash))
	if err != nil {
		return nil, fmt.Errorf("failed to append block: %w", err)
	}

	return block, nil
}

// Print the current Blockchain to Stdout, it is good for testing purposes
func (bc *BlockChain) Print(header bool, count int) error {
	fmt.Printf("Difficulty: %d\n store: %T\n", bc.Difficulty, bc.store)
	if header {
		return nil
	}

	errEnough := fmt.Errorf("enough")
	
	err := Iterate(bc.store, func(b *Block) error {
		if count > 0 {
			count--
		}
		fmt.Print(b)
		if count == 0 {

		}

		return nil
	})

	if errors.Is(err, errEnough) {
		return nil
	}

	return err
}

// NewBlockChain returns a new BlockChain
func NewBlockChain(genesis []byte, difficulty int, store Store) (*BlockChain, error) {
	mask := GenerateMask(difficulty)
	blockChain := BlockChain{
		Difficulty: difficulty,
		Mask:       mask,
		store:      store,
	}

	_, err := store.LastHash()

	if !errors.Is(err, ErrNoInitialized) {
		return nil, fmt.Errorf("getting the last hash failed: %w", err)
	}

	genesisBlock := NewBlock(
		genesis,
		blockChain.Mask,
		[]byte{},
	)
	if _, err := store.Append(genesisBlock); err != nil {
		return nil, err
	}
	return &blockChain, nil
}

// OpenBlockChain open blockchain
func OpenBlockChain(difficulty int, store Store) (*BlockChain, error) {
	mask := GenerateMask(difficulty)
	blockChain := BlockChain{
		Difficulty: difficulty,
		Mask:       mask,
		store:      store,
	}

	_, err := store.LastHash()
	if err != nil {
		return nil, fmt.Errorf("blockchain is not valid for getting LastHash: %w", err)
	}

	return &blockChain, nil
}

// Validate checks the PrevHash and hash of the blocks
func (bc *BlockChain) Validate() error {
	return Iterate(bc.store, func(b *Block) error {
		if err := b.Validate(bc.Mask); err != nil {
			return fmt.Errorf("block is invalid: %w", err)
		}
		return nil
	})
}
