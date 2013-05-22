package event

// receiver just implements the Receiver Interface.
// It can be used to make any class into an event
// receiver easily just by mixin in this class.
type receiver struct {
	evtHandler     func(interface{})
	evtChanHandler func(<-chan interface{})
	evtChan        chan interface{}
}

func NewReceiver() Receiver {
	return new(receiver)
}

func NewReceiverFor(sender Sender) Receiver {
	rec := NewReceiver()
	rec.ListenTo(sender)
	return rec
}

func (r *receiver) ListenTo(sender Sender) {
	sender.AddReceiver(r)
}

func (r *receiver) UnlistenTo(sender Sender) {
	sender.RemoveReceiver(r)
}

// SetHandler sets a function set is used to handle all
// received events.
//
// An alternative method for handling events is to mix-in Receiver
// into another class and to overwrite the ReceiveEvent method.
// However this will totally change the default behaviour.
// So it is not recommended.
func (r *receiver) SetEvtHandler(handler func(interface{})) {
	r.evtHandler = handler
}

func (r *receiver) SetEvtChanHandler(handler func(evtChan <-chan interface{})) {
	// XXX Dirk: Buffer should grow as needed.
	// A buffer that is too small could cause a deadlock.
	r.evtChan = make(chan interface{}, 100)
	r.evtChanHandler = handler
	go r.evtChanHandler(r.evtChan)
}

func (r *receiver) ReceiveEvent(evt interface{}) {
	if r.evtHandler != nil {
		r.evtHandler(evt)
	}
	if r.evtChanHandler != nil {
		r.evtChan <- evt
	}
}
