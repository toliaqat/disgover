package disgover

import (
	"encoding/json"
	"fmt"

	"github.com/golang/groupcache/lru"
	kbucket "github.com/libp2p/go-libp2p-kbucket"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

// NewDisgover -
func NewDisgover(contact *Contact, seed []*Contact, transport ITransport) *Disgover {
	seedNodes := map[string]*Contact{}
	for _, peer := range seed {
		seedNodes[peer.Id] = peer
	}

	disgover := &Disgover{
		Contact:   contact,
		Transport: transport,

		lruCache: lru.New(0),
		nodes:    seedNodes,
		kdht: kbucket.NewRoutingTable(
			1000,
			kbucket.ConvertPeerID(peer.ID(contact.Id)),
			1000,
			peerstore.NewMetrics(),
		),
	}

	disgover.Transport.OnPeerRPC(func(data []byte) (result []byte, err error) {
		return disgover.onPeerRPC(data)
	})

	disgover.addOrUpdate(contact)

	for _, contact := range seed {
		disgover.addOrUpdate(contact)
	}

	return disgover
}

func (disgover *Disgover) findViaPeers(nodeID string, sender *Contact) (contact *Contact, err error) {
	fmt.Println(fmt.Sprintf("TRACE: findViaPeers(): %s", nodeID))

	peerIDs := disgover.kdht.NearestPeers([]byte(disgover.Contact.Id), len(disgover.nodes))

	payload := DisgoverRpc{
		Request: "findNode",
		Contact: Contact{
			Id: nodeID,
		},
	}

	for _, peerID := range peerIDs {
		peerIDAsString := string(peerID)
		if peerIDAsString == disgover.Contact.Id {
			continue
		}

		respose := disgover.Transport.ExecRPC(disgover.nodes[peerIDAsString], payload)
		if len(respose) != 0 {
			fmt.Println("TRACE: findViaPeers() RESULT")
			fmt.Println("       ", string(respose[:]))

			contact = &Contact{}
			json.Unmarshal(respose, contact)

			if len(contact.Id) == 0 {
				contact = nil
			} else {
				disgover.addOrUpdate(contact)
			}

			return
		}
	}

	return nil, nil
}

func (disgover *Disgover) onPeerRPC(data []byte) (result []byte, err error) {
	fmt.Println("TRACE: onPeerRPC()")
	fmt.Println("      ", string(data[:]))

	rpc := DisgoverRpc{}
	err = json.Unmarshal(data, &rpc)

	if err == nil {
		if rpc.Request == "findNode" {
			node, err := disgover.Find(rpc.Contact.Id, disgover.Contact)
			if err != nil {
				return nil, err
			}

			bytes, _ := json.Marshal(node)
			return bytes, nil
		}
		if rpc.Request == "ping" {
			disgover.addOrUpdate(&rpc.Contact)
		}
	} else {
		return nil, err
	}

	return nil, nil
}

func (disgover *Disgover) addOrUpdate(contact *Contact) {
	disgover.nodes[contact.Id] = contact
	disgover.kdht.Update(peer.ID(contact.Id))
}

func (disgover *Disgover) pingSeedList() {
	fmt.Println(fmt.Sprintf("TRACE: pingSeedList()"))

	peerIDs := disgover.kdht.NearestPeers([]byte(disgover.Contact.Id), len(disgover.nodes))

	payload := DisgoverRpc{
		Request: "ping",
		Contact: *disgover.Contact,
	}

	for _, peerID := range peerIDs {
		peerIDAsString := string(peerID)
		if peerIDAsString == disgover.Contact.Id {
			continue
		}

		respose := disgover.Transport.ExecRPC(disgover.nodes[peerIDAsString], payload)
		if len(respose) != 0 {
			fmt.Println("TRACE: pingSeedList() RESULT")
			fmt.Println("       ", string(respose[:]))
		}
	}
}
