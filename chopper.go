package main

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrInvalidAddress = fmt.Errorf("not a valid address")
var ErrInvalidSeed = fmt.Errorf("invalid seed")

func (a *App) Chopper(input string) (string, map[string]string, error) {
	if !base.IsValidAddress(input) {
		return "", map[string]string{}, ErrInvalidAddress
	}

	if strings.HasSuffix(input, ".eth") {
		var ok bool
		if input, ok = a.conn.GetEnsAddress(input); !ok {
			return "", map[string]string{}, ErrInvalidAddress
		}
	}

	hash := hexutil.Encode(crypto.Keccak256([]byte(input)))
	seed := hash[2:] + input[2:]
	if len(seed) < 104 {
		return "", map[string]string{}, ErrInvalidSeed
	}

	keys := []string{"adverb", "adjective", "emotionshort", "emotion", "literary", "noun", "style", "style2", "color1", "color2", "color3", "variant1", "variant2", "variant3", "background", "orientation"}
	starts := []int{0, 12, 24, 36, 48, 60, 72, 84, 92, 80, 68, 56, 44, 32, 20, 8}
	ends := []int{12, 24, 36, 48, 60, 72, 84, 96, 104, 92, 80, 68, 56, 44, 32, 20}

	segments := make(map[string]string, len(keys))
	for i, key := range keys {
		segments[key] = seed[starts[i]:ends[i]]
	}
	return seed, segments, nil
}
