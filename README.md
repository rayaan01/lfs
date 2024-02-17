## LFS - High performance Key Value Storage Engine

### Records are written to an append-only binary file (AOF) and in-memory hash-tables are used for indexing records. Records are stored in this format: `<keyLength>key<valueLength>value` where the key length and value length are fixed uint16 bit integers. Hash maps are used to store the offset of a given key for efficient reading.