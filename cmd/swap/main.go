package main

import (
	"bufio"
	"context"
	"eth-trasnsfer/lib"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := "https://bsc-dataseed1.binance.org:443"
	secret := "058df786cdb66a72ad1408d75be42b37bb9b0f7b4e1863ac6945dceb244a5419"
	pk, err := lib.HexToPk(secret)
	if err != nil {
		fmt.Printf("secret invalid: %v\n", err)
		return
	}

	sender, err := lib.NewSender(ctx, url, pk, 1)
	if err != nil {
		fmt.Printf("sender error: %v\n", err)
		return
	}

	contract := common.HexToAddress("0x55d398326f99059ff775485246999027b3197955")

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

		bal, err := sender.GetERC20Balance(ctx, contract, lib.PkToAddress(p))
		if err != nil {
			fmt.Printf("get erc20 balance failed: %v\n", err)
			return
		} else {
			fmt.Println(bal.String())
		}
	}

	//to := common.HexToAddress("0x58b28a163C326205Fa31aC7084E6E8aFCad181E0")
	//err = sender.SendETH(ctx, to, "0.15")
	//if err != nil {
	//	fmt.Println(3, err)
	//	return
	//} else {
	//	fmt.Println("success")
	//}
}
