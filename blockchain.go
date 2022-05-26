package golockchain

import (
	"encoding/hex"
	"errors"
	"fmt"
)

type BlockChain struct {
	Difficulty int
	Mask       []byte

	store Store
}

// Add a block to the BlockChain
func (bc *BlockChain) Add(data ...*Transaction) (*Block, error) {
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

func (bc *BlockChain) UnspentTxn(address []byte) (map[string]*Transaction, map[string][]int, int, error) {
	spent := make(map[string][]int)
	txom := make(map[string][]int)
	txns := make(map[string]*Transaction)
	acc := 0
	err := Iterate(bc.store, func(b *Block) error {
		for _, txn := range b.Transactions {
			txnID := hex.EncodeToString(txn.ID)

			for i := range txn.VOut {
				if txn.VOut[i].TryUnlock(address) && !inArray(i, spent[txnID]) {
					txns[txnID] = txn
					txom[txnID] = append(txom[txnID], i)
					acc += txn.VOut[i].Value
				}
			}

			delete(spent, txnID)

			if txn.IsCoinBase() {
				continue
			}

			for i := range txn.VIn {
				if txn.VIn[i].MatchLock(address) {
					outID := hex.EncodeToString(txn.VIn[i].TXID)
					spent[outID] = append(spent[outID], txn.VIn[i].VOut)
				}
			}

		}

		return nil
	})
	if err != nil {
		return nil, nil, 0, fmt.Errorf("iterate error: %w", err)
	}

	return txns, txom, acc, nil
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
	gbTx := NewCoinBaseTx(genesis, nil)
	genesisBlock := NewBlock(
		[]*Transaction{gbTx},
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
