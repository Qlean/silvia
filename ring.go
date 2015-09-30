package main

import "log"

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

func (ring *Ring) Add(body []byte, nginxlog *NginxLog, event *Event, info string) {

	bodyString := string(body)

	log.Println(bodyString)

	ringItem := &RingItem{
		Body: bodyString,
		Event: event,
		Info: info,
	}

	ring.Ring = append(ring.Ring, ringItem)

	if len(ring.Ring) > serverConfig.RingSize { ring.Ring = ring.Ring[1:] }
}
