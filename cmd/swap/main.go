package main

import (
	"bufio"
	"context"
	"eth-trasnsfer/lib"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"math/rand"
	"os"
	"time"
)

//func main() {
//	f, err := os.Create("./pk2.txt")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer f.Close()
//
//	for i := 0; i < 500; i++ {
//		pk, err := lib.CreateNewPk()
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			f.WriteString(fmt.Sprintf("%s,%s\n", lib.PkToAddress(pk), lib.PkToHex(pk)))
//		}
//	}
//}

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
//	//usdt := common.HexToAddress("0x55d398326f99059ff775485246999027b3197955")
//	//total := big.NewInt(40000)
//	//total.Mul(total, big.NewInt(1000_000_000_000_000_000))
//
//	i := 0
//	scanner := bufio.NewScanner(pkfile)
//	for scanner.Scan() {
//		i++
//		p, err := lib.HexToPk(scanner.Text())
//		if err != nil {
//			fmt.Printf("invalid pk: %v\n", err)
//			return
//		}
//
//		to := lib.PkToAddress(p)
//		si, err := lib.NewSendInfo(to, "0.002", 18)
//		if err != nil {
//			fmt.Printf("invalid pk: %v\n", err)
//			return
//		}
//		err = sender.SendETH(ctx, si)
//		if err != nil {
//			fmt.Printf("send error: %v\n", err)
//			return
//		} else {
//			fmt.Printf("%d: send 0.002 eth to %s\n", i, si.To)
//		}
//
//		//to := lib.PkToAddress(p)
//		//amount := 10 + rand.Float64() * 90
//		//si, err := lib.NewSendInfo(to, strconv.FormatFloat(amount, 'f', 10, 64), 18)
//		//if err != nil {
//		//	fmt.Printf("invalid pk: %v\n", err)
//		//	return
//		//}
//		//
//		//if total.Cmp(si.Amount) < 0 {
//		//	fmt.Println("finished")
//		//	return
//		//}
//		//total.Sub(total, si.Amount)
//		//err = sender.SendERC20(ctx, usdt, si)
//		//if err != nil {
//		//	fmt.Printf("invalid pk: %v\n", err)
//		//	return
//		//} else {
//		//	fmt.Printf("%d: send usdt %s to %s successfully\n", i, si.Amount, si.To)
//		//}
//	}
//}

//func main() {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	url := "https://bsc-dataseed1.binance.org:443"
//	swapRouter := common.HexToAddress("0x10ed43c718714eb63d5aa57b78b54704e256024e")
//	base := common.HexToAddress("0x55d398326f99059ff775485246999027b3197955")
//	pkfile, err := os.Open("./pk.txt")
//	if err != nil {
//		fmt.Printf("open file error: %v\n", err)
//		return
//	}
//	defer pkfile.Close()
//
//	i := 0
//	scanner := bufio.NewScanner(pkfile)
//	for scanner.Scan() {
//		i++
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
//
//		err = sender.ApproveERC20(ctx, base, swapRouter)
//		if err != nil {
//			fmt.Printf("%d: %s approve erc20 balance failed: %v\n", i, sender.From(), err)
//			time.After(time.Second * 5)
//		} else {
//			fmt.Printf("%d: %s approve erc20 success\n", i, sender.From())
//		}
//	}
//}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := "https://bsc-dataseed1.binance.org:443"
	swapRouter := common.HexToAddress("0x10ed43c718714eb63d5aa57b78b54704e256024e")
	quote := common.HexToAddress("0x651bb1691f691e2ff2ef8dbc23dec8cf3c121366")
	base := common.HexToAddress("0x55d398326f99059ff775485246999027b3197955")
	pkfile, err := os.Open("./pk.txt")
	if err != nil {
		fmt.Printf("open file error: %v\n", err)
		return
	}
	defer pkfile.Close()

	i := 0
	scanner := bufio.NewScanner(pkfile)
	for scanner.Scan() {
		<-time.After(time.Duration(60 + rand.Int63n(120)) * time.Second)
		i++
		p, err := lib.HexToPk(scanner.Text())
		if err != nil {
			fmt.Printf("invalid pk: %v\n", err)
			return
		}

		sender, err := lib.NewSender(ctx, url, p, 1)
		if err != nil {
			fmt.Printf("%d: sender %s error: %v\n", i, lib.PkToAddress(p), err)
			continue
		}
		//get balance
		bal, err := sender.GetERC20Balance(ctx, base, sender.From())
		if err != nil {
			fmt.Printf("%d: %s get erc20 balance failed: %v\n", i, sender.From(), err)
			continue
		}
		//swap
		err = sender.SwapERC20(ctx, swapRouter, &lib.SwapInfo{
			AmountIn:     bal,
			AmountOutMin: big.NewInt(0),
			Path:         []common.Address{base, quote},
			To:           sender.From(),
			Deadline:     big.NewInt(1622476800),
		})
		if err != nil {
			fmt.Printf("%d: %s swap erc20 failed: %v\n", i, sender.From(), err)
		} else {
			fmt.Printf("%d: %s swap erc20 success\n", i, sender.From())
		}
	}
}
