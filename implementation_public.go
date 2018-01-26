package main

import (
	"errors"
	"time"

	"github.com/golang/groupcache/lru"
	"github.com/libp2p/go-libp2p-kbucket"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

// NewDisgover -
func NewDisgover(bucketsize int, localID string, latency time.Duration, transport interface{}) *Disgover {
	return &Disgover{
		InlineTrace: false,
		Seeds:       []*Contact{},
		Transport:   transport,

		lruCache: lru.New(0),
		buckets:  map[string]*Contact{},
		routingTable: kbucket.NewRoutingTable(
			bucketsize,
			[]byte(localID),
			latency,
			pstore.NewMetrics(),
		),

		CONCURRENCY_CONSTANT: 3,
		maxCacheSize: 1000,
	}
}

//
// Add -
//		`remoteContact`	: object to add and is not managed
//		`id`			: base64 contact id
//		`data`			: data to be included with the contact,
//						  will be returned for anyone querying for this `contact` by `id`.
//
//		Return: `Contact` that was added
//
func (disgover *Disgover) Add(remoteContact *Contact) (*Contact, error) {
	if len(remoteContact.Id) == 0 {
		return nil, errors.New("Invalid or missing contact.Id")
	}

	// even if we don't have buckets to update, we can still store information
	// in LRU cache (check using arbiter to update cache with latest only)
	var cached, ok = disgover.lruCache.Get(remoteContact.Id)
	var cachedAsContant = cached.(*Contact)
	var selection = disgover.arbiter(cachedAsContant, remoteContact)
	if selection != cached {
		disgover.lruCache.Remove(remoteContact.Id)
		disgover.lruCache.Add(remoteContact.Id, remoteContact)
	}

	if len(disgover.buckets) == 0 {
		return nil, nil // no k-buckets to update
	}

	// first, check if remote contact id is locally registered
	if _, ok = disgover.buckets[remoteContact.Id]; ok {
		// remote contact id is same as locally registered contact id
		// need to arbiter which contact version should be retained
		// (we already calculated the selection)
		disgover.buckets[remoteContact.Id] = selection
	} else {
		// we pick the closest kBucket to the node id of our contact to store the
		// data in, since they have the most space to accomodate near-by node ids
		// (inherent KBucket property)

		id, _ := peer.IDFromString(remoteContact.Id)
		disgover.routingTable.Update(id)
	}

	return remoteContact, nil
}

// ExecuteQuery -
//		`query`		: containing query state for this request
//		`nodeId`	: base64 node id to find
//		`nodes`		: contacts to query for `nodeId` arranged from closest to farthest
//		`node`		:
//			`id`		: base64 contact id.
//			`nodesMap`	: dictionary to the same contacts already present in `nodes` for O(1) access.
//		`callback`	: result callback to call with
func (disgover *Disgover) ExecuteQuery(query string, callback QueryCallback) {
	if query.done {
		return
	}

    // if we have no nodes, we can't query anything
	if len(query.nodes) == 0 {
		callback(error.Error("No known nodes to query"))
		return
	}

	if query.index < 0 {
		query.index = 0
	}
	if query.closest == nil {
		query.closest = query.nodes[0];
	}
    if query.ongoingRequests < 0 {
		query.ongoingRequests = 0
	}
    // if (query.newNodes === undefined)
    //     query.newNodes = [];

    // we listen for `node` events that contain the nodeId we asked for
    // this helps to decouple discover from the transport and allows us to
    // benefit from other ongoing queries (TODO: "prove" this)
    //
    // because executeQuery can be called multiple times on the same query,
    // we keep the state
    if (!query.listener) {
        // TODO: maybe there is an opportunity here to generate events
        // uniquely named by "nodeId" so I don't have to have tons of listeners
        // listen to everything and throw away what they don't want?
        query.listener = function (error, contact, nodeId, response) {
            // filter other queries
            if (nodeId != query.nodeId)
                return;

            // query already successfully completed
            if (query.done)
                return;

            // request has been handled
            // TODO: what happens if two requests for the same nodeId are
            //       happening at the same time?
            // maybe do a check prior to executeQuery to not duplicate searches
            // for the same nodeId across the network?
            query.ongoingRequests--;

            if (error) {
                if (disgover.tracing) {
                    disgover.trace('error response from ' + util.inspect(contact) +
                        ' looking for ' + nodeId + ': ' + util.inspect(error));
                }
                var contactRecord = query.nodesMap[contact.id];

                if (!contactRecord)
                    return;

                if (contactRecord.kBucket) {
                    // we have a kBucket to report unreachability to
                    // remove from kBucket
                    var kBucketInfo = disgover.kBuckets[contactRecord.kBucket.id];
                    if (!kBucketInfo) {
                        return;
                    }

                    var kBucket = kBucketInfo.kBucket;
                    if (!kBucket) {
                        return;
                    }

                    var contactRecordToRemove = clone(contactRecord);
                    contactRecordToRemove.id =
                        new Buffer(contactRecord.id, 'base64');
                    kBucket.remove(contactRecordToRemove);
                }

                contactRecord.contacted = true;

                // console.dir(query);

                // initiate next request if there are still queries to be made
                if (query.index < query.nodes.length
                    && query.ongoingRequests < disgover.CONCURRENCY_CONSTANT) {
                    process.nextTick(function () {
                        disgover.executeQuery(query, callback);
                    });
                } else {
                    disgover.queryCompletionCheck(query, callback);
                }
                return; // handled error
            }

            // we have a response, it could be an Object or Array

            if (disgover.tracing) {
                disgover.trace('response from ' + util.inspect(contact) +
                    ' looking for ' + nodeId + ': ' + util.inspect(response));
            }
            if (Array.isArray(response)) {
                // add the closest contacts to new nodes
                query.newNodes = query.newNodes.concat(response);

                // TODO: same code inside error handler
                // initiate next request if there are still queries to be made
                if (query.index < query.nodes.length
                    && query.ongoingRequests < disgover.CONCURRENCY_CONSTANT) {
                    process.nextTick(function () {
                        disgover.executeQuery(query, callback);
                    });
                } else {
                    disgover.queryCompletionCheck(query, callback);
                }
                return;
            }

            // we have a response Object, found the contact!
            // add the new contact to the closestKBucket
            var finalClosestKBuckets = disgover.getClosestKBuckets(response.id);
            if (finalClosestKBuckets.length > 0) {
                var finalClosestKBucket =
                    disgover.kBuckets[finalClosestKBuckets[0].id].kBucket;
                var contact = clone(response);
                contact.id = new Buffer(contact.id, "base64");
                finalClosestKBucket.add(contact);
            }

            // return the response and stop querying
            var latency = disgover.timerEndInMilliseconds('find.ms', nodeId);
            var roundLatency = disgover.timerEndInMilliseconds('find.round.ms', nodeId);
            callback(null, response);
            query.done = true;
            disgover.transport.removeListener('node', query.listener);
            disgover.emit('stats.timers.find.ms', latency);
            disgover.emit('stats.timers.find.round.ms', roundLatency);
            return;
        };
        disgover.transport.on('node', query.listener);
    }

	for ; (query.index < len(query.nodes)) & (query.ongoingRequests < disgover.CONCURRENCY_CONSTANT); query.index++) {
        query.ongoingRequests++;
        disgover.transport.findNode(query.nodes[query.index], query.nodeId, query.sender);
    }

    disgover.queryCompletionCheck(query, callback);
}
