torronto - A BitTorrent Implementation in Go
============================================
## TODO

* error checking
* ask for missing files
* replication structure
* send out unique chunks when joining and leaving

## Torronto Messaging Documentation
* * *

_Using [json](www.json.org) for data interchange_

_Each message is a header with some max size_
_
_When a joins the network, it send out a message that it's joining and a list of its files._
### Joining network
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "join",
  "files":
    [
      {
        "file": "<fileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>, ...]"
      },
      {
        "file": "<fileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>, ...]"
      }, ...
    ]
}
```

### Leaving network
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "leave"
}

```

* * *
 _The peers respond with a file list._

### Have a new file
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "Have",
  "files":
    [
      {
        "fileName": "<fileName>",
        "chunks": [<fileSize>, <chunkNumber>]
      }
    ]
}
```

### Returning file list
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "files",
  "files":
    [
      {
        "fileName": "<fileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>, ...]"
      },
      {
        "fileName": "<fileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>, ...]"
      }, ...
    ]
}
```

_Each peer will then update the status of the files it has, and send out requests to download the files it doesn't have_
### Request to download files
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "download",
  "files":
    [
      {
        "fileName": "<fileName>",
        "chunks": [<fileSize>, <chunkNumber>]
      }
    ]
}
```

### Sending a file chunk (this is followed by the payload)
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "upload",
  "files":
    [
      {
        "fileName": "<fileName>",
        "chunks": [<fileSize>, <chunkNumber>]
      }
    ]
}
```

* * *

_Incomplete files are saved with their chunk number in the name_
.`<fileName>:<chunkNumber>`
