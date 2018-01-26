# `disgover` - Dispatch KDHT based resource discovery engine
Distributed, master-less, node discovery mechanism that enables locating any entity (server, worker, drone, actor) based on node id.
The intent is to not be a data storage/distribution mechanism.
Meaning we implement only `PING` and `FIND-NODE` rpc.

One `disgover` instance in the node:
- stores info about numerous nodes
- functions as a gateway to outside local network

Data Strcutre
```javascript
var contact = {
	id: string
	data: {} // any data needed when it is retrieved by others on the network
	transport: {} // info that the transport mechanism requires for operation
}
```
Node ids in `disgover` are represented as base64 encoded Strings. This is because the default generated node ids (20 random bytes) could be unsafe to print. base64 encoding was picked over hex encoding because it takes up less space when printed or serialized in ASCII over the wire.


Example:
```javascript
var contact = {
	id: "Zm9v",
	data: "fooo",
	transport: {
		host: "foo.bar.com", // or "localhost", "127.0.0.1", etc...
		port: 6742
	}
}
```

- `contact.transport` content
	- Only required for contacts that are seeds, meaning their transport information is known ahead of time so that a `disgover` node can connect to them.
	- For non-seed contacts, the content will be provided by the particular transport implementation.
- `contact.data` content
	- Intent is to support the discovery mechanism, by storing a minimal amount of information required for connecting to the node endpoint for the application's purpose
	- Given that `contact.transport` contains information for how a `disgover` transport can connect to another `disgover` transport, this is not very useful if one is trying to figure out the endpoint address of another node for application level purposes. It may not correspond at all to what's in `contact.transport`.

For example, if we want a DNS-like functionality, we could look for a contact with id of `my.secret.dns.com`. This could correspond to the following contact:
```javascript
var contact = {
	id: "bXkuc2VjcmV0LmRucy5jb20=", // Base64 encoding of "my.secret.dns.com"
	data: {
		host: "10.22.1.37",
		port: 8080
	},
	transport: {
		host: "10.22.1.37",
		port: 6742
	}
}
```
This would tell us that we can connect to `my.secret.dns.com` at IP address `10.22.1.37` and port 8080.

As another example and to illustrate perhaps less familiar intents, if we want to find an actor "receptionist" in the global actor system, we could look for a contact that looks like this:
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
This would tell us that we can access the actor using the published webkey at IP address `10.13.211.201` and port `9999` and to encrypt our communication using provided public key.

Uses of contact.data that are not "minimal" in this way can result in poor system behavior.

# Arbiter
A conflict resolution mechanism using an `arbiter` function.
It can choose between two contact objects with the same id but different properties and determine which one should be stored. As the `arbiter` function returns the actual object to be stored, it does not need to make an either/or choice, but instead could perform some sort of operation and return the result as a new object that would then be stored. 

`arbiterDefaults` function makes sure that `contact` has the appropriate defualt properties for the arbiter function to work correctly.

`arbiter` function is used in three places
- as the `k-bucket` arbiter function
- to determine whether a new `remote contact` should be inserted into the `LRU` cache
- to determine if unregistering a contact will succeed
	- if `arbiter` returns `contact` === to the `stored contact` and `stored contact` !== `contact` we want to unregister, then unregister will fail

Example `arbiter` function implementing a `vectorClock` mechanism
```javascript
var contact = {
    id: new Buffer('contactId'),
    vectorClock: 0
};

function arbiterDefaults(contact) {
    if (!contact.vectorClock) {
        contact.vectorClock = 0;
    }
    return contact;
};

function arbiter(incumbent, candidate) {
    if (!incumbent
        || (incumbent && !incumbent.vectorClock)
        || (incumbent && incumbent.vectorClock && candidate.vectorClock
            && (candidate.vectorClock >= incumbent.vectorClock))) {

        return candidate;
    }
    return incumbent;
};
```

Example `arbiter` that implements a`Grow-Only-Set` `CRDT` mechanism.
Assuming that each worker node has a globally unique id and that each value for a worker node id will be written only once.
```javascript
var contact = {
    id: new Buffer('workerService'),
    data: {
        workerNodes: {
            '17asdaf7effa2': { host: '127.0.0.1', port: 1337 },
            '17djsyqeryasu': { host: '127.0.0.1', port: 1338 }
        }
    }
};

function arbiterDefaults(contact) {
    if (!contact.data) {
        contact.data = {};
    }
    if (!contact.data.workerNodes) {
        contact.data.workerNodes = {};
    }
    return contact;
}

function arbiter(incumbent, candidate) {
    if (!incumbent || !incumbent.data || !incumbent.data.workerNodes) {
        return candidate;
    }

    if (!candidate || !candidate.data || !candidate.data.workerNodes) {
        return incumbent;
    }

    // we create a new object so that our selection is guaranteed to replace
    // the incumbent
    var merged = {
        id: incumbent.id, // incumbent.id === candidate.id within an arbiter
        data: {
            workerNodes: incumbent.data.workerNodes
        }
    };

    Object.keys(candidate.data.workerNodes).forEach(function (workerNodeId) {
        merged.data.workerNodes[workerNodeId] =
            candidate.data.workerNodes[workerNodeId];
    });

    return merged;
}
```