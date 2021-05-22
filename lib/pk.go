package lib

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func CreateNewPk() (*ecdsa.PrivateKey, error) {
	pk, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	} else {
		return pk, nil
	}
}

func HexToPk(secret string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(secret)
}

func PkToHex(key *ecdsa.PrivateKey) string {
	return hex.EncodeToString(crypto.FromECDSA(key))
}

func PkToAddress(pk *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(pk.PublicKey)
}
