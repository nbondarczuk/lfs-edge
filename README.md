## lfs-edge

`lfs-edge` allows a tenant specific local mirror for a remote file server.

`lfs-edge` will transparently mirror remote changes while also allowing
for on-demand fetches in case a file is requested outside the sync window.

`lfs-edge` is designed to serve a specialized group of clients intially.
These clients will know about the files they need by `file id` via
an alternate channel. The clients in turn will query `lfs-edge` with respective `file ids`.

`lfs-edge` will cache on first request or provide the file readily
if the file was already in local storage via a sync operation.

`lfs-edge` is aiming to reduce traffic seen by cloud file service especially
when it comes to large files.

### References
- [wiki](https://github.azc.ext.hp.com/Krypton/lfs-edge/wiki)
