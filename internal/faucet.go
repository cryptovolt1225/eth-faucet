package internal

import (
	"log"
	"math/big"
)

type faucet struct {
	payout    *big.Int
	queue     chan string
	txBuilder ITxBuilder
}

func NewFaucet(capacity int, builder ITxBuilder) *faucet {
	return &faucet{
		queue:     make(chan string, capacity),
		txBuilder: builder,
	}
}

func (f *faucet) TryEnqueue(job string) bool {
	select {
	case f.queue <- job:
		return true
	default:
		return false
	}
}

func (f faucet) GetPayoutWei() *big.Int {
	return f.payout
}

func (f *faucet) SetPayoutEther(amount int64) {
	payoutWei := new(big.Int).Mul(big.NewInt(amount), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	f.payout = payoutWei
}

func (f faucet) transferFund(to string) (string, error) {
	tx, err := f.txBuilder.buildUnsignedTx(to, f.payout, nil)
	if err != nil {
		return "", err
	}

	if err := f.txBuilder.submitSignedTx(tx); err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func (f *faucet) Run() {
	for address := range f.queue {
		txHash, err := f.transferFund(address)
		if err != nil {
			log.Println(err)
		}
		log.Println(txHash)
	}
}

func (f *faucet) Close() {
	close(f.queue)
}
