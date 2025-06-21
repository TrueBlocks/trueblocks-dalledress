package names

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

func CreateTestNames() []Name {
	return []Name{
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

type MockSDKName struct {
	names []Name
}

func NewMockSDKName() *MockSDKName {
	return &MockSDKName{
		names: CreateTestNames(),
	}
}

func (m *MockSDKName) SetNames(names []Name) {
	m.names = names
}
