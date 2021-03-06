package golockchain

import (
	"bytes"
	"fmt"
)

type mapStore struct {
	data map[string]*Block
	last []byte
}

func (ms *mapStore) Load(hash []byte) (*Block, error) {
	hashBase16 := fmt.Sprintf("%x", hash)
	if block, ok := ms.data[hashBase16]; ok {
		return block, nil
	}

	return nil, fmt.Errorf("block is not in this store")
}

func (ms *mapStore) Append(block *Block) (*Block, error) {
	hashBase16 := fmt.Sprintf("%x", block.Hash)

	if !bytes.Equal(ms.last, block.PrevHash) {
		return nil, fmt.Errorf("store is out of sync")
	}

	if _, ok := ms.data[hashBase16]; ok {
		return nil, fmt.Errorf("duplicate")
	}

	ms.data[hashBase16] = block
	ms.last = block.Hash
	return ms.data[hashBase16], nil
}

func (ms *mapStore) LastHash() ([]byte, error) {
	if len(ms.last) == 0 {
		return nil, ErrNoInitialized
	}

	return ms.last, nil
}

// NewMapStore is a in memory storage for blockchain it's fast and
// it is good for testing purposes
func NewMapStore() Store {
	return &mapStore{
		data: make(map[string]*Block),
	}
}
