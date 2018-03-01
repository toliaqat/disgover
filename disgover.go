package disgover

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/golang/groupcache/lru"
	kbucket "github.com/libp2p/go-libp2p-kbucket"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"google.golang.org/grpc"
	grpcPeer "google.golang.org/grpc/peer"
)

/*
TODO: This was in disgo/server.go which needs to be in Disgover

type CmdParams struct {
	NodeId   string
	ThisIp   string
	SeedList []string
}


	var cmdParams = CmdParams{}

	// Parse CMD Args
	for _, arg := range os.Args {
		if strings.Index(arg, "-nodeId=") == 0 {
			cmdParams.NodeId = strings.Replace(arg, "-nodeId=", "", -1)
		} else if strings.Index(arg, "-thisIp=") == 0 {
			cmdParams.ThisIp = strings.Replace(arg, "-thisIp=", "", -1)
		} else if strings.Index(arg, "-seedList=") == 0 {
			var seedList = strings.Replace(arg, "-seedList=", "", -1)
			cmdParams.SeedList = strings.Split(seedList, ";")
		}
	}

	// Set THIS Contact/Node on the network
	var thisContact = disgover.NewContact()
	if len(cmdParams.NodeId) > 0 {
		thisContact.Id = cmdParams.NodeId
	}
	if len(cmdParams.ThisIp) > 0 {
		thisContact.Endpoint.Host = cmdParams.ThisIp
	}
	thisContact.Endpoint.Port = int64(properties.Properties.GrpcPort)

	// Check if we have a seed list
	var seedList = []*disgover.Contact{}
	for _, seedIP := range cmdParams.SeedList {
		seedList = append(seedList, &disgover.Contact{
			Endpoint: &disgover.Endpoint{
				Host: seedIP,
				Port: int64(properties.Properties.GrpcPort),
			},
		})
	}

	// Instantiate the node
	disgover.SetInstance(disgover.NewDisgover(thisContact, seedList))

*/


type Disgover struct {
	ThisContact *Contact
	Nodes       map[string]*Contact

	lruCache *lru.Cache
	kdht     *kbucket.RoutingTable
}

type IDisgover interface {
	Run()
	RunOnExisting(listener net.Listener)

	Find(contactId string, sender *Contact) (*Contact, error)
}

var DisgoverSingleton *Disgover = nil

func GetInstance() *Disgover {
	return DisgoverSingleton
}

func SetInstance(disgover *Disgover) {
	DisgoverSingleton = disgover
}

func NewDisgover(thisContact *Contact, seed []*Contact) *Disgover {
	seedNodes := map[string]*Contact{}
	for _, peer := range seed {
		seedNodes[peer.Id] = peer
	}

	disgoverInstance := &Disgover{
		ThisContact: thisContact,

		lruCache: lru.New(0),
		Nodes:    seedNodes,
		kdht: kbucket.NewRoutingTable(
			1000,
			kbucket.ConvertPeerID(peer.ID(thisContact.Id)),
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

func (disgover *Disgover) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", disgover.ThisContact.Endpoint.Port))
	if err != nil {
		log.Fatalf("unable to listen on %d: %v", disgover.ThisContact.Endpoint.Port, err)
	}

	disgover.RunOnExisting(listener)
}

func (disgover *Disgover) RunOnExisting(listener net.Listener) {
	server := grpc.NewServer()
	RegisterDisgoverRPCServer(server, disgover)
	go server.Serve(listener)

	disgover.Go()
}

func (disgover *Disgover) Go() {
	fmt.Println(fmt.Sprintf("Disgover-TRACE: Disgover[%s @ %s:%d]",
		disgover.ThisContact.Id,
		disgover.ThisContact.Endpoint.Host,
		disgover.ThisContact.Endpoint.Port,
	))

	log.WithFields(log.Fields{
		"method": "Disgover.Go",
	}).Info("running...")

	disgover.pingSeedList()
}

func (disgover *Disgover) PeerPing(ctx context.Context, contact *Contact) (*Contact, error) {
	thePeer, ok := grpcPeer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Disgover-TRACE: failed to get peer from ctx")
	}
	if thePeer.Addr == net.Addr(nil) {
		return nil, fmt.Errorf("Disgover-TRACE: failed to get peer address")
	}

	var peerAddressWithPort = thePeer.Addr.String()
	contact.Endpoint.Host = peerAddressWithPort[0:strings.Index(peerAddressWithPort, ":")]

	fmt.Println(fmt.Sprintf("Disgover-TRACE: PeerPing(): %s @ [%s : %d]",
		contact.Id,
		contact.Endpoint.Host,
		contact.Endpoint.Port,
	))

	disgover.addOrUpdate(contact)
	return disgover.ThisContact, nil
}

func (disgover *Disgover) PeerFind(ctx context.Context, findRequest *FindRequest) (*Contact, error) {
	return disgover.Find(findRequest.ContactId, findRequest.Sender)
}

func (disgover *Disgover) Find(contactId string, sender *Contact) (*Contact, error) {
	fmt.Println(fmt.Sprintf("Disgover-TRACE: Find(): %s in %s by %s", contactId, disgover.ThisContact.Id, sender.Id))

	if contact, ok := disgover.Nodes[contactId]; ok {
		return contact, nil
	}

	return disgover.findViaPeers(contactId, sender)
}

func (disgover *Disgover) findViaPeers(nodeID string, sender *Contact) (*Contact, error) {
	peerIDs := disgover.kdht.NearestPeers([]byte(disgover.ThisContact.Id), len(disgover.Nodes))

	for _, peerID := range peerIDs {
		peerIDAsString := string(peerID)
		if peerIDAsString == disgover.ThisContact.Id {
			continue
		}

		if peerIDAsString == sender.Id {
			continue
		}

		contact := disgover.Nodes[peerIDAsString]

		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", contact.Endpoint.Host, contact.Endpoint.Port), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("cannot dial server: %v", err)
		}

		client := NewDisgoverRPCClient(conn)
		respose, _ := client.PeerFind(context.Background(), &FindRequest{
			ContactId: nodeID,
			Sender:    sender,
		})

		fmt.Println("Disgover-TRACE: findViaPeers() RESULT")
		if respose != nil {
			fmt.Println(fmt.Sprintf("       %s, on [%s : %d]", respose.Id, respose.Endpoint.Host, respose.Endpoint.Port))
			disgover.addOrUpdate(contact)
			return respose, nil
		} else {
			fmt.Println("       NOT FOUND")
		}
	}

	return nil, nil
}

func (disgover *Disgover) addOrUpdate(contact *Contact) {
	disgover.Nodes[contact.Id] = contact
	disgover.kdht.Update(peer.ID(contact.Id))
}

func (disgover *Disgover) pingSeedList() {
	fmt.Println(fmt.Sprintf("Disgover-TRACE: pingSeedList()"))

	// peerIDs := disgover.kdht.NearestPeers([]byte(disgover.ThisContact.Id), len(disgover.Nodes))
	// contact := disgover.Nodes[peerID]

	for peerID, contact := range disgover.Nodes {
		if peerID == disgover.ThisContact.Id {
			continue
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", contact.Endpoint.Host, contact.Endpoint.Port), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("cannot dial server: %v", err)
		}

		client := NewDisgoverRPCClient(conn)
		seedNode, err := client.PeerPing(context.Background(), disgover.ThisContact)
		if err != nil {
			fmt.Println(fmt.Sprintf("Disgover-TRACE-ERROR: pingSeedList() -> %s", err))
		} else {
			disgover.Nodes[peerID].Id = seedNode.Id
		}
	}
}

func NewContact() *Contact {
	data := make([]byte, 10)
	rand.Read(data)

	return &Contact{
		Id: base64.StdEncoding.EncodeToString(data),
		Endpoint: &Endpoint{
			Port: 1975,
			Host: getLocalIP(),
		},
	}
}

func getLocalIP() string {
	return strings.Join(getLocalIPList(), ",")
}

func getLocalIPList() []string {
	var ipList = []string{}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			var ipAddress = ip.String()

			// var isUnspecified = ip.IsUnspecified()
			// var isLoopback = ip.IsLoopback()
			// var isMulticast = ip.IsMulticast()
			// var isInterfaceLocalMulticast = ip.IsInterfaceLocalMulticast()
			// var isLinkLocalMulticast = ip.IsLinkLocalMulticast()
			// var isLinkLocalUnicast = ip.IsLinkLocalUnicast()
			// var isGlobalUnicast = ip.IsGlobalUnicast()

			if ip.IsGlobalUnicast() {
				ipList = append(ipList, ipAddress)
			}
		}
	}

	return ipList

	// name, err := os.Hostname()
	// if err != nil {
	// 	fmt.Printf("Oops: %v\n", err)
	// 	return ""
	// }

	// addrs, err := net.LookupHost(name)
	// if err != nil {
	// 	fmt.Printf("Oops: %v\n", err)
	// 	return ""
	// }
	// fmt.Printf("Local IP: %s\n", addrs[0])

	// return addrs[0]
}
