# Monitors View

// EXISTING_CODE
// EXISTING_CODE

## Facets:
- Monitors Facet uses Monitors store.

## Stores:

- **Monitors Store (8 members):**

  - address: the address of this monitor
  - name: the name of this monitor (if any)
  - nRecords: the number of appearances for this monitor
  - fileSize: the size of this monitor on disc
  - isEmpty: true if the monitor has no appearances
  - lastScanned: the last scanned block number
  - deleted: if this monitor has been deleted, `false` otherwise
  - isStaged: if the monitor file in on the stage, `false` otherwise

// EXISTING_CODE
// EXISTING_CODE
