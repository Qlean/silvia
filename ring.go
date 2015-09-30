package main

type RingItem struct {
	Body     string
	Event    *Event
	Info     string
}

type Ring struct {
	Ring []*RingItem
}

var successRing *Ring = &Ring{}
var failRing *Ring = &Ring{}

func (ring *Ring) Add(body string, nginxlog *NginxLog, event *Event, info string) {
	ringItem := &RingItem{
		Body: body,
		Event: event,
		Info: info,
	}

	ring.Ring = append(ring.Ring, ringItem)

	if len(ring.Ring) > serverConfig.RingSize { ring.Ring = ring.Ring[1:] }
}
