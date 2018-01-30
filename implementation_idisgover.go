package disgover

func (disgover *Disgover) Run() {
	disgover.Transport.Listen()
}

func (disgover *Disgover) Find(nodeId string, sender *Contact) (contact *Contact, err error) {
	if contact, ok := disgover.nodes[nodeId]; ok {
		return contact, nil
	}

	return disgover.findViaPeers(nodeId, sender)
}
