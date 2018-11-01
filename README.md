# P2P messaging simulation toolkit
---
This repository holds different simulators for exploring and researching p2p networks and messaging related to Status.

Original intent of these simulators is to provide stats and resulting traces/logs for further analysis and visualization.

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

See READMEs for each package/program separately.

## Propagation
Propagation simulators send message over known network topology and record its propagation data/stats.

Currently implemented:

| Simulator   | Description |
|---|---|
| WhisperV6  | Master branch if go-ethereum Whisper implementation  |
| Naive  | Naive p2p propagation  |

## Usage


## License
MIT
