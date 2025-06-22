package store

import "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"

type CollectionKey struct {
	Chain   string       // may be empty
	Address base.Address // may be empty
}

func GetCollectionKey(chain, address string) CollectionKey {
	return CollectionKey{
		Chain:   chain,
		Address: base.HexToAddress(address)}
}
