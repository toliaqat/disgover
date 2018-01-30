package disgover

import (
	"fmt"
)

func (disgover *Disgover) Run() {
	disgover.Transport.Listen()
	disgover.pingSeedList()
}

func (disgover *Disgover) Find(nodeID string, sender *Contact) (contact *Contact, err error) {
	fmt.Println(fmt.Sprintf("TRACE: Find(): %s", nodeID))

	if contact, ok := disgover.nodes[nodeID]; ok {
		return contact, nil
	}

	return disgover.findViaPeers(nodeID, sender)
}
