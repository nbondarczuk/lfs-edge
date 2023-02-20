# Files handling controller

## Purpose (general concept)

It handles pending file gRPC requests from proxy. It starts download in async
mode. It sends a pending file request to channel where download agent
waits in a loop.

### Download agent (detailed description)

The agent waits on a channel for pending file requests. It handles them
by connecting to FS and downloading them using presigned url. The status
and metadata of the file are stored in the local database.

The files being downloaded are named as {file}.pending. They are stored
in the storage path directory. After successfull download they are renamed
to the final name and their status in the database is enriched with
file metadata like: name, size, md5 checksum. In case of error in transfer
onlys status is set to error.

### Sync agent (detailed description)

The new/deleted files are handles by sync agent which runs as a goroutine.
It is supposed to connect to source in a timely manner in order to get
the list of new/ deleted files since last  access. I uses download agent
to download them. the time stamp of the last access is stored in
a marker file. This file contains time stamp in unix format.