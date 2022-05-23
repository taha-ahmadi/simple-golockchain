package golockchain

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestBlockCreation(t *testing.T) {
	data := []*Transaction{NewCoinBaseTx([]byte("bita"), nil)}
	mask := GenerateMask(2)
	prev := GenerateHash("Prev hash")

	b := NewBlock(data, mask, prev)
	if err := b.Validate(mask); err != nil {
		t.Errorf("Validation failed: %q", err)
		return
	}

	if !strings.Contains(b.String(), hex.EncodeToString(prev)) {
		t.Errorf("The prev hash is not in string")
	}

	mask2 := GenerateMask(8)
	if err := b.Validate(mask2); err == nil {
		t.Errorf("Block should be invalid with mask %x, but is not", mask2)
	}

	// Invalidate block
	b.Nonce++
	if err := b.Validate(mask); err == nil {
		t.Errorf("Block should be invalid, but is not")
	}
	b.Nonce--

	b.Transactions = []*Transaction{NewCoinBaseTx([]byte("forud"), nil)}
	if err := b.Validate(mask); err == nil {
		t.Errorf("Block should be invalid, but is not")
	}
	b.Transactions = data

	b.PrevHash = GenerateHash("Something else")
	if err := b.Validate(mask); err == nil {
		t.Errorf("Block should be invalid, but is not")
	}
	b.PrevHash = prev

	hash := b.Hash
	b.Hash, _ = DifficultHash(mask, "Something else")
	if err := b.Validate(mask); err == nil {
		t.Errorf("Block should be invalid, but is not")
	}
	b.Hash = hash
}
