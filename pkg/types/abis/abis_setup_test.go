package abis

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	testingPkg "github.com/TrueBlocks/trueblocks-dalledress/pkg/testing"
)

func CreateTestAbis() []Abi {
	return []Abi{
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

func CreateTestFunctions() []Function {
	return []Function{
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

type MockSDKAbi struct {
	abis      []Abi
	functions []Function
}

func NewMockSDKAbi() *MockSDKAbi {
	return &MockSDKAbi{
		abis:      CreateTestAbis(),
		functions: CreateTestFunctions(),
	}
}

func (m *MockSDKAbi) SetAbis(abis []Abi) {
	m.abis = abis
}

func (m *MockSDKAbi) SetFunctions(functions []Function) {
	m.functions = functions
}

type MessageCapture = testingPkg.MessageCapture
type CapturedMessage = testingPkg.CapturedMessage

func NewMessageCapture() *MessageCapture {
	return testingPkg.NewMessageCapture()
}
