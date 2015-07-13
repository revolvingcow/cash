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

	padLength := 20
	for _, key := range keys {
		if len(key) > padLength {
			padLength = len(key) + 2
		}
	}

	boom := ""
	for _, key := range keys {
		extra := ""
		if balance[key].Sign() == 1 {
			extra = "+"
		}
		boom += fmt.Sprintf("%s%s%s\n", padRight(key, " ", padLength), extra, balance[key].FloatString(2))
	}
	return boom
}

func padRight(source, c string, length int) string {
	for i := len(source); i < length; i++ {
		source += c
	}
	return source
}

func padLeft(source, c string, length int) string {
	for i := len(source); i < length; i++ {
		source = c + source
	}
	return source
}
