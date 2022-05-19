package golockchain

import "fmt"

var (
	ErrNoInitialized = fmt.Errorf("blockchain not initialized yet")
)

// Store is an interface for blockchain storage drivers.
type Store interface {
	// Load should return the block from the store based on the requested hash
	Load(hash []byte) (*Block, error)
	Append(*Block) (*Block, error)

	// LastHash returns the last hash in the storage, if there is no block
	// it returns the ErrNotInitialized
	LastHash() ([]byte, error)
}

// Iterate over the blocks in the Store, if the callback returns an error it
// stops the loop and return the error to the caller
func Iterate(store Store, fn func(b *Block) error) error {
	last, err := store.LastHash()
	if err != nil {
		return err
	}

	for {
		block, err := store.Load(last)
		if err != nil {
			return err
		}
		if err := fn(block); err != nil {
			return err
		}

		if len(block.PrevHash) == 0 {
			return nil
		}

		last = block.PrevHash
	}
}
