torronto - A BitTorrent Implementation in Go
============================================
## TODO

_Start testing_
`go build`

`./torronto <hostName>:<portNumber>`

* error checking
* ask for missing files
* replication structure
* send out unique chunks when joining and leaving

## Torronto Messaging Documentation

_Using [json](www.json.org) for data interchange_

_Each message is a header with some max size_

## Interface Messaging
_We tell the node to join the network._
### Joining network
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "Join",
}
```

_We tell the node to leave the network_
### Leaving network
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "Leave"
}
```

_We tell the node to leave the network_
### Querying status
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "Query"
}
```

_Giving the node the path to a file to insert_
### Inserting a file
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "Insert",
  "files":
    [
      {
        "fileName": "<fileName>",
      }
    ]
}
```
* * *
## Peer Messaging

_When a joins the network, it send out a message that it's joining and a list of its files._
### Joining network
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "Add",
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
  "action": "Remove"
}
```

 _The peers respond with a file list._

### Returning file list
```
{
  "hostName": "<hostName>",
  "portNumber": "<portNumber>",
  "action": "Files",
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
  "action": "Download",
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
  "action": "Upload",
  "files":
    [
      {
        "fileName": "<fileName>",
        "chunks": [<fileSize>, <chunkNumber>]
      }
    ]
}
```

 _Peer messages whenever it receives a new chunk_

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

* * *
