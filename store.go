package golockchain

import "fmt"

var (
	ErrNoInitialized = fmt.Errorf("blockchain not initialized yet")
)

type Store interface {
	Load(hash []byte) (*Block, error)
	Append(*Block) (*Block, error)

	LastHash() ([]byte, error)
}

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
