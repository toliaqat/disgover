package main

import (
	"github.com/golang/groupcache/lru"
	"github.com/libp2p/go-libp2p-kbucket"
)

type Contact struct {
	Id       string
	Data     interface{}
	Endpoint interface{}

	vectorClock int64
}

type Query struct {
	done            bool
	nodeId          string
	index           int64
	ongoingRequests int64
	closest         *Contact
	nodes           []*Contact
	nodesMap        map[string]*Contact
	sender          *Contact
	newNodes        []*Contact
	listener
}
type QueryCallback func(err error, contact *Contact)
type NodeFoundListener func(err error, contact *Contact, nodeId string, response)

type ITransport interface {
	findNode(contact *Contact, nodeId string, sender *Contact)
}
type Disgover struct {
	InlineTrace bool
	Seeds       []*Contact
	Transport   ITransport

	lruCache     *lru.Cache
	buckets      map[string]*Contact
	routingTable *kbucket.RoutingTable

	CONCURRENCY_CONSTANT int64
	maxCacheSize         int64
}
