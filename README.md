torronto - A BitTorrent Implementation in Go
============================================
## TODO
* race condition between leaving and downloading a file.
* bigger tests

## Torronto Messaging Documentation

_Using [json](www.json.org) for data interchange_

_Each message is a header with some max size_

Interface Actions
```
  Join   = 0
  Leave  = 1
  Insert = 2
  Query  = 3

```

Peer Actions
```
  Add      = 4
  Remove   = 5
  Files    = 6
  Download = 7
  Upload   = 8
  Have     = 9
```

Error Codes
```
  ErrOK            = 0
  ErrConnected     = 1
  ErrDisconnected  = 2
  ErrFileExists    = 3
  ErrFileMissing   = 4
  ErrBadPermission = 5
```

## Interface Messaging
_We tell the peer to join the network._
### Joining network
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 0,
}
```

_We tell the peer to leave the network_
### Leaving network
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 1
}
```

_Giving the peer the path to a file to insert_
### Inserting a file
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 2,
  "Files":
    [
      {
        "FileName": "<FileName>",
      }
    ]
}
```

_We tell the peer to leave the network_
### Querying status
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 3
}
```

_The peers responds with a status json object_
### Response to a query
```
{
  "numFiles": <numFiles>,
  "local": [<frActionPresentLocally>, <frActionPresentLocally>, ...],
  "system": [<frActionPresent>, <frActionPresent>, ...],
  "leastReplication": [<minimumReplicationLevel>, <minimumReplicationLevel>, ...],
  "weightedLeastReplication": [<averageReplicationLevel>, <averageReplicationLevel>, ...]
}
```

* * *
## Peer Messaging

_When a joins the network, it send out a message that it's joining and a list of its Files._
### Joining network
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 4,
  "Files":
    [
      {
        "file": "<FileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>, ...]"
      },
      {
        "file": "<FileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>, ...]"
      }, ...
    ]
}
```

### Leaving network
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 5
}
```

 _The peers respond with a file list._

### Returning file list
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 6,
  "Files":
    [
      {
        "FileName": "<FileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>, ...]"
      },
      {
        "FileName": "<FileName>",
        "chunks": "[<chunkNumber>, <chunkNumber>, ...]"
      }, ...
    ]
}
```

_Each peer will then update the status of the Files it has, and send out requests to download the Files it doesn't have_
### Request to download Files
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 7,
  "Files":
    [
      {
        "FileName": "<FileName>",
        "chunks": [<Filesize>, <chunkNumber>]
      }
    ]
}
```

### Sending a file chunk (this is followed by the payload)
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 8,
  "Files":
    [
      {
        "FileName": "<FileName>",
        "chunks": [<Filesize>, <chunkNumber>]
      }
    ]
}
```

 _Peer messages whenever it receives a new chunk_

### Have a new file
```
{
  "HostName": "<HostName>",
  "PortNumber": "<PortNumber>",
  "Action": 9,
  "Files":
    [
      {
        "FileName": "<FileName>",
        "chunks": [<Filesize>, <chunkNumber>]
      }
    ]
}
```

* * *
