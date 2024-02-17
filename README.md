## LFS - High Performance Key Value Storage Engine

Records are written to an append-only binary file (AOF) and in-memory hash-tables are used for indexing records.

Records are stored in the following format: 

`<keyLength>key<valueLength>value` 

The `keyLength` and `valueLength` are binary encoded `uint16` bit integers. Both `key` and `value` should be of type `string`. Hash maps are used to store the offset of a given key for efficient reading.