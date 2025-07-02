# Abis View

// EXISTING_CODE

Welcome to the **Abis** view! This section provides information about managing abis in your application.

// EXISTING_CODE

## Facets:
- Downloaded Facet uses Abis store.
- Known Facet uses Abis store.
- Functions Facet uses Functions store.
- Events Facet uses Functions store.

## Stores:

- **Abis Store (12 members):**

  - address: the address for the ABI
  - name: the filename of the ABI (likely the smart contract address)
  - fileSize: the size of this file on disc
  - nFunctions: the number of functions in the ABI
  - nEvents: the number of events in the ABI
  - path: the folder holding the abi file
  - lastModDate: the last update date of the file
  - isKnown: true if this is the ABI for a known smart contract or protocol
  - isEmpty: true if the ABI could not be found (and won't be looked for again)
  - hasConstructor: if verbose and the abi has a constructor, then `true`, else `false`
  - hasFallback: if verbose and the abi has a fallback, then `true`, else `false`
  - functions: the functions for this address

- **Functions Store (10 members):**

  - encoding: the signature encoded with keccak
  - name: the name of the interface
  - type: the type of the interface, either 'event' or 'function'
  - signature: the canonical signature of the interface
  - anonymous: 
  - constant: 
  - stateMutability: 
  - message: 
  - inputs: the input parameters to the function, if any
  - outputs: the output parameters to the function, if any

// EXISTING_CODE
// EXISTING_CODE
