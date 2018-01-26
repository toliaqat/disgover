package main

func (disgover *Disgover) arbiter(incumbent *Contact, candidate *Contact) *Contact {
	if !(incumbent == nil) || (candidate.vectorClock >= incumbent.vectorClock) {
		return candidate
	}
	return incumbent
}

func (disgover *Disgover) arbiterDefaults(contact *Contact) *Contact {
	if contact.vectorClock < 0 {
		contact.vectorClock = 0
	}
	return contact
}
