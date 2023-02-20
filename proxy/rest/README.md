# Proxy REST API

## Methods

The rest API provides the following methods:

- GET:/health - system health

- GET:/metrics - system Prometheus metrics

- HEAD:/api/v1/files/{id} - file raw stat metadata like size, cksum MD5, last modified ts (no created ts)

- GET:/api/v1/files/{id}?device={device} - render file {id} contents for {device}
