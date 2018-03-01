package disgover

import (
	"net"
	"strconv"
	"testing"
)

var hostNumber = 0

// Verify that multiple nodes can find each other, and find method returns in
// reasonable time if non-existent node is being tried.
func TestMultiNodeDisgover(t *testing.T) {
	var nodes = setupNodes(3)
	joinAllNodes(nodes)

	// Verify that all nodes can find each other
	for _, peer1 := range nodes {
		for _, peer2 := range nodes {
			if canFindNode(peer1, peer2) == false {
				t.FailNow()
			}
		}
	}

	// Verify that disjoint nodes cannot be found and find routine returns in time.
	var newNodes = setupNodes(3)
	if canFindNode(nodes[0], newNodes[0]) == true {
		t.FailNow()
	}
}

func joinAllNodes(nodes []Disgover) {
	for _, peer := range nodes {
		joinNodes(peer, nodes)
	}
}

func joinNodes(node Disgover, peers []Disgover) {
	for _, peer := range peers {
		node.addOrUpdate(peer.ThisContact)
	}
}

func setupNodes(count int) ([]Disgover) {
	var nodes []Disgover
	for i := 0; i < count; i++ {
		hostNumber++
		var dsg = NewDisgover(
			&Contact{
				Id: "host-" + strconv.Itoa(hostNumber),
				Endpoint: &Endpoint{
					Host: "127.0.0.1",
					Port: getNewPort(),
				},
			},
			nil,
		)
		dsg.Run()
		nodes = append(nodes, *dsg)
	}
	return nodes
}

func canFindNode(node Disgover, nodeToFind Disgover) bool {
	peer, _ := node.Find(nodeToFind.ThisContact.Id, node.ThisContact)
	return peer != nil
}

func getNewPort() (int64) {
	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0
	}
	defer listener.Close()
	var port = int64(listener.Addr().(*net.TCPAddr).Port)
	return port
}
