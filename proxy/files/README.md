# Files handling controller

## Methods

- GetFileInfo()

Provides file static info like size, mod time, cksum.

- GetFileContent()

Provides file descriptor opened on the file to be rendered by the handler.
The Range part extraction and file type detection is handled by the handler
using standard library function http.ServeContent() so no additional
file chopping is needed. IF file does not exist locally, the sync rpc service
is conntacted in order to start async download. In this case file is
registered by sync as pending and rest server retuerns Too Manu Requests
response.