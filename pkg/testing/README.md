# Testing Package

This package provides shared testing utilities and mocks for the TrueBlocks Dalledress project.

## Overview

The testing package consolidates common testing patterns, mock implementations, and utilities that were previously duplicated across multiple test files. This includes:

- **Message Capture functionality** - Thread-safe message capturing for testing event emissions
- **Mock SDK implementations** - Base mock structures for SDK interactions  
- **Test data factories** - Standardized test data creation for different types
- **Testing utilities** - Common helper functions for timeouts and assertions

## Components

### MessageCapture

Provides thread-safe capturing of messages emitted during tests:

```go
capture := testing.NewMessageCapture()
// Use capture.CaptureLoaded() to capture messages
messages := capture.GetMessages()
```

### Mock SDK Base

Generic base class for creating type-specific mock SDKs:

```go
type MockSDKAbi struct {
    *testing.MockSDKBase[coreTypes.Abi]
    // Additional abi-specific fields
}

func NewMockSDKAbi() *MockSDKAbi {
    return &MockSDKAbi{
        MockSDKBase: testing.NewMockSDKBase(testing.CreateTestAbis()),
    }
}
```

### Test Data Factories

Standardized factory functions for creating consistent test data:

- `CreateTestAbis()` - Returns standard ABI test data
- `CreateTestMonitors()` - Returns standard Monitor test data  
- `CreateTestNames()` - Returns standard Name test data
- `CreateTestFunctions()` - Returns standard Function test data

### Testing Utilities

Common testing helper functions:

- `WaitWithTimeout()` - Wait for a condition with timeout
- `AssertEventually()` - Assert that a condition becomes true within timeout
- `CreateSafeTimeout()` - Create conservative timeouts for tests
- `GetTestTimeout()` - Get standard test timeout duration

## Usage in Test Files

The consolidated testing package is used in the setup_test.go files across different packages:

```go
import testingPkg "github.com/TrueBlocks/trueblocks-dalledress/pkg/testing"

// Use consolidated types
type MessageCapture = testingPkg.MessageCapture
type CapturedMessage = testingPkg.CapturedMessage

// Create mocks using base functionality
type MockSDKAbi struct {
    *testingPkg.MockSDKBase[coreTypes.Abi]
    functions []coreTypes.Function
}

func NewMockSDKAbi() *MockSDKAbi {
    return &MockSDKAbi{
        MockSDKBase: testingPkg.NewMockSDKBase(testingPkg.CreateTestAbis()),
        functions:   testingPkg.CreateTestFunctions(),
    }
}
```

## Benefits

1. **Reduced Duplication** - Eliminates identical code across multiple setup_test.go files
2. **Consistency** - Ensures all tests use the same mock data and patterns
3. **Maintainability** - Changes to testing patterns only need to be made in one place
4. **Type Safety** - Uses Go generics for type-safe mock implementations
5. **Extensibility** - Easy to add new test utilities and mock types

## Migration

The existing setup_test.go files have been updated to use this consolidated package while maintaining the same public interface, ensuring all existing tests continue to work without modification.
