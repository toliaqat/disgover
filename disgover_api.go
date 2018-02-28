package disgover

import (
	"github.com/dispatchlabs/disgo_commons/types"
	"fmt"
	"github.com/golang/groupcache/lru"
	"github.com/libp2p/go-libp2p-kbucket"
	"github.com/libp2p/go-libp2p-peerstore"
	peer "github.com/libp2p/go-libp2p-peer"
	"sync"
)

var once sync.Once
var disgover *Disgover

type Disgover struct {
	ThisContact *types.Contact
	Nodes       map[string]*types.Contact
	lruCache 	*lru.Cache
	kdht     	*kbucket.RoutingTable
}

//TODO: Need to add common code for seed list
func getDisgover() *Disgover {
	once.Do(func() {
		disgover = NewDisgover(types.NewContact(), []*types.Contact{})
		disgover.addOrUpdate(disgover.ThisContact)
	})
	return disgover
}

func NewDisgover(thisContact *types.Contact, seed []*types.Contact) *Disgover {
	seedNodes := map[string]*types.Contact{}
	for _, peer := range seed {
		seedNodes[peer.Address] = peer
	}

	disgoverInstance := &Disgover{
		ThisContact: thisContact,

		lruCache: lru.New(0),
		Nodes:    seedNodes,
		kdht: kbucket.NewRoutingTable(
			1000,
			kbucket.ConvertPeerID(peer.ID(thisContact.Address)),
			1000,
			peerstore.NewMetrics(),
		),
	}
	disgoverInstance.addOrUpdate(thisContact)

	for _, contact := range seed {
		disgoverInstance.addOrUpdate(contact)
	}
	return disgoverInstance
}

func (disgover *Disgover) addOrUpdate(contact *types.Contact) {
	id := contact.Address
	getDisgover().Nodes[id] = contact
	getDisgover().kdht.Update(peer.ID(id))
}

func (disgover *Disgover) pingSeedList() {
	fmt.Println(fmt.Sprintf("Disgover-TRACE: pingSeedList()"))

	for peerID, contact := range getDisgover().Nodes {
		if peerID == getDisgover().ThisContact.Address {
			continue
		}
		seedNode := PeerPingWithGrpcClient(contact, getDisgover().ThisContact)
		getDisgover().Nodes[peerID].Address = seedNode.Address
		getDisgover().addOrUpdate(seedNode)
	}
}

func (disGoverService *DisGoverService) Find(contactId string, sender *types.Contact) (*types.Contact, error) {
	fmt.Println(fmt.Sprintf("Disgover-TRACE: Find(): %s", contactId))

	if contact, ok := getDisgover().Nodes[contactId]; ok {
		return contact, nil
	}
	return getDisgover().findViaPeers(contactId, sender)
}

func (disGoverService *DisGoverService) FindAll() (*[]types.Contact, error) {
	fmt.Println(fmt.Sprintf("Disgover-TRACE: FindAll()"))
	peerIDs := getDisgover().kdht.NearestPeers([]byte(getDisgover().ThisContact.Address), len(getDisgover().Nodes))
	contacts := make([]types.Contact, 0)
	for _, peerID := range peerIDs {
		peerIDAsString := string(peerID)
		contact := getDisgover().Nodes[peerIDAsString]
		contacts = append(contacts, *contact)
	}
	return &contacts, nil
}

func (disGoverService *DisGoverService) PeerFind(idToFind string, contact *types.Contact) (*types.Contact, error) {
	fmt.Println(fmt.Sprintf("Disgover-TRACE: PeerFind(): %s", contact.Address))

	if contact, ok := getDisgover().Nodes[contact.Address]; ok {
		return contact, nil
	}
	return getDisgover().findViaPeers(contact.Address, contact)
}

func (disgover *Disgover) findViaPeers(idToFind string, sender *types.Contact) (*types.Contact, error) {
	fmt.Println(fmt.Sprintf("Disgover-TRACE: findViaPeers(): %s", idToFind))

	peerIDs := getDisgover().kdht.NearestPeers([]byte(getDisgover().ThisContact.Address), len(getDisgover().Nodes))

	for _, peerID := range peerIDs {
		peerIDAsString := string(peerID)
		if peerIDAsString == getDisgover().ThisContact.Address {
			continue
		}

		contact := getDisgover().Nodes[peerIDAsString]
		respose := FindPeerWithGrpcClient(idToFind, sender)

		if respose != nil {
			fmt.Println(fmt.Sprintf(" %s, on [%s : %d]", respose.Address, respose.Endpoint.Host, respose.Endpoint.Port))

			getDisgover().addOrUpdate(contact)
			return respose, nil
		} else {
			fmt.Println("       NOT FOUND")
		}
	}
	return nil, nil
}
