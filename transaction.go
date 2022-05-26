package golockchain

import (
	"bytes"
	"fmt"
	"time"
)

const (
	coinBaseReward = 1000000
)

type Transaction struct {
	ID []byte

	VOut []TXOutput
	VIn  []TXInput
}

// TXOutput we can spend outputs
type TXOutput struct {
	Value  int
	PubKey []byte // anyone who has the key can spend TXOutput
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

func (txi *TXInput) MatchLock(key []byte) bool {
	return bytes.Equal(txi.Sig, key)
}

func (txo *TXOutput) TryUnlock(key []byte) bool {
	return bytes.Equal(txo.PubKey, key)
}

func (txn *Transaction) IsCoinBase() bool {
	return len(txn.VOut) == 1 &&
		len(txn.VIn) == 1 &&
		txn.VIn[0].VOut == -1 &&
		len(txn.VIn[0].TXID) == 0
}

func NewTransaction(bc *BlockChain, from, to []byte, amount int) (*Transaction, error) {
	txns, txom, acc, err := bc.UnspentTxn(from)
	if err != nil {
		return nil, fmt.Errorf("get unused txn failed: %w", err)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("negative transfer?")
	}

	if acc < amount {
		return nil, fmt.Errorf("not enough money, want %d have %d", amount, acc)
	}

	var (
		vin      []TXInput
		required = amount
	)

bigLoop:
	for id, txn := range txns {
		for _, v := range txom[id] {
			required -= txn.VOut[v].Value
			vin = append(vin, TXInput{
				TXID: txn.ID,
				VOut: v,
				Sig:  from, // TODO : real sign
			})

			if required <= 0 {
				break bigLoop
			}
		}
	}

	vout := []TXOutput{
		TXOutput{
			Value:  amount,
			PubKey: to,
		},
	}
	if required < 0 {
		vout = append(vout, TXOutput{
			Value:  -required,
			PubKey: from,
		})
	}

	txn := &Transaction{
		VIn:  vin,
		VOut: vout,
	}

	txn.ID = calculateTxID(txn)

	return txn, nil
}
