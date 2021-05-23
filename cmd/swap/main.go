package main

import (
	"bufio"
	"context"
	"eth-trasnsfer/lib"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

//func main()  {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	url := "https://bsc-dataseed1.binance.org:443"
//	secret := "058df786cdb66a72ad1408d75be42b37bb9b0f7b4e1863ac6945dceb244a5419"
//	pk, err := lib.HexToPk(secret)
//	if err != nil {
//		fmt.Printf("secret invalid: %v\n", err)
//		return
//	}
//
//	sender, err := lib.NewSender(ctx, url, pk, 1)
//	if err != nil {
//		fmt.Printf("sender error: %v\n", err)
//		return
//	}
//
//	pkfile, err := os.Open("./pk.txt")
//	if err != nil {
//		fmt.Printf("open file error: %v\n", err)
//		return
//	}
//	defer pkfile.Close()
//
//	scanner := bufio.NewScanner(pkfile)
//	for scanner.Scan() {
//		p, err := lib.HexToPk(scanner.Text())
//		if err != nil {
//			fmt.Printf("invalid pk: %v\n", err)
//			return
//		}
//
//		bal, err := sender.GetERC20Balance(ctx, contract, lib.PkToAddress(p))
//		if err != nil {
//			fmt.Printf("get erc20 balance failed: %v\n", err)
//			return
//		} else {
//			fmt.Println(bal.String())
//		}
//		to := lib.PkToAddress(p)
//		si, err := lib.NewSendInfo(to, "0.001", 18)
//		if err != nil {
//			fmt.Printf("invalid pk: %v\n", err)
//			return
//		}
//		err = sender.SendETH(ctx, si)
//		if err != nil {
//			fmt.Printf("send error: %v\n", err)
//			return
//		} else {
//			fmt.Printf("success send to %s\n", to)
//			<-time.After(time.Second)
//		}
//	}
//}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := "https://bsc-dataseed1.binance.org:443"
	swapRouter := common.HexToAddress("0x10ed43c718714eb63d5aa57b78b54704e256024e")
	base := common.HexToAddress("0x55d398326f99059ff775485246999027b3197955")
	pkfile, err := os.Open("./pk.txt")
	if err != nil {
		fmt.Printf("open file error: %v\n", err)
		return
	}
	defer pkfile.Close()

	scanner := bufio.NewScanner(pkfile)
	for scanner.Scan() {
		p, err := lib.HexToPk(scanner.Text())
		if err != nil {
			fmt.Printf("invalid pk: %v\n", err)
			return
		}

		sender, err := lib.NewSender(ctx, url, p, 1)
		if err != nil {
			fmt.Printf("sender %s error: %v\n", lib.PkToAddress(p), err)
			return
		}

		err = sender.ApproveERC20(ctx, base, swapRouter)
		if err != nil {
			fmt.Printf("%s approve erc20 balance failed: %v\n", sender.From(), err)
			return
		} else {
			fmt.Printf("%s approve erc20 success\n", sender.From())
		}
	}
}

//func main() {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	url := "https://bsc-dataseed1.binance.org:443"
//	swapRouter := common.HexToAddress("0x10ed43c718714eb63d5aa57b78b54704e256024e")
//	quote := common.HexToAddress("0x651bb1691f691e2ff2ef8dbc23dec8cf3c121366")
//	base := common.HexToAddress("0x55d398326f99059ff775485246999027b3197955")
//	pkfile, err := os.Open("./pk.txt")
//	if err != nil {
//		fmt.Printf("open file error: %v\n", err)
//		return
//	}
//	defer pkfile.Close()
//
//	scanner := bufio.NewScanner(pkfile)
//	for scanner.Scan() {
//		p, err := lib.HexToPk(scanner.Text())
//		if err != nil {
//			fmt.Printf("invalid pk: %v\n", err)
//			return
//		}
//
//		sender, err := lib.NewSender(ctx, url, p, 1)
//		if err != nil {
//			fmt.Printf("sender %s error: %v\n", lib.PkToAddress(p), err)
//			return
//		}
//		//get balance
//		bal, err := sender.GetERC20Balance(ctx, base, sender.From())
//		if err != nil {
//			fmt.Printf("%s get erc20 balance failed: %v\n", sender.From(), err)
//			return
//		}
//		//approve
//		//err = sender.ApproveERC20(ctx, base, swapRouter)
//		//if err != nil {
//		//	fmt.Printf("%s approve erc20 balance failed: %v\n", sender.From(), err)
//		//	return
//		//}
//		//swap
//		err = sender.SwapERC20(ctx, swapRouter, &lib.SwapInfo{
//			AmountIn:     bal,
//			AmountOutMin: big.NewInt(0),
//			Path:         []common.Address{base, quote},
//			To:           sender.From(),
//			Deadline:     big.NewInt(1622476800),
//		})
//		if err != nil {
//			fmt.Printf("%s swap erc20 balance failed: %v\n", sender.From(), err)
//			return
//		} else {
//			fmt.Printf("%s swap erc20 success\n", sender.From())
//			<-time.After(time.Duration(30 + rand.Int63n(170)) * time.Second)
//		}
//	}
//}
