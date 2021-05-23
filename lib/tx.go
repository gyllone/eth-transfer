package lib

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/atomic"
	"math/big"
)

var (
	SendSig, _ = hex.DecodeString("a9059cbb")
	ApproveSig, _ = hex.DecodeString("095ea7b3")
	BalanceSig, _ = hex.DecodeString("70a08231")
	SwapSig, _ = hex.DecodeString("38ed1739")
)

type Sender struct {
	multiplier 	*big.Rat
	from 		common.Address
	nonce 		*atomic.Uint64
	sign 		func(*types.Transaction) (*types.Transaction, error)
	client		*rpc.Client
}

func NewSender(ctx context.Context, url string, pk *ecdsa.PrivateKey, multiplier float64) (*Sender, error) {
	client, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}

	chainID, err := ethclient.NewClient(client).ChainID(ctx)
	if err != nil {
		return nil, err
	}
	signer := types.NewEIP155Signer(chainID)

	from := PkToAddress(pk)
	nonce, err := ethclient.NewClient(client).PendingNonceAt(ctx, from)
	if err != nil {
		return nil, err
	}

	return &Sender{
		multiplier: new(big.Rat).SetFloat64(multiplier),
		from: from,
		nonce: atomic.NewUint64(nonce),
		sign: func(tx *types.Transaction) (*types.Transaction, error) {
			return types.SignTx(tx, signer, pk)
		},
		client: client,
	}, nil
}

func (s *Sender) From() common.Address {
	return s.from
}

func (s *Sender) GetPendingETHBalance(ctx context.Context, addr common.Address) (*big.Int, error) {
	return ethclient.NewClient(s.client).PendingBalanceAt(ctx, addr)
}

func (s *Sender) GetETHBalance(ctx context.Context, addr common.Address) (*big.Int, error) {
	return ethclient.NewClient(s.client).BalanceAt(ctx, addr, nil)
}

func (s *Sender) SendETH(ctx context.Context, si *SendInfo) error {
	return s.sendMessage(ctx, si.To, si.Amount, nil)
}

func (s *Sender) SendERC20(ctx context.Context, contract common.Address, si *SendInfo) error {
	// ERC20 data encode
	buffer1 := make([]byte, 32)
	copy(buffer1[12:], si.To.Bytes())
	buffer2 := make([]byte, 32)
	si.Amount.FillBytes(buffer2)
	data := make([]byte, 4, 4+32*2)
	copy(data, SendSig)
	data = append(data, buffer1...)
	data = append(data, buffer2...)

	return s.sendMessage(ctx, contract, nil, data)
}

func (s *Sender) ApproveERC20(ctx context.Context, contract, spender common.Address) error {
	buffer1 := make([]byte, 32)
	copy(buffer1[12:], spender.Bytes())
	buffer2 := bytes.Repeat([]byte{255}, 32)
	data := make([]byte, 4, 4+32*2)
	copy(data, ApproveSig)
	data = append(data, buffer1...)
	data = append(data, buffer2...)

	return s.sendMessage(ctx, contract, nil, data)
}

func (s *Sender) SwapERC20(ctx context.Context, contract common.Address, si *SwapInfo) error {
	buffer1 := make([]byte, 32)
	si.AmountIn.FillBytes(buffer1)
	buffer2 := make([]byte, 32)
	si.AmountOutMin.FillBytes(buffer2)
	buffer3 := make([]byte, 32)
	buffer3[31] = 0xa0
	buffer4 := make([]byte, 32)
	copy(buffer4[12:], si.To.Bytes())
	buffer5 := make([]byte, 32)
	si.Deadline.FillBytes(buffer5)
	buffer6 := make([]byte, 32)
	buffer6[31] = byte(len(si.Path))
	buffer7 := make([]byte, 32*len(si.Path))
	for i, p := range si.Path {
		copy(buffer7[12+32*i:], p.Bytes())
	}
	data := make([]byte, 4, 4+32*(6+len(si.Path)))
	copy(data, SwapSig)
	data = append(data, buffer1...)
	data = append(data, buffer2...)
	data = append(data, buffer3...)
	data = append(data, buffer4...)
	data = append(data, buffer5...)
	data = append(data, buffer6...)
	data = append(data, buffer7...)

	return s.sendMessage(ctx, contract, nil, data)
}

func (s *Sender) GetERC20Balance(ctx context.Context, contract, addr common.Address) (*big.Int, error) {
	// ERC20 data encode
	buffer := make([]byte, 32)
	copy(buffer[12:], addr.Bytes())
	data := make([]byte, 4, 4+32)
	copy(data, BalanceSig)
	data = append(data, buffer...)

	res, err := ethclient.NewClient(s.client).PendingCallContract(ctx, ethereum.CallMsg{
		From:   s.from,
		To:     &contract,
		Data: 	data,
	})
	if err != nil {
		return nil, err
	} else {
		return new(big.Int).SetBytes(res), nil
	}
}

func (s *Sender) BatchSendETH(ctx context.Context, sis []*SendInfo) error {
	gasPrice, err := ethclient.NewClient(s.client).SuggestGasPrice(ctx)
	if err != nil {
		return err
	}
	gp := new(big.Rat).SetInt(gasPrice)
	gp.Mul(gp, s.multiplier)

	nonce := s.nonce.Load()
	payload := make([]rpc.BatchElem, len(sis))
	for i, si := range sis {
		tx, err := s.sign(types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &si.To,
			Value:    si.Amount,
			Gas:      21_000,
			GasPrice: gp.Num(),
		}))
		if err != nil {
			return err
		}

		data, err := tx.MarshalBinary()
		if err != nil {
			return err
		} else {
			payload[i] = rpc.BatchElem{
				Method: "eth_sendRawTransaction",
				Args:   []interface{}{hexutil.Encode(data)},
			}
			nonce++
		}
	}

	err = s.client.BatchCallContext(ctx, payload)
	if err != nil {
		return err
	} else {
		s.nonce.Store(nonce)
		return nil
	}
}

func (s *Sender) sendMessage(ctx context.Context, to common.Address, value *big.Int, data []byte) error {
	gasLimit, err := ethclient.NewClient(s.client).EstimateGas(ctx, ethereum.CallMsg{
		From:   s.from,
		To:     &to,
		Value:  value,
		Data: 	data,
	})
	if err != nil {
		return err
	}

	gasPrice, err := ethclient.NewClient(s.client).SuggestGasPrice(ctx)
	if err != nil {
		return err
	}
	gp := new(big.Rat).SetInt(gasPrice)
	gp.Mul(gp, s.multiplier)

	tx, err := s.sign(types.NewTx(&types.LegacyTx{
		Nonce:    s.nonce.Load(),
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gp.Num(),
		Data:  	  data,
	}))
	if err != nil {
		return err
	}

	err = ethclient.NewClient(s.client).SendTransaction(ctx, tx)
	if err != nil {
		return err
	} else {
		s.nonce.Add(1)
		return nil
	}
}

func (s *Sender) pendingCallContract(ctx context.Context, contract common.Address, data []byte) ([]byte, error) {
	return ethclient.NewClient(s.client).PendingCallContract(ctx, ethereum.CallMsg{
		From:   s.from,
		To:     &contract,
		Data: 	data,
	})
}

func (s *Sender) callContract(ctx context.Context, contract common.Address, data []byte) ([]byte, error) {
	return ethclient.NewClient(s.client).CallContract(ctx, ethereum.CallMsg{
		From:   s.from,
		To:     &contract,
		Data: 	data,
	}, nil)
}