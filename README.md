# Dispatch KDHT based node discovery engine
Distributed, node discovery mechanism that enables locating any 
entity (server, worker, drone, actor) based on node id.

The intent is to not be a data storage/distribution mechanism.
Meaning we implement only `PING` and `FIND-NODE` rpc.

One `disgover` instance in the node:
- stores info about numerous nodes
- functions as a gateway to outside local network

An example and to illustrate perhaps less familiar intents, if we 
want to find an actor "receptionist" in the global actor system, 
we could look for a contact that looks like this:
```javascript
var contact = {
    id: "tmqjRAfBILbEC6aaHoz3AurtluM=", // Base64 encoded receptionist address
    data: {
        webkey: "c9bf857b35ed4750ca35c0a4f41e56644df59547",
        host: "10.13.211.201",
        port: 9999,
        publicKey: "mQINBFJhVUwBEADRwsK6hvXoZU/niqZU2k9NXVNA9kAiVBfhUZ...WYco9YzK2K1Q="
    },
    transport: {
        host: "10.13.211.201",
        port: 6742
    }
};
```

# Core Concept
The Conceptual part is around
```go
type Contact struct {
	Id       string
	Endpoint Endpoint
	Data     interface{}
}
```

- `Id` - node ID
- `Endpoint` - where the node can be contacted
- `Data` - any additional paypload needed for extended functionalities and queries


# Samples
![Test](nodes.png "Test")

