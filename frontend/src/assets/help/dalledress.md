<!--
Copyright 2016, 2026 The Authors. All rights reserved.
Use of this source code is governed by a license that can
be found in the LICENSE file.

Parts of this file were auto generated. Edit only those parts of
the code inside of 'EXISTING_CODE' tags.
-->
# DalleDress View

// EXISTING_CODE
// EXISTING_CODE

## Facets

- Generator Facet uses Generator store.
- Series Facet uses Series store.
- Databases Facet uses Databases store.
- Events Facet uses Logs store.
- Gallery Facet uses Logs store.

## Stores

- **Databases Store (0 members)**


- **Generator Store (0 members)**


- **Logs Store (14 members)**

  - blockNumber: the number of the block
  - transactionIndex: the zero-indexed position of the transaction in the block
  - logIndex: the zero-indexed position of this log relative to the block
  - timestamp: the timestamp of the block this log appears in
  - date: the timestamp as a date
  - address: the smart contract that emitted this log
  - name: the name of the smart contract if any
  - topics: the first topic hashes event signature of the log, up to 3 additional index parameters may appear
  - data: any remaining un-indexed parameters to the event
  - transactionHash: the hash of the transction
  - blockHash: the hash of the block
  - articulatedLog: a human-readable version of the topic and data fields
  - compressedLog: a truncated, more readable version of the articulation
  - isNFT: true if the log is an NFT transfer

- **Series Store (0 members)**


// EXISTING_CODE
// EXISTING_CODE
