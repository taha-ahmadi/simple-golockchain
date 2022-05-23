package golockchain

import (
	"bytes"
	"fmt"
	"time"
)

type Block struct {
	Timestamp    time.Time `json:"timestamp"`
	Transactions []*Transaction

	Nonce    int32  `json:"nonce"`
	PrevHash []byte `json:"prev_hash"`
	Hash     []byte `json:"hash"`
}

// NewBlock returns a new Block with mask that is for
// difficulty level of Block
func NewBlock(txs []*Transaction, mask, prefHash []byte) *Block {
	b := Block{
		Timestamp:    time.Now(),
		Transactions: txs,
		PrevHash:     prefHash,
	}
	b.Hash, b.Nonce = DifficultHash(mask, b.Timestamp.UnixNano(), calculateTxHash(b.Transactions...), b.PrevHash)

	return &b
}

// Validate try to validate the current block with mask for validating
// the hash difficulty
func (b *Block) Validate(mask []byte) error {
	hash := GenerateHash(b.Timestamp.UnixNano(), calculateTxHash(b.Transactions...), b.PrevHash, b.Nonce)

	if !bytes.Equal(hash, b.Hash) {
		return fmt.Errorf("the hash is invalid it should %x but it is %x", hash, b.Hash)
	}

	if !ValidDifficulty(mask, hash) {
		return fmt.Errorf("hash is not good enough with mask %x", mask)
	}

	return nil
}

func (b *Block) String() string {
	return fmt.Sprintf(
		"Time: %s\nData: %d\nHash: %x\nPervHash: %x\nNonce:%d\n",
		b.Timestamp, len(b.Transactions), b.Hash, b.PrevHash, b.Nonce,
	)
}
