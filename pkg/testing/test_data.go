package testing

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

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
