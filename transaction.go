package golockchain

import "time"

const (
	coinBaseReward = 1000000
)

type Transaction struct {
	ID []byte

	VOut []TXOutput
	VIn  []TXInput
}

type TXOutput struct {
	Value  int
	PubKey []byte
}

type TXInput struct {
	TXID []byte
	VOut int // index of Transaction.VOut
	Sig  []byte
}

func calculateTxID(tx *Transaction) []byte {
	return GenerateHash(tx.VOut, tx.VIn)
}

func calculateTxHash(txs ...*Transaction) []byte {
	data := make([]interface{}, len(txs))

	for i := range txs {
		data[i] = txs[i].ID
	}

	return GenerateHash(data...)
}

func NewCoinBaseTx(to, data []byte) *Transaction {
	if len(data) == 0 {
		data = GenerateHash(to, time.Now())
	}

	txi := TXInput{
		TXID: []byte{},
		VOut: -1,
		Sig:  data,
	}

	txo := TXOutput{
		Value:  coinBaseReward,
		PubKey: to,
	}

	tx := &Transaction{
		VOut: []TXOutput{txo},
		VIn:  []TXInput{txi},
	}
	tx.ID = calculateTxID(tx)

	return tx
}
