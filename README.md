# P2P messaging simulation toolkit
---
This repository holds different simulators for exploring and researching p2p networks and messaging related to Status.

Original intent of these simulators is to provide stats and resulting traces/logs for further analysis and visualization.

## Design
```

+------------------+   +----------------+   +-------------+   +------------------+                 
| Choose           |   |  Run nodes in  |   |             |   |                  |                 
| network topology |----  simulated     |   |             |   |                  |                 
+------------------+   |  environment   |   | Propagate   |   | Collect network  |                 
                       |   - in-memory  |---- message(s)  |---- events &         |                 
+------------------+   |   - exec       |   |             |   | generate stats   |                 
|  Choose          |----   - docker     |   |             |   |                  |                 
|  Simulator       |   |                |   |             |   |                  |                 
+------------------+   +----------------+   +-------------+   +------------------+                 
```

### Simulators support

| Simulator   | State | Description |
|---|---|---|
| **WhisperV6** | Done | Master branch if go-ethereum Whisper implementation  |
| **Naive**  | Done | Naive p2p propagation  |
| PSS | TBD | Swarm's PSS messaging |

### Network environments support

| Node type  | State | Description |
|---|---|---|
| **In-Memory** | Done | Single node in-memory network  |
| Exec  | TBD | Single node native binary network with localhost connection |
| Docker | TBD | Docker-based network |

## Usage


## License
MIT
