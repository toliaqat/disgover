package disgover

import (
	"github.com/golang/groupcache/lru"
	"github.com/libp2p/go-libp2p-kbucket"
)

type Endpoint struct {
	Host string
	Port int64
}

type Contact struct {
	Id       string
	Endpoint Endpoint
	Data     interface{}
}

type DisgoverRpc struct {
	Request string
	NodeId  string
}

// Transport
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
type DisgoverRpcDelegate func(data []byte) (result []byte, err error)

type ITransport interface {
	ExecRPC(destination *Contact, rpc DisgoverRpc) []byte
	OnPeerRPC(delegate DisgoverRpcDelegate)

	Listen()
}

// Disgover
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
type Disgover struct {
	Contact   *Contact
	Transport ITransport

	lruCache *lru.Cache
	nodes    map[string]*Contact
	kdht     *kbucket.RoutingTable
}

type IDisgover interface {
	Run()

	Find(nodeId string, sender *Contact) (contact *Contact, err error)
}
