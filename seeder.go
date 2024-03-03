package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrCouldNotResolveEns = fmt.Errorf("could not resolve ENS address")
var ErrNotAHexString = fmt.Errorf("input is not a valid hex string")

func (a *App) SeedBuilder(input string) (string, string, error) {
	if strings.HasSuffix(input, ".eth") {
		var ok bool
		if input, ok = a.conn.GetEnsAddress(input); !ok {
			return input, "", ErrCouldNotResolveEns
		}
	}

	isHexValid := func(s string) bool {
		trimmedString := strings.TrimPrefix(s, "0x")
		_, err := strconv.ParseInt(trimmedString, 16, 0)
		return err == nil
	}
	if !isHexValid(input) {
		return "", "", ErrNotAHexString
	}

	hash := hexutil.Encode(crypto.Keccak256([]byte(input)))
	if len(hash) < 2 || len(input) < 2 {
		return "", "", ErrInvalidSeed
	}

	seed := hash[2:] + input[2:]
	if len(seed) < 104 {
		return "", "", ErrInvalidSeed
	}

	return "0x" + input, "0x" + hash + input, nil
}
