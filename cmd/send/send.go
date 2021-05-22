package main

import (
	"errors"
	"eth-trasnsfer/lib"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
	"os"
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
			return sender.SendERC20Token(cctx.Context, contract, si)
		}
	},
}