package event

// sender just implements the Sender interface.
type sender struct {
	receivers map[Receiver]bool
}

func NewSender() Sender {
	sender := new(sender)
	sender.receivers = make(map[Receiver]bool)
	return sender
}

func (s *sender) AddReceiver(receiver Receiver) {
	s.receivers[receiver] = true
}

func (s *sender) RemoveReceiver(receiver Receiver) {
	delete(s.receivers, receiver)
}

func (s *sender) SendEvent(evt interface{}) {
	for receiver, _ := range s.receivers {
		receiver.ReceiveEvent(evt)
	}
}
