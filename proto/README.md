# GRPC remote file fetcher interface

## Purpose

The interface checks for a file by id. If the file is missing in
local store, it initiates async transfer of the missing file.
When a file is not immediately available, status code `FAILED_PRECONDITION` is returned
to the clients indicating that the file is not ready.

For status codes such as `FAILED_PRECONDITION`, a `RetryAfterSeconds` is applicable.
Clients are adviced to retry after the specified seconds. This is a way
for the sync service to manage client expectations on large file downloads.

If a file is not available at source,`NOT_FOUND` is returned. Clients should
not attempt this file again unless they resolve the original error.
`RetryAfterSeconds` will be set to invalid value `-1` in cases where retrying
is not likely to change the outcome.

### How to build
use `make` command. The protoc version used is `libprotoc 3.15.8`
