## LFS - High Performance Key Value Storage Engine

### Storage

Records are written to an append-only binary file (AOF) and in-memory hash-tables are used for indexing records. Hash tables keep track of the offset of a record for a given key allowing for fast access.

Data is stored in the following format: 

`<keyLength>key<valueLength>value` 

The `keyLength` and `valueLength` are binary encoded `uint16` integers.

Both `key` and `value` are of type `string`.

### Usage

Connect to the server using any TCP client such as netcat (nc)

`nc localhost 8080`

Once connected, you can use any of the available commands to interact with the engine - 

1. `set [key] [value]`
2. `get [key]`
3. `exit`