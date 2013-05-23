torronto - A BitTorrent Implementation in Go
============================================

## Torronto Messaging Documentation
* * *

_Using [json](www.json.org) for data interchange_

_Each message is a header with some max size_

### Joining network
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "join"
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
_When a joins the network, it send out a message requesting a list of files from all of it's peers. The peers respond with a file list, and send a request for the new node's file list._

### Fetch file list
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "fetch"
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
        "file": "<fileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>]"
      },
      {
        "file": "<fileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>]"
      }
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
  "file": "file"
}
```

### Sending a file chunk (this is followed by the payload)
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "upload",
  "file": "<fileName>",
  "chunk": "<chunkNumber>",
}
```
