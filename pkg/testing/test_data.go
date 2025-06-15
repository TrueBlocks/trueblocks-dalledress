package testing

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

func CreateTestAbis() []coreTypes.Abi {
	return []coreTypes.Abi{
		{
			Address:  base.HexToAddress("0x1234567890123456789012345678901234567890"),
			Name:     "Test ABI 1",
			IsKnown:  false,
			FileSize: 1024,
		},
		{
			Address:  base.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			Name:     "Test ABI 2",
			IsKnown:  true,
			FileSize: 2048,
		},
	}
}

func CreateTestFunctions() []coreTypes.Function {
	return []coreTypes.Function{
		{
			Name:         "transfer",
			FunctionType: "function",
			Signature:    "transfer(address,uint256)",
			Encoding:     "0xa9059cbb",
		},
		{
			Name:         "Transfer",
			FunctionType: "event",
			Signature:    "Transfer(address,address,uint256)",
			Encoding:     "0xddf252ad1be2c18e81c25cf750fee4f7ad3e5c8f70b1d4a6b4f8e12e4a6f8c23",
		},
		{
			Name:         "approve",
			FunctionType: "function",
			Signature:    "approve(address,uint256)",
			Encoding:     "0x095ea7b3",
		},
		{
			Name:         "Approval",
			FunctionType: "event",
			Signature:    "Approval(address,address,uint256)",
			Encoding:     "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
		},
	}
}

func CreateTestMonitors() []coreTypes.Monitor {
	return []coreTypes.Monitor{
		{
			Address:     base.HexToAddress("0x1234567890123456789012345678901234567890"),
			Name:        "Test Monitor 1",
			NRecords:    100,
			FileSize:    1024,
			LastScanned: 12345,
			IsEmpty:     false,
			IsStaged:    true,
			Deleted:     false,
		},
		{
			Address:     base.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			Name:        "Test Monitor 2",
			NRecords:    250,
			FileSize:    2048,
			LastScanned: 67890,
			IsEmpty:     true,
			IsStaged:    false,
			Deleted:     false,
		},
	}
}

func CreateTestNames() []coreTypes.Name {
	return []coreTypes.Name{
		{
			Address:    base.HexToAddress("0x1234567890123456789012345678901234567890"),
			Name:       "Test Name 1",
			Tags:       "testing",
			Source:     "test",
			Symbol:     "TEST1",
			Decimals:   18,
			IsCustom:   true,
			IsPrefund:  false,
			IsContract: true,
			IsErc20:    true,
			IsErc721:   false,
			Parts:      coreTypes.Custom,
		},
		{
			Address:    base.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			Name:       "Test Name 2",
			Tags:       "prefund",
			Source:     "genesis",
			Symbol:     "TEST2",
			Decimals:   0,
			IsCustom:   false,
			IsPrefund:  true,
			IsContract: false,
			IsErc20:    false,
			IsErc721:   false,
			Parts:      coreTypes.Prefund,
		},
		{
			Address:    base.HexToAddress("0x9876543210987654321098765432109876543210"),
			Name:       "Test Regular Name",
			Tags:       "regular",
			Source:     "etherscan",
			Symbol:     "TRN",
			Decimals:   8,
			IsCustom:   false,
			IsPrefund:  false,
			IsContract: true,
			IsErc20:    true,
			IsErc721:   false,
			Parts:      coreTypes.Regular,
		},
		{
			Address:    base.HexToAddress("0xfedcba0987654321fedcba0987654321fedcba09"),
			Name:       "Bad Address",
			Tags:       "baddress",
			Source:     "manual",
			Symbol:     "BAD",
			Decimals:   0,
			IsCustom:   false,
			IsPrefund:  false,
			IsContract: false,
			IsErc20:    false,
			IsErc721:   false,
			Parts:      coreTypes.Baddress,
		},
	}
}
