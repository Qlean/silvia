package main

type RingItem struct {
	Body          string
	AdjustEvent   *AdjustEvent
	SnowplowEvent *SnowplowEvent
	Info          string
}

type Ring struct {
	Ring  []*RingItem
	Total int
	Size  int
}

func (ring *Ring) AddAdjust(body string, nginxlog *NginxLog, event interface, info string) {
	ringItem := &RingItem{
		Body: body,
		AdjustEvent: event,
		Info: info,
	}

	ring.Ring = append(ring.Ring, ringItem)
	ring.Total++

	if len(ring.Ring) > ring.Size { ring.Ring = ring.Ring[1:] }
}
