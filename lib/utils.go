package lib

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
)

type SendInfo struct {
	To 			common.Address
	Amount 		*big.Int
}

func NewSendInfo(to common.Address, amount string, prec int64) (*SendInfo, error) {
	value, ok := new(big.Rat).SetString(amount)
	if !ok {
		return nil, fmt.Errorf("invalid amount %s", amount)
	}
	value = value.Mul(value, big.NewRat(math.BigPow(10, prec).Int64(), 1))
	return &SendInfo{
		To:     to,
		Amount: value.Num(),
	}, nil
}

type SwapInfo struct {
	AmountIn 		*big.Int
	AmountOutMin 	*big.Int
	Path 			[]common.Address
	To 				common.Address
	Deadline 		*big.Int
}