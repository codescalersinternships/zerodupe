# Zerodupe

zerodupe is a content-addressable storage server designed for space efficiency through block-level deduplication. Clients interact with the server via a REST API to upload and download data. Files are split into fixed-size blocks, each identified by the hash of its content. The server stores each unique block only once. Metadata linking files to their constituent blocks is maintained separately.
---

## Usage

### Flags

| Flag        | Description                                                 |
| ----------- | ----------------------------------------------------------- | 
| `-upload`   | Path to the file you want to upload                         | 
| `-download` | Hash of the file you want to download                       | 
| `-out`      | Output path for the downloaded file (default: `downloaded`) | 


### Examples
```bash
cd cmd/client
go run main.go -upload ../../testdata/files/test1.txt
go run main.go -download 9d259d5d9057fec99f14f4025f76188bff1de029cffe020440a474e8739a7719 -out file.txt
```
