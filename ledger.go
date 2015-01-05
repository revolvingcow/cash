package main

import (
	"fmt"
	"math/big"
	"sort"
)

type Ledger struct {
	Transactions []Transaction
}

func (l *Ledger) ToString() string {
	balance := make(map[string]*big.Rat)

	for _, t := range l.Transactions {
		err := t.CheckBalance()
		check(err)

		for _, a := range t.Accounts {
			if b, ok := balance[a.Name]; ok {
				if a.Debit {
					b.Add(b, a.Amount)
				} else {
					b.Sub(b, a.Amount)
				}
			} else {
				if a.Debit {
					balance[a.Name] = a.Amount
				} else {
					neg := new(big.Rat)
					neg.SetInt64(-1)
					balance[a.Name] = a.Amount.Mul(a.Amount, neg)
				}
			}
		}
	}

	keys := make([]string, len(balance))
	i := 0
	for k := range balance {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	boom := ""
	for _, key := range keys {
		boom += fmt.Sprintf("%s\t\t%s\n", key, balance[key].FloatString(2))
	}
	return boom
}
