# Purpose

A file server for edge clients providing file contents stored
in a local directory. It uses gRPC connection to fetch missing files
from sync server which, in its turn, downloads the missing files
from file service. The files are transferred in async mode
so that the request method does not wait for the called party.
The response contains a promise treshold value for recommended
time before next request can be done.

## REST API

It provides the following methods:
- GET:/health - system health
- GET:/metrics - system Prometheus metrics
- HEAD:/api/v1/files/${id} - file raw stat metadata like size, cksum SHA256
- GET:/api/v1/files/${id} - file contents with possible Range tags

### GRPC FileProxyIPC

It provides methods to access the GRPC sync service:
- GetFile(GetFileRequest):GetFileResponse - getting the files by id by async download
- Ping(PingRequest):PingResponse - sanity check

## Local File Storage

The files are stored id a dir with file names equal to id.

## Config - config.yml

It provides configuration options for:
- logging
- rest service options
- grpc client for server connection
- file storage dir path
