package main

import (
	"bufio"
	"errors"
	"eth-trasnsfer/lib"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	app := &cli.App{
		Name:    "eth-sender",
		Usage:   "ETH automatic sending server",
		Version: "v0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "pk",
				Usage: "private key",
			},
			&cli.StringFlag{
				Name:  "url",
				Usage: "chain endpoint to interact with",
			},
			&cli.Float64Flag{
				Name:  "multiplier",
				Usage: "use (estimateGasPrice * multiplier) as actual gas price",
				Value: 1,
			},
		},
		Commands: []*cli.Command{
			sendETHCmd,
			sendERC20Cmd,
			batchSendETHCmd,
			batchSendERC20Cmd,
			cancelCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

var sendETHCmd = &cli.Command{
	Name:  "send-eth",
	Usage: "send eth",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "to",
			Usage: "send eth to address",
		},
		&cli.StringFlag{
			Name:  "amount",
			Usage: "send eth amount",
		},
	},
	Action: func(cctx *cli.Context) error {
		if !cctx.IsSet("pk") {
			return errors.New("--pk is not passed")
		}
		if !cctx.IsSet("url") {
			return errors.New("--url is not passed")
		}
		if !cctx.IsSet("to") {
			return errors.New("--to is not passed")
		}
		if !cctx.IsSet("amount") {
			return errors.New("--amount is not passed")
		}

		pk, err := lib.HexToPk(cctx.String("pk"))
		if err != nil {
			return err
		}

		sender, err := lib.NewSender(cctx.Context, cctx.String("url"), pk, cctx.Float64("multiplier"))
		if err != nil {
			return err
		}

		to := common.HexToAddress(cctx.String("to"))
		si, err := lib.NewSendInfo(to, cctx.String("amount"), 18)
		if err != nil {
			return err
		} else {
			return sender.SendETH(cctx.Context, si)
		}
	},
}

var sendERC20Cmd = &cli.Command{
	Name:  "send-erc20",
	Usage: "send erc20 token",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "contract",
			Usage: "erc20 contract address",
		},
		&cli.StringFlag{
			Name:  "to",
			Usage: "send eth to address",
		},
		&cli.StringFlag{
			Name:  "amount",
			Usage: "send eth amount",
		},
		&cli.Int64Flag{
			Name:  "prec",
			Usage: "token precision",
			Value: 18,
		},
	},
	Action: func(cctx *cli.Context) error {
		if !cctx.IsSet("pk") {
			return errors.New("--pk is not passed")
		}
		if !cctx.IsSet("url") {
			return errors.New("--url is not passed")
		}
		if !cctx.IsSet("to") {
			return errors.New("--to is not passed")
		}
		if !cctx.IsSet("amount") {
			return errors.New("--amount is not passed")
		}
		if !cctx.IsSet("contract") {
			return errors.New("--contract is not passed")
		}

		pk, err := lib.HexToPk(cctx.String("pk"))
		if err != nil {
			return err
		}

		sender, err := lib.NewSender(cctx.Context, cctx.String("url"), pk, cctx.Float64("multiplier"))
		if err != nil {
			return err
		}

		contract := common.HexToAddress(cctx.String("contract"))
		to := common.HexToAddress(cctx.String("to"))
		si, err := lib.NewSendInfo(to, cctx.String("amount"), cctx.Int64("prec"))
		if err != nil {
			return err
		} else {
			return sender.SendERC20(cctx.Context, contract, si)
		}
	},
}

var batchSendETHCmd = &cli.Command{
	Name:  "batch-send-eth",
	Usage: "batch send eth",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "to-path",
			Usage: "addresses file path",
		},
		&cli.Float64Flag{
			Name:  "min-amount",
			Usage: "min amount",
		},
		&cli.Float64Flag{
			Name:  "max-amount",
			Usage: "max amount",
		},
		&cli.Int64Flag{
			Name: "interval",
			Usage: "sending interval (seconds)",
			Value: 1,
		},
	},
	Action: func(cctx *cli.Context) error {
		if !cctx.IsSet("pk") {
			return errors.New("--pk is not passed")
		}
		if !cctx.IsSet("url") {
			return errors.New("--url is not passed")
		}
		if !cctx.IsSet("to-path") {
			return errors.New("--to-path is not passed")
		}
		if !cctx.IsSet("min-amount") {
			return errors.New("--min-amount is not passed")
		}
		if !cctx.IsSet("max-amount") {
			return errors.New("--max-amount is not passed")
		}

		pk, err := lib.HexToPk(cctx.String("pk"))
		if err != nil {
			return err
		}

		sender, err := lib.NewSender(cctx.Context, cctx.String("url"), pk, cctx.Float64("multiplier"))
		if err != nil {
			return err
		}

		toFile, err := os.Open(cctx.String("to-path"))
		if err != nil {
			return fmt.Errorf("open to file error: %v", err)
		}
		defer toFile.Close()

		a := cctx.Float64("min-amount")
		b := cctx.Float64("max-amount") - a
		interval := time.Duration(cctx.Int64("interval")) * time.Second

		scanner := bufio.NewScanner(toFile)
		for scanner.Scan() {
			<-time.After(interval)
			to := common.HexToAddress(scanner.Text())
			amount := strconv.FormatFloat(a + rand.Float64() * b, 'f', 10, 64)
			si, err := lib.NewSendInfo(to, amount, 18)
			if err != nil {
				return err
			}
			err = sender.SendETH(cctx.Context, si)
			if err != nil {
				return err
			} else {
				fmt.Printf("send eth %s to %s ok\n", amount, to)
			}
		}
		return nil
	},
}

var batchSendERC20Cmd = &cli.Command{
	Name:  "batch-send-erc20",
	Usage: "batch send erc20",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "contract",
			Usage: "erc20 contract address",
		},
		&cli.StringFlag{
			Name:  "to-path",
			Usage: "addresses file path",
		},
		&cli.Float64Flag{
			Name:  "min-amount",
			Usage: "min amount",
		},
		&cli.Float64Flag{
			Name:  "max-amount",
			Usage: "max amount",
		},
		&cli.Int64Flag{
			Name: "interval",
			Usage: "sending interval (seconds)",
			Value: 1,
		},
		&cli.Int64Flag{
			Name:  "prec",
			Usage: "token precision",
			Value: 18,
		},
	},
	Action: func(cctx *cli.Context) error {
		if !cctx.IsSet("contract") {
			return errors.New("--contract is not passed")
		}
		if !cctx.IsSet("pk") {
			return errors.New("--pk is not passed")
		}
		if !cctx.IsSet("url") {
			return errors.New("--url is not passed")
		}
		if !cctx.IsSet("to-path") {
			return errors.New("--to-path is not passed")
		}
		if !cctx.IsSet("min-amount") {
			return errors.New("--min-amount is not passed")
		}
		if !cctx.IsSet("max-amount") {
			return errors.New("--max-amount is not passed")
		}

		pk, err := lib.HexToPk(cctx.String("pk"))
		if err != nil {
			return err
		}

		sender, err := lib.NewSender(cctx.Context, cctx.String("url"), pk, cctx.Float64("multiplier"))
		if err != nil {
			return err
		}

		toFile, err := os.Open(cctx.String("to-path"))
		if err != nil {
			return fmt.Errorf("open to file error: %v", err)
		}
		defer toFile.Close()

		contract := common.HexToAddress(cctx.String("contract"))
		a := cctx.Float64("min-amount")
		b := cctx.Float64("max-amount") - a
		interval := time.Duration(cctx.Int64("interval")) * time.Second
		prec := cctx.Int64("prec")

		scanner := bufio.NewScanner(toFile)
		for scanner.Scan() {
			<-time.After(interval)
			to := common.HexToAddress(scanner.Text())
			amount := strconv.FormatFloat(a + rand.Float64() * b, 'f', 10, 64)
			si, err := lib.NewSendInfo(to, amount, prec)
			if err != nil {
				return err
			}
			err = sender.SendERC20(cctx.Context, contract, si)
			if err != nil {
				return err
			} else {
				fmt.Printf("send erc20 token %s to %s ok\n", amount, to)
			}
		}
		return nil
	},
}

var cancelCmd = &cli.Command{
	Name:  "cancel",
	Usage: "cancel pending transactions",
	Action: func(cctx *cli.Context) error {
		if !cctx.IsSet("pk") {
			return errors.New("--pk is not passed")
		}
		if !cctx.IsSet("url") {
			return errors.New("--url is not passed")
		}

		pk, err := lib.HexToPk(cctx.String("pk"))
		if err != nil {
			return err
		}

		sender, err := lib.NewSender(cctx.Context, cctx.String("url"), pk, cctx.Float64("multiplier"))
		if err != nil {
			return err
		} else {
			return sender.CancelPending(cctx.Context)
		}
	},

}