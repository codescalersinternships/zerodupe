## zerodupe is a Deduplicated Storage Server.

zerodupe is a content-addressable storage server designed for space efficiency through block-level deduplication. Clients interact with the server via a REST API to upload and download data. Files are split into fixed-size blocks, each identified by the hash of its content. The server stores each unique block only once. Metadata linking files to their constituent blocks is maintained separately.

### Notes
- Provide a REST API for block and file operations.
- Store data blocks based on their content hash (SHA-256).
- Avoid storing duplicate block data.
- Maintain metadata associating an ordered list of block hashes with a unique file identifier (file hash).
- Provide a CLI client for uploading and downloading files.
- Client should perform client-side checks to avoid uploading existing blocks.
- Client should support parallel block uploads (after sequential implementation).
- Utilize the local filesystem for both block and metadata storage.
- Ensure correctness through unit and end-to-end tests.


## Next step
- Directory uploads/downloads (upload-dir, download-dir).


## Definitions 

- Block: A fixed-size segment of a file's data (e.g., 4 MiB).
- Block Hash: The SHA-256 hash of a block's content, represented as a lowercase hexadecimal string. - This is the unique identifier for the block's data.
- File Hash: The SHA-256 hash of the concatenated sequence of its constituent block hashes 
- Content-Addressable Storage: Storing and retrieving data based on its content hash, rather than a location or user-defined name.


## Data storage
Data is stored directly on the filesystem within a configurable base directory (<storage_root>).
### Block Storage:
Stores the actual data chunks.
- Path: `<storage_root>/blocks/<first_4_chars_of_hash>/<full_block_hash>`
- Example: Block hash `f8ab123c4d`..., path `/path/to/storage/blocks/f8ab/f8ab123c4d`...
- Content: Raw binary data of the block.

### Metadata Storage:

#### in the file system


```
<storage_root>/
├── blocks/
│   ├── f8ab/
│   │   │   └── f8ab123c4d...         # Block data
└── meta/
    ├── a1b2/
    │       └── a1b2c3d4e5....    # File metadata (list of block hashes)

```

Stores the mapping between a file hash and its ordered list of block hashes.
- Path: `<storage_root>/meta/<first_4_chars_of_file_hash>/<full_file_hash>`
- Example: File hash `a1b2c3d4e5...,` path `/path/to/storage/meta/a1b2/a1b2c3d4e5..`
- Content: Plain text file, with each line containing one block hash in the correct order.
```
f8ab123c4d...
e7cd345a6b...
99ef567b8a...
```
OR
Alternative Content: JSON array `{"block_hashes": ["hash1", "hash2", ...]}`


#### in a sqlite database

We can have a table in such format

```
id     block_order block_hash  file_hash

```
and there we can retrieve all of the blocks relevant to a specific file_hash to download them and order them by block_order so we could rebuild the file after downloading the blocks	

ETA: 10 days