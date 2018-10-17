# Propagation simulator server
---

This simulator server implements an API for running message propagation through network.

Currently supported:
 - whisperv6
 - naive p2p propagation


Server expects a network topology as an input, and returns propagation log data.


# Installation

```
go get github.com/status-im/simulation/cmd/propagation_server
```

# Usage

Just run:
```
propagation_server
```

You can specify different bind address using `-h` command line flag. See `./propagation_server -h` for usage info.


# Request JSON

``json
{
  "nodes": [
    {
      "id": "192.168.1.2"
    },
    {
      "id": "192.168.1.4"
    }
  ],
  "links": [
    {
      "source": "192.168.1.2",
      "target": "192.168.1.4"
    }
  ]
}
```

# Response format

Plog (propagation log)
TBD (see code)
